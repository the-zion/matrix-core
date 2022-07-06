package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "github.com/the-zion/matrix-core/api/bff/interface/v1"
	"github.com/the-zion/matrix-core/app/bff/interface/internal/biz"
)

var ProviderSet = wire.NewSet(NewBffService)

type BffService struct {
	v1.UnimplementedBffServer
	ac  *biz.ArticleUseCase
	cc  *biz.CreationUseCase
	uc  *biz.UserUseCase
	log *log.Helper
}

func NewBffService(uc *biz.UserUseCase, cc *biz.CreationUseCase, ac *biz.ArticleUseCase, logger log.Logger) *BffService {
	return &BffService{
		log: log.NewHelper(log.With(logger, "module", "bff/interface")),
		uc:  uc,
		ac:  ac,
		cc:  cc,
	}
}
