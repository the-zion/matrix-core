package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "github.com/the-zion/matrix-core/api/user/service/v1"
	"github.com/the-zion/matrix-core/app/info/interface/internal/biz"
)

var ProviderSet = wire.NewSet(NewUserService)

type InfoService struct {
	v1.UnimplementedUserServer
	uc  *biz.UserUseCase
	log *log.Helper
}

func NewUserService(uc *biz.UserUseCase, logger log.Logger) *InfoService {
	return &InfoService{
		log: log.NewHelper(log.With(logger, "module", "user/service")),
		uc:  uc,
	}
}
