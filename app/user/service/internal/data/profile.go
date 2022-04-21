package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Cube-v2/cube-core/app/user/service/internal/biz"
	v2 "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
	"time"
)

var _ biz.ProfileRepo = (*profileRepo)(nil)

var profileCacheKey = func(id string) string {
	return "profile_" + id
}

type profileRepo struct {
	data *Data
	log  *log.Helper
}

func NewProfileRepo(data *Data, logger log.Logger) biz.ProfileRepo {
	return &profileRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "user/data/profile")),
	}
}

type Profile struct {
	gorm.Model
	UserId          int64  `gorm:"uniqueIndex"`
	Username        string `gorm:"uniqueIndex;size:200"`
	Sex             string `gorm:"size:100"`
	Introduce       string `gorm:"size:200"`
	Address         string `gorm:"size:100"`
	Industry        string `gorm:"size:100"`
	PersonalProfile string `gorm:"size:500"`
	Tag             string `gorm:"size:100"`
	Background      string `gorm:"size:500"`
	Image           string `gorm:"size:500"`
}

func (r *profileRepo) GetProfile(ctx context.Context, id int64) (*biz.Profile, error) {
	key := profileCacheKey(strconv.FormatInt(id, 10))
	target, err := r.getProfileFromCache(ctx, key)
	if v2.IsNotFound(err) {
		profile := &Profile{}
		err = r.data.db.WithContext(ctx).Where("user_id = ?", id).First(profile).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v2.NotFound("profile not found from db", fmt.Sprintf("user_id(%v)", id))
		}
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("db query system error: user_id(%v)", id))
		}
		target = profile
		r.setProfileToCache(ctx, profile, key)
	}
	if err != nil {
		return nil, err
	}
	return &biz.Profile{
		Username:        target.Username,
		Sex:             target.Sex,
		Introduce:       target.Introduce,
		Address:         target.Address,
		Industry:        target.Industry,
		PersonalProfile: target.PersonalProfile,
		Tag:             target.Tag,
		Background:      target.Background,
		Image:           target.Image,
	}, nil
}

// SetProfile set redis ?
func (r *profileRepo) SetProfile(ctx context.Context, id int64, sex, introduce, industry, address, profile, tag string) error {
	p := Profile{
		Sex:             sex,
		Introduce:       introduce,
		Industry:        industry,
		Address:         address,
		PersonalProfile: profile,
		Tag:             tag,
	}
	err := r.data.db.WithContext(ctx).Model(&Profile{}).Where("user_id = ?", id).Select("Sex", "Introduce", "Industry", "Address", "PersonalProfile", "Tag").Updates(p).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set user profile to db: profile(%v), user_id(%v)", p, id))
	}
	return nil
}

// SetName set redis ?
func (r *profileRepo) SetName(ctx context.Context, id int64, name string) error {
	err := r.data.db.WithContext(ctx).Model(&Profile{}).Where("user_id = ?", id).Update("username", name).Error
	if err != nil {
		//r.log.Errorf("fail to set user name to db:name(%v) error(%v)", name, err)
		return errors.Wrapf(err, fmt.Sprintf("fail to set user name to db: name(%s), user_id(%v)", name, id))
	}
	return nil
}

func (r *profileRepo) getProfileFromCache(ctx context.Context, key string) (*Profile, error) {
	result, err := r.data.redisCli.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, v2.NotFound("profile not found from cache", fmt.Sprintf("key(%s)", key))
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

func (r *profileRepo) setProfileToCache(ctx context.Context, profile *Profile, key string) {
	marshal, err := json.Marshal(profile)
	if err != nil {
		r.log.Errorf("fail to set user profile to json: json.Marshal(%v), error(%v)", profile, err)
	}
	err = r.data.redisCli.Set(ctx, key, string(marshal), time.Minute*30).Err()
	if err != nil {
		r.log.Errorf("fail to set user profile to cache: redis.Set(%v), error(%v)", profile, err)
	}
}
