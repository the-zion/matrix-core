package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "github.com/the-zion/matrix-core/api/bff/interface/v1"
	"github.com/the-zion/matrix-core/app/bff/interface/internal/biz"
)

var ProviderSet = wire.NewSet(NewUserService)

type BffService struct {
	v1.UnimplementedBffServer
	uc  *biz.UserUseCase
	log *log.Helper
}

func NewUserService(uc *biz.UserUseCase, logger log.Logger) *BffService {
	return &BffService{
		log: log.NewHelper(log.With(logger, "module", "user/service")),
		uc:  uc,
	}
}