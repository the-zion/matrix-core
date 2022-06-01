package biz

import (
	"context"
	"errors"
	v1 "github.com/Cube-v2/matrix-core/api/user/service/v1"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrAchievementNotFound = errors.New("achievement not found")
)

type Achievement struct {
	Follow   int64
	Followed int64
	Agree    int64
	Collect  int64
	View     int64
}

type AchievementRepo interface {
	GetAchievement(ctx context.Context, id int64) (*Achievement, error)
}

type AchievementUseCase struct {
	repo AchievementRepo
	log  *log.Helper
}

func NewAchievementUseCase(repo AchievementRepo, logger log.Logger) *AchievementUseCase {
	return &AchievementUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "user/biz/AchievementUseCase")),
	}
}

func (r *AchievementUseCase) GetUserAchievement(ctx context.Context, id int64) (*Achievement, error) {
	ach, err := r.repo.GetAchievement(ctx, id)
	if err != nil {
		return nil, v1.ErrorGetAchievementFailed("get user achievement failed: %s", err.Error())
	}
	return &Achievement{
		Follow:   ach.Follow,
		Followed: ach.Followed,
		Agree:    ach.Agree,
		Collect:  ach.Collect,
		View:     ach.View,
	}, nil
}
