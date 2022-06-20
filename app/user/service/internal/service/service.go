package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "github.com/the-zion/matrix-core/api/user/service/v1"
	"github.com/the-zion/matrix-core/app/user/service/internal/biz"
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

func NewUserService(uc *biz.UserUseCase, ac *biz.AuthUseCase, pc *biz.ProfileUseCase, logger log.Logger) *UserService {
	return &UserService{
		log: log.NewHelper(log.With(logger, "module", "user/service")),
		uc:  uc,
		ac:  ac,
		pc:  pc,
	}
}
