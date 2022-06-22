package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/the-zion/matrix-core/app/user/service/internal/biz"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

var _ biz.UserRepo = (*userRepo)(nil)

var userCacheKey = func(username string) string {
	return "user_" + username
}

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "user/data/user")),
	}
}

func (r *userRepo) GetUser(ctx context.Context, id int64) (*biz.User, error) {
	key := userCacheKey(strconv.FormatInt(id, 10))
	target, err := r.getUserFromCache(ctx, key)
	if kerrors.IsNotFound(err) {
		user := &User{}
		err = r.data.db.WithContext(ctx).Where("id = ?", id).First(user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, kerrors.NotFound("user not found from db", fmt.Sprintf("user_id(%v)", id))
		}
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("db query system error: user_id(%v)", id))
		}
		target = user
		r.setUserToCache(ctx, user, key)
	}
	if err != nil {
		return nil, err
	}
	return &biz.User{Phone: target.Phone, Email: target.Email, Wechat: target.Wechat, Github: target.Github}, nil
}

func (r *userRepo) GetProfile(ctx context.Context, uuid string) (*biz.Profile, error) {
	key := "profile_" + uuid
	target, err := r.getProfileFromCache(ctx, key)
	if kerrors.IsNotFound(err) {
		profile := &Profile{}
		err = r.data.DB(ctx).WithContext(ctx).Where("uuid = ?", uuid).First(profile).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, kerrors.NotFound("profile not found from db", fmt.Sprintf("uuid(%v)", uuid))
		}
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("db query system error: uuid(%v)", uuid))
		}
		target = profile
		r.setProfileToCache(ctx, profile, key)
	}
	if err != nil {
		return nil, err
	}
	return &biz.Profile{
		Uuid:      target.Uuid,
		Username:  target.Username,
		Avatar:    target.Avatar,
		School:    target.School,
		Company:   target.Company,
		Homepage:  target.Homepage,
		Introduce: target.Introduce,
	}, nil
}

func (r *userRepo) GetUserProfileUpdate(ctx context.Context, uuid string) (*biz.ProfileUpdate, error) {
	profile := &ProfileUpdate{}
	pu := &biz.ProfileUpdate{}
	err := r.data.DB(ctx).WithContext(ctx).Where("uuid = ?", uuid).First(profile).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, kerrors.NotFound("profile update not found from db", fmt.Sprintf("uuid(%v)", uuid))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("db query system error: uuid(%v)", uuid))
	}
	pu.Username = profile.Username
	pu.Avatar = profile.Avatar
	pu.School = profile.School
	pu.Company = profile.Company
	pu.Homepage = profile.Homepage
	pu.Introduce = profile.Introduce
	pu.Status = profile.Status
	return pu, nil
}

func (r *userRepo) SetUserProfile(ctx context.Context, profile *biz.ProfileUpdate) (*biz.ProfileUpdate, error) {
	pu := &ProfileUpdate{}
	pu.Username = profile.Username
	pu.School = profile.School
	pu.Company = profile.Company
	pu.Homepage = profile.Homepage
	pu.Introduce = profile.Introduce
	pu.Status = 2
	err := r.data.DB(ctx).Model(&ProfileUpdate{}).Where("uuid = ?", profile.Uuid).Updates(pu).Error
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "Duplicate") {
			return nil, kerrors.Conflict("username conflict", fmt.Sprintf("profile(%v)", profile))
		} else {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to update a profile: profile(%v)", profile))
		}
	}
	profile.UpdatedAt = pu.UpdatedAt
	return profile, nil
}

func (r *userRepo) SetProfileUpdateRetry(profile *biz.ProfileUpdate) {
	pur := &ProfileUpdateRetry{}
	pur.UpdatedAt = profile.UpdatedAt
	pur.Uuid = profile.Uuid
	pur.Username = profile.Username
	pur.School = profile.School
	pur.Company = profile.Company
	pur.Homepage = profile.Homepage
	pur.Introduce = profile.Introduce
	err := r.data.db.Model(&ProfileUpdateRetry{}).Where("uuid = ?", profile.Uuid).Updates(pur).Error
	if err != nil {
		r.log.Errorf("fail to save profile to retry table, error: %v", err)
	}
}

func (r *userRepo) SendProfileToMq(_ context.Context, profile *biz.ProfileUpdate) error {
	data, err := json.Marshal(profile)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "profile",
		Body:  data,
	}
	msg.WithKeys([]string{profile.Uuid})
	err = r.data.profileMqPro.producer.SendAsync(context.Background(), func(ctx context.Context, result *primitive.SendResult, e error) {
		if e != nil {
			r.log.Errorf("mq receive message error: %v", e)
			r.SetProfileUpdateRetry(profile)
		} else {
			fmt.Printf(result.String())
		}
	}, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to use async producer: %v", err))
	}
	return nil
}

func (r *userRepo) getProfileFromCache(ctx context.Context, key string) (*Profile, error) {
	result, err := r.data.redisCli.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, kerrors.NotFound("profile not found from cache", fmt.Sprintf("key(%s)", key))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get profile from cache: redis.Get(%v)", key))
	}
	var cacheProfile = &Profile{}
	err = json.Unmarshal([]byte(result), cacheProfile)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("json unmarshal error: profile(%v)", result))
	}
	return cacheProfile, nil
}

func (r *userRepo) setProfileToCache(ctx context.Context, profile *Profile, key string) {
	marshal, err := json.Marshal(profile)
	if err != nil {
		r.log.Errorf("fail to set user profile to json: json.Marshal(%v), error(%v)", profile, err)
	}
	err = r.data.redisCli.Set(ctx, key, string(marshal), time.Minute*30).Err()
	if err != nil {
		r.log.Errorf("fail to set user profile to cache: redis.Set(%v), error(%v)", profile, err)
	}
}

// SetUserPhone set redis ?
func (r *userRepo) SetUserPhone(ctx context.Context, id int64, phone string) error {
	err := r.data.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("phone", phone).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("db query system error: user_id(%v), phone(%s)", id, phone))
	}
	return nil
}

// SetUserEmail set redis ?
func (r *userRepo) SetUserEmail(ctx context.Context, id int64, email string) error {
	err := r.data.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("email", email).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("db query system error: user_id(%v), email(%s)", id, email))
	}
	return nil
}

func (r *userRepo) getUserFromCache(ctx context.Context, key string) (*User, error) {
	result, err := r.data.redisCli.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, kerrors.NotFound("user not found from cache", fmt.Sprintf("key(%s)", key))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user from cache: redis.Set(%v)", key))
	}
	var cacheUser = &User{}
	err = json.Unmarshal([]byte(result), cacheUser)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("json unmarshal error: user(%v)", result))
	}
	return cacheUser, nil
}

func (r *userRepo) setUserToCache(ctx context.Context, user *User, key string) {
	marshal, err := json.Marshal(user)
	if err != nil {
		r.log.Errorf("json marshal error: json.Marshal(%v), error(%v)", user, err)
	}
	err = r.data.redisCli.Set(ctx, key, string(marshal), time.Minute*30).Err()
	if err != nil {
		r.log.Errorf("fail to set user to cache: redis.Set(%v), error(%v)", user, err)
	}
}
