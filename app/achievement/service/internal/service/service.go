package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "github.com/the-zion/matrix-core/api/achievement/service/v1"
	"github.com/the-zion/matrix-core/app/achievement/service/internal/biz"
)

var ProviderSet = wire.NewSet(NewAchievementService)

type AchievementService struct {
	v1.UnimplementedAchievementServer
	ac  *biz.AchievementUseCase
	log *log.Helper
}

func NewAchievementService(ac *biz.AchievementUseCase, logger log.Logger) *AchievementService {
	return &AchievementService{
		log: log.NewHelper(log.With(logger, "module", "achievement/service")),
		ac:  ac,
	}
}
