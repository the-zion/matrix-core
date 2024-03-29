package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "github.com/the-zion/matrix-core/api/bff/interface/v1"
	"github.com/the-zion/matrix-core/app/bff/interface/internal/biz"
	"google.golang.org/protobuf/types/known/emptypb"
)

var ProviderSet = wire.NewSet(NewBffService)

type BffService struct {
	v1.UnimplementedBffServer
	ac    *biz.ArticleUseCase
	tc    *biz.TalkUseCase
	cc    *biz.CreationUseCase
	coc   *biz.ColumnUseCase
	uc    *biz.UserUseCase
	achc  *biz.AchievementUseCase
	nc    *biz.NewsUseCase
	commc *biz.CommentUseCase
	mc    *biz.MessageUseCase
	log   *log.Helper
}

func NewBffService(uc *biz.UserUseCase, cc *biz.CreationUseCase, tc *biz.TalkUseCase, ac *biz.ArticleUseCase, coc *biz.ColumnUseCase, achc *biz.AchievementUseCase, nc *biz.NewsUseCase, commc *biz.CommentUseCase, mc *biz.MessageUseCase, logger log.Logger) *BffService {
	return &BffService{
		log:   log.NewHelper(log.With(logger, "module", "bff/interface")),
		uc:    uc,
		ac:    ac,
		tc:    tc,
		cc:    cc,
		coc:   coc,
		achc:  achc,
		nc:    nc,
		mc:    mc,
		commc: commc,
	}
}

func (s *BffService) GetHealth(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
