package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
	"github.com/the-zion/matrix-core/app/creation/service/internal/biz"
	"google.golang.org/protobuf/types/known/emptypb"
)

var ProviderSet = wire.NewSet(NewCreationService)

type CreationService struct {
	v1.UnimplementedCreationServer
	ac  *biz.ArticleUseCase
	tc  *biz.TalkUseCase
	cc  *biz.CreationUseCase
	coc *biz.ColumnUseCase
	nc  *biz.NewsUseCase
	log *log.Helper
}

func NewCreationService(ac *biz.ArticleUseCase, tc *biz.TalkUseCase, cc *biz.CreationUseCase, coc *biz.ColumnUseCase, nc *biz.NewsUseCase, logger log.Logger) *CreationService {
	return &CreationService{
		log: log.NewHelper(log.With(logger, "module", "creation/service")),
		ac:  ac,
		tc:  tc,
		cc:  cc,
		coc: coc,
		nc:  nc,
	}
}

func (s *CreationService) GetHealth(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
