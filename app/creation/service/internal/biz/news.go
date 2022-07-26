package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
)

type NewsRepo interface {
	GetNewsFromTianXing(ctx context.Context, page int32, key string) ([]*News, error)
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

func (r *NewsUseCase) GetNewsFromTianXing(ctx context.Context, page int32, kind string) ([]*News, error) {
	news, err := r.repo.GetNewsFromTianXing(ctx, page, kind)
	if err != nil {
		return nil, v1.ErrorGetNewsFailed("get news failed: %s", err.Error())
	}
	return news, nil
}
