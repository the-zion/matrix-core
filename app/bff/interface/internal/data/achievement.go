package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	achievementV1 "github.com/the-zion/matrix-core/api/achievement/service/v1"
	"github.com/the-zion/matrix-core/app/bff/interface/internal/biz"
	"golang.org/x/sync/singleflight"
)

var _ biz.AchievementRepo = (*achievementRepo)(nil)

type achievementRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

func NewAchievementRepo(data *Data, logger log.Logger) biz.AchievementRepo {
	return &achievementRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "bff/data/achievement")),
		sg:   &singleflight.Group{},
	}
}

func (r *achievementRepo) GetAchievementList(ctx context.Context, uuids []string) ([]*biz.Achievement, error) {
	reply := make([]*biz.Achievement, 0)
	achievementList, err := r.data.ac.GetAchievementList(ctx, &achievementV1.GetAchievementListReq{
		Uuids: uuids,
	})
	if err != nil {
		return nil, err
	}
	for _, item := range achievementList.Achievement {
		reply = append(reply, &biz.Achievement{
			Uuid:     item.Uuid,
			View:     item.View,
			Agree:    item.Agree,
			Follow:   item.Follow,
			Followed: item.Followed,
		})
	}
	return reply, nil
}

func (r *achievementRepo) GetUserAchievement(ctx context.Context, uuid string) (*biz.Achievement, error) {
	achievement, err := r.data.ac.GetUserAchievement(ctx, &achievementV1.GetUserAchievementReq{
		Uuid: uuid,
	})
	if err != nil {
		return nil, err
	}
	return &biz.Achievement{
		Agree:    achievement.Agree,
		View:     achievement.View,
		Collect:  achievement.Collect,
		Follow:   achievement.Follow,
		Followed: achievement.Followed,
	}, nil
}
