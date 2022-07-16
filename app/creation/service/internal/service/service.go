package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
	"github.com/the-zion/matrix-core/app/creation/service/internal/biz"
)

var ProviderSet = wire.NewSet(NewCreationService)

type CreationService struct {
	v1.UnimplementedCreationServer
	ac  *biz.ArticleUseCase
	tc  *biz.TalkUseCase
	cc  *biz.CreationUseCase
	log *log.Helper
}

func NewCreationService(ac *biz.ArticleUseCase, tc *biz.TalkUseCase, cc *biz.CreationUseCase, logger log.Logger) *CreationService {
	return &CreationService{
		log: log.NewHelper(log.With(logger, "module", "creation/service")),
		ac:  ac,
		tc:  tc,
		cc:  cc,
	}
}
