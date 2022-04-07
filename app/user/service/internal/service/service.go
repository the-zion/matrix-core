package service

import (
	"github.com/Cube-v2/cube-core/api/user/service/v1"
	"github.com/Cube-v2/cube-core/app/user/service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewUserService)

type UserService struct {
	v1.UnimplementedUserServer
	uc  *biz.UserUseCase
	log *log.Helper
}

func NewUserService(
	uc *biz.UserUseCase,
	logger log.Logger) *UserService {
	return &UserService{
		log: log.NewHelper(log.With(logger, "module", "service/user")),
		uc:  uc,
	}
}
