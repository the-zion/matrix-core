package data

import (
	"context"
	"fmt"
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
	result, err, _ := r.sg.Do(fmt.Sprintf("user_achievement_%s", uuid), func() (interface{}, error) {
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
			Score:    achievement.Score,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.Achievement), nil
}

func (r *achievementRepo) GetUserMedal(ctx context.Context, uuid string) (*biz.Medal, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("user_medal_%s", uuid), func() (interface{}, error) {
		medal, err := r.data.ac.GetUserMedal(ctx, &achievementV1.GetUserMedalReq{
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		return &biz.Medal{
			Creation1: medal.Creation1,
			Creation2: medal.Creation2,
			Creation3: medal.Creation3,
			Creation4: medal.Creation4,
			Creation5: medal.Creation5,
			Creation6: medal.Creation6,
			Creation7: medal.Creation7,
			Agree1:    medal.Agree1,
			Agree2:    medal.Agree2,
			Agree3:    medal.Agree3,
			Agree4:    medal.Agree4,
			Agree5:    medal.Agree5,
			Agree6:    medal.Agree6,
			View1:     medal.View1,
			View2:     medal.View2,
			View3:     medal.View3,
			Comment1:  medal.Comment1,
			Comment2:  medal.Comment2,
			Comment3:  medal.Comment3,
			Collect1:  medal.Collect1,
			Collect2:  medal.Collect2,
			Collect3:  medal.Collect3,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.Medal), nil
}

func (r *achievementRepo) GetUserActive(ctx context.Context, uuid string) (*biz.Active, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("user_active_%s", uuid), func() (interface{}, error) {
		active, err := r.data.ac.GetUserActive(ctx, &achievementV1.GetUserActiveReq{
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		return &biz.Active{
			Agree: active.Agree,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.Active), nil
}

func (r *achievementRepo) SetUserMedal(ctx context.Context, medal, uuid string) error {
	_, err := r.data.ac.SetUserMedal(ctx, &achievementV1.SetUserMedalReq{
		Uuid:  uuid,
		Medal: medal,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *achievementRepo) CancelUserMedalSet(ctx context.Context, medal, uuid string) error {
	_, err := r.data.ac.CancelUserMedalSet(ctx, &achievementV1.CancelUserMedalSetReq{
		Uuid:  uuid,
		Medal: medal,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *achievementRepo) AccessUserMedal(ctx context.Context, medal, uuid string) error {
	_, err := r.data.ac.AccessUserMedal(ctx, &achievementV1.AccessUserMedalReq{
		Uuid:  uuid,
		Medal: medal,
	})
	if err != nil {
		return err
	}
	return nil
}
