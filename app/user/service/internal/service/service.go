package service

import (
	v1 "github.com/Cube-v2/cube-core/api/user/service/v1"
	"github.com/Cube-v2/cube-core/app/user/service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewUserService)

type UserService struct {
	v1.UnimplementedUserServer
	uc   *biz.UserUseCase
	ac   *biz.AuthUseCase
	pc   *biz.ProfileUseCase
	achc *biz.AchievementUseCase
	log  *log.Helper
}

func NewUserService(uc *biz.UserUseCase, ac *biz.AuthUseCase, pc *biz.ProfileUseCase, achc *biz.AchievementUseCase, logger log.Logger) *UserService {
	return &UserService{
		log:  log.NewHelper(log.With(logger, "module", "user/service/user")),
		uc:   uc,
		ac:   ac,
		pc:   pc,
		achc: achc,
	}
}
