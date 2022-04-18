package data

import (
	"context"
	"encoding/json"
	"github.com/Cube-v2/cube-core/app/user/service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
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
	if err != nil {
		profile := &Profile{}
		if err = r.data.db.WithContext(ctx).Where("user_id = ?", id).First(profile).Error; err != nil {
			r.log.Errorf("fail to get user profile from db: id(%v) error(%v)", id, err.Error())
			return nil, biz.ErrProfileNotFound
		}
		target = profile
		r.setProfileToCache(ctx, profile, key)
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

func (r *profileRepo) getProfileFromCache(ctx context.Context, key string) (*Profile, error) {
	result, err := r.data.redisCli.Get(ctx, key).Result()
	if err != nil {
		r.log.Errorf("fail to get user profile from cache:redis.Get(profile, %v) error(%v)", key, err)
		return nil, err
	}
	var cacheProfile = &Profile{}
	err = json.Unmarshal([]byte(result), cacheProfile)
	if err != nil {
		return nil, err
	}
	return cacheProfile, nil
}

func (r *profileRepo) setProfileToCache(ctx context.Context, profile *Profile, key string) {
	marshal, err := json.Marshal(profile)
	if err != nil {
		r.log.Errorf("fail to set user profile to json:json.Marshal(%v) error(%v)", profile, err)
	}
	err = r.data.redisCli.Set(ctx, key, string(marshal), time.Minute*30).Err()
	if err != nil {
		r.log.Errorf("fail to set user profile to cache:redis.Set(%v) error(%v)", profile, err)
	}
}
