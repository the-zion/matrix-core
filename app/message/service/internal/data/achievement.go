package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
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

func (r *achievementRepo) SetAchievementAgree(ctx context.Context, uuid string) error {
	_, err := r.data.ac.SetAchievementAgree(ctx, &achievementv1.SetAchievementAgreeReq{
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *achievementRepo) CancelAchievementAgree(ctx context.Context, uuid string) error {
	_, err := r.data.ac.SetAchievementAgree(ctx, &achievementv1.SetAchievementAgreeReq{
		Uuid: uuid,
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
