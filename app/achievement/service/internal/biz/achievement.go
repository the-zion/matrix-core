package biz

import (
	"github.com/go-kratos/kratos/v2/log"
)

type AchievementRepo interface {
}

type AchievementUseCase struct {
	repo AchievementRepo
	tm   Transaction
	log  *log.Helper
}

func NewAchievementUseCase(repo AchievementRepo, tm Transaction, logger log.Logger) *AchievementUseCase {
	return &AchievementUseCase{
		repo: repo,
		tm:   tm,
		log:  log.NewHelper(log.With(logger, "module", "achievement/biz/AchievementUseCase")),
	}
}
