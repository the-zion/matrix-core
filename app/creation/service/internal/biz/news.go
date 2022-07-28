package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
)

type NewsRepo interface {
	GetNews(ctx context.Context, page int32) ([]*News, error)
}

type NewsUseCase struct {
	repo NewsRepo
	tm   Transaction
	log  *log.Helper
}

func NewNewsUseCase(repo NewsRepo, tm Transaction, logger log.Logger) *NewsUseCase {
	return &NewsUseCase{
		repo: repo,
		tm:   tm,
		log:  log.NewHelper(log.With(logger, "module", "creation/biz/newsUseCase")),
	}
}

func (r *NewsUseCase) GetNews(ctx context.Context, page int32) ([]*News, error) {
	news, err := r.repo.GetNews(ctx, page)
	if err != nil {
		return nil, v1.ErrorGetNewsFailed("get news failed: %s", err.Error())
	}
	return news, nil
}
