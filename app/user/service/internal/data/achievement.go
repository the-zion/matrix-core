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

var _ biz.AchievementRepo = (*achievementRepo)(nil)

var achCacheKey = func(id string) string {
	return "achievement_" + id
}

type Achievement struct {
	gorm.Model
	UserId   int64 `gorm:"uniqueIndex"`
	Follow   int64
	Followed int64
	Agree    int64
	Collect  int64
	View     int64
}

type Follower struct {
	gorm.Model
	Follow   int64 `gorm:"index"`
	Followed int64 `gorm:"index"`
}

type achievementRepo struct {
	data *Data
	log  *log.Helper
}

func NewAchievementRepo(data *Data, logger log.Logger) biz.AchievementRepo {
	return &achievementRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "user/data/achievement")),
	}
}

func (r *achievementRepo) GetAchievement(ctx context.Context, id int64) (*biz.Achievement, error) {
	key := achCacheKey(strconv.FormatInt(id, 10))
	target, err := r.getAchievementFromCache(ctx, key)
	if err != nil {
		achievement := &Achievement{}
		if err = r.data.db.WithContext(ctx).Where("user_id = ?", id).First(achievement).Error; err != nil {
			//r.log.Errorf("fail to get user achievement from db: id(%v) error(%v)", id, err.Error())
			return nil, biz.ErrAchievementNotFound
		}
		target = achievement
		r.setAchievementToCache(ctx, achievement, key)
	}
	return &biz.Achievement{
		Follow:   target.Follow,
		Followed: target.Followed,
		Agree:    target.Agree,
		Collect:  target.Collect,
		View:     target.View,
	}, nil
}

func (r *achievementRepo) getAchievementFromCache(ctx context.Context, key string) (*Achievement, error) {
	result, err := r.data.redisCli.Get(ctx, key).Result()
	if err != nil {
		//r.log.Errorf("fail to get user achievement from cache:redis.Get(achievement, %v) error(%v)", key, err)
		return nil, err
	}
	var cacheAchievement = &Achievement{}
	err = json.Unmarshal([]byte(result), cacheAchievement)
	if err != nil {
		return nil, err
	}
	return cacheAchievement, nil
}

func (r *achievementRepo) setAchievementToCache(ctx context.Context, achievement *Achievement, key string) {
	marshal, err := json.Marshal(achievement)
	if err != nil {
		//r.log.Errorf("fail to set user achievement to json:json.Marshal(%v) error(%v)", achievement, err)
	}
	err = r.data.redisCli.Set(ctx, key, string(marshal), time.Minute*30).Err()
	if err != nil {
		//r.log.Errorf("fail to set user achievement to cache:redis.Set(%v) error(%v)", achievement, err)
	}
}
