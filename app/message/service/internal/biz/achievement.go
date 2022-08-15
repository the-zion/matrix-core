package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type AchievementRepo interface {
	SetAchievementAgree(ctx context.Context, uuid string) error
	CancelAchievementAgree(ctx context.Context, uuid string) error
	SetAchievementView(ctx context.Context, uuid string) error
	SetAchievementCollect(ctx context.Context, uuid string) error
	CancelAchievementCollect(ctx context.Context, uuid string) error
	SetAchievementFollow(ctx context.Context, follow, followed string) error
	CancelAchievementFollow(ctx context.Context, follow, followed string) error
	AddAchievementScore(ctx context.Context, uuid string, score int32) error
}

type AchievementCase struct {
	repo AchievementRepo
	log  *log.Helper
}

func NewAchievementUseCase(repo AchievementRepo, logger log.Logger) *AchievementCase {
	return &AchievementCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "message/biz/achievementUseCase")),
	}
}

func (r *AchievementCase) SetAchievementAgree(ctx context.Context, uuid string) error {
	return r.repo.SetAchievementAgree(ctx, uuid)
}

func (r *AchievementCase) CancelAchievementAgree(ctx context.Context, uuid string) error {
	return r.repo.CancelAchievementAgree(ctx, uuid)
}

func (r *AchievementCase) SetAchievementView(ctx context.Context, uuid string) error {
	return r.repo.SetAchievementView(ctx, uuid)
}

func (r *AchievementCase) SetAchievementCollect(ctx context.Context, uuid string) error {
	return r.repo.SetAchievementCollect(ctx, uuid)
}

func (r *AchievementCase) CancelAchievementCollect(ctx context.Context, uuid string) error {
	return r.repo.CancelAchievementCollect(ctx, uuid)
}

func (r *AchievementCase) SetAchievementFollow(ctx context.Context, follow, followed string) error {
	return r.repo.SetAchievementFollow(ctx, follow, followed)
}

func (r *AchievementCase) CancelAchievementFollow(ctx context.Context, follow, followed string) error {
	return r.repo.CancelAchievementFollow(ctx, follow, followed)
}

func (r *AchievementCase) AddAchievementScore(ctx context.Context, uuid string, score int32) error {
	return r.repo.AddAchievementScore(ctx, uuid, score)
}
