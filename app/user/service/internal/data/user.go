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

func (r *userRepo) GetAccount(ctx context.Context, uuid string) (*biz.User, error) {
	user := &User{}
	err := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, kerrors.NotFound("user not found from db", fmt.Sprintf("uuid(%v)", uuid))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("db query system error: uuid(%v)", uuid))
	}
	return &biz.User{
		Phone:  user.Phone,
		Email:  user.Email,
		Qq:     user.Qq,
		Wechat: user.Wechat,
		Weibo:  user.Weibo,
		Github: user.Github,
	}, nil
}

func (r *userRepo) GetProfile(ctx context.Context, uuid string) (*biz.Profile, error) {
	key := "profile_" + uuid
	target, err := r.getProfileFromCache(ctx, key)
	if kerrors.IsNotFound(err) {
		profile := &Profile{}
		err = r.data.db.WithContext(ctx).Where("uuid = ?", uuid).First(profile).Error
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
		Updated:   strconv.FormatInt(target.Updated, 10),
		Uuid:      target.Uuid,
		Username:  target.Username,
		Avatar:    target.Avatar,
		School:    target.School,
		Company:   target.Company,
		Job:       target.Job,
		Homepage:  target.Homepage,
		Introduce: target.Introduce,
	}, nil
}

func (r *userRepo) GetProfileUpdate(ctx context.Context, uuid string) (*biz.ProfileUpdate, error) {
	profile := &ProfileUpdate{}
	pu := &biz.ProfileUpdate{}
	err := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).First(profile).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, kerrors.NotFound("profile update not found from db", fmt.Sprintf("uuid(%v)", uuid))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("db query system error: uuid(%v)", uuid))
	}
	pu.Updated = strconv.FormatInt(profile.Updated, 10)
	pu.Uuid = profile.Uuid
	pu.Username = profile.Username
	pu.Avatar = profile.Avatar
	pu.School = profile.School
	pu.Company = profile.Company
	pu.Job = profile.Job
	pu.Homepage = profile.Homepage
	pu.Introduce = profile.Introduce
	pu.Status = profile.Status
	return pu, nil
}

func (r *userRepo) SetProfile(ctx context.Context, profile *biz.ProfileUpdate) error {
	updateTime, err := strconv.ParseInt(profile.Updated, 10, 64)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to convert string to int64, update: %v", profile.Updated))
	}
	p := &Profile{}
	p.Updated = updateTime
	p.Username = profile.Username
	p.School = profile.School
	p.Company = profile.Company
	p.Job = profile.Job
	p.Homepage = profile.Homepage
	p.Introduce = profile.Introduce
	err = r.data.DB(ctx).Model(&Profile{}).Where("uuid = ?", profile.Uuid).Updates(p).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set profile: profile(%v)", profile))
	}
	r.setProfileToCache(ctx, p, "profile_"+profile.Uuid)
	return nil
}

func (r *userRepo) SetProfileUpdate(ctx context.Context, profile *biz.ProfileUpdate, status int32) (*biz.ProfileUpdate, error) {
	pu := &ProfileUpdate{}
	pu.Updated = time.Now().Unix()
	pu.Username = profile.Username
	pu.School = profile.School
	pu.Company = profile.Company
	pu.Job = profile.Job
	pu.Homepage = profile.Homepage
	pu.Introduce = profile.Introduce
	pu.Status = status
	err := r.data.DB(ctx).Model(&ProfileUpdate{}).Where("uuid = ?", profile.Uuid).Updates(pu).Error
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "Duplicate") {
			return nil, kerrors.Conflict("username conflict", fmt.Sprintf("profile(%v)", profile))
		} else {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to update a profile: profile(%v)", profile))
		}
	}
	profile.Updated = strconv.FormatInt(pu.Updated, 10)
	return profile, nil
}

func (r *userRepo) SetProfileUpdateRetry(profile *biz.ProfileUpdate) {
	updateTime, err := strconv.ParseInt(profile.Updated, 10, 64)
	if err != nil {
		r.log.Errorf("fail to transform string to int64, error: %v", err)
		return
	}
	pur := &ProfileUpdateRetry{}
	pur.Updated = updateTime
	pur.Uuid = profile.Uuid
	pur.Username = profile.Username
	pur.School = profile.School
	pur.Company = profile.Company
	pur.Job = profile.Job
	pur.Homepage = profile.Homepage
	pur.Introduce = profile.Introduce
	err = r.data.db.Model(&ProfileUpdateRetry{}).Where("uuid = ?", profile.Uuid).Updates(pur).Error
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

func (r *userRepo) ModifyProfileUpdateStatus(ctx context.Context, uuid, update string) error {
	updateTime, err := strconv.ParseInt(update, 10, 64)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to convert string to int64, update: %v", update))
	}
	err = r.data.DB(ctx).Model(&ProfileUpdate{}).Where("uuid = ? and updated = ?", uuid, updateTime).Update("status", 1).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to modify profile update status: uuid(%v)", uuid))
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

// SetUserEmail set redis ?
func (r *userRepo) SetUserEmail(ctx context.Context, id int64, email string) error {
	err := r.data.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("email", email).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("db query system error: user_id(%v), email(%s)", id, email))
	}
	return nil
}
