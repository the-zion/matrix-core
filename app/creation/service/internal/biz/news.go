package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
)

type NewsRepo interface {
	GetNews(ctx context.Context, page int32) ([]*News, error)
	GetNewsSearch(ctx context.Context, page int32, search, time string) ([]*NewsSearch, int32, error)
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

func (r *NewsUseCase) GetNewsSearch(ctx context.Context, page int32, search, time string) ([]*NewsSearch, int32, error) {
	newsList, total, err := r.repo.GetNewsSearch(ctx, page, search, time)
	if err != nil {
		return nil, 0, v1.ErrorGetNewsSearchFailed("get news search failed: %s", err.Error())
	}
	return newsList, total, nil
}
