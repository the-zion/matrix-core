package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	achievementV1 "github.com/the-zion/matrix-core/api/achievement/service/v1"
	achievementv1 "github.com/the-zion/matrix-core/api/achievement/service/v1"
	"github.com/the-zion/matrix-core/app/message/service/internal/biz"
)

type achievementRepo struct {
	data *Data
	log  *log.Helper
}

func NewAchievementRepo(data *Data, logger log.Logger) biz.AchievementRepo {
	return &achievementRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "message/data/achievement")),
	}
}

func (r *achievementRepo) SetAchievementAgree(ctx context.Context, uuid, userUuid string) error {
	_, err := r.data.ac.SetAchievementAgree(ctx, &achievementv1.SetAchievementAgreeReq{
		Uuid:     uuid,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *achievementRepo) CancelAchievementAgree(ctx context.Context, uuid, userUuid string) error {
	_, err := r.data.ac.CancelAchievementAgree(ctx, &achievementv1.CancelAchievementAgreeReq{
		Uuid:     uuid,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *achievementRepo) SetAchievementView(ctx context.Context, uuid string) error {
	_, err := r.data.ac.SetAchievementView(ctx, &achievementv1.SetAchievementViewReq{
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *achievementRepo) SetAchievementCollect(ctx context.Context, uuid string) error {
	_, err := r.data.ac.SetAchievementCollect(ctx, &achievementv1.SetAchievementCollectReq{
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *achievementRepo) CancelAchievementCollect(ctx context.Context, uuid string) error {
	_, err := r.data.ac.CancelAchievementCollect(ctx, &achievementv1.CancelAchievementCollectReq{
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *achievementRepo) SetAchievementFollow(ctx context.Context, follow, followed string) error {
	_, err := r.data.ac.SetAchievementFollow(ctx, &achievementv1.SetAchievementFollowReq{
		Follow:   follow,
		Followed: followed,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *achievementRepo) SetUserMedalDbAndCache(ctx context.Context, medal, uuid string) error {
	_, err := r.data.ac.SetUserMedalDbAndCache(ctx, &achievementv1.SetUserMedalDbAndCacheReq{
		Medal: medal,
		Uuid:  uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *achievementRepo) CancelAchievementFollow(ctx context.Context, follow, followed string) error {
	_, err := r.data.ac.CancelAchievementFollow(ctx, &achievementv1.CancelAchievementFollowReq{
		Follow:   follow,
		Followed: followed,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *achievementRepo) CancelUserMedalDbAndCache(ctx context.Context, medal, uuid string) error {
	_, err := r.data.ac.CancelUserMedalDbAndCache(ctx, &achievementv1.CancelUserMedalDbAndCacheReq{
		Medal: medal,
		Uuid:  uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *achievementRepo) AddAchievementScore(ctx context.Context, uuid string, score int32) error {
	_, err := r.data.ac.AddAchievementScore(ctx, &achievementv1.AddAchievementScoreReq{
		Uuid:  uuid,
		Score: score,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *achievementRepo) GetUserAchievement(ctx context.Context, uuid string) (int32, int32, error) {
	achievement, err := r.data.ac.GetUserAchievement(ctx, &achievementV1.GetUserAchievementReq{
		Uuid: uuid,
	})
	if err != nil {
		return 0, 0, err
	}
	return achievement.Agree, achievement.View, nil
}

func (r *achievementRepo) GetUserActive(ctx context.Context, uuid string) (int32, error) {
	active, err := r.data.ac.GetUserActive(ctx, &achievementV1.GetUserActiveReq{
		Uuid: uuid,
	})
	if err != nil {
		return 0, err
	}
	return active.Agree, nil
}

func (r *achievementRepo) AccessUserMedalDbAndCache(ctx context.Context, medal, uuid string) error {
	_, err := r.data.ac.AccessUserMedalDbAndCache(ctx, &achievementv1.AccessUserMedalReq{
		Medal: medal,
		Uuid:  uuid,
	})
	if err != nil {
		return err
	}
	return nil
}
