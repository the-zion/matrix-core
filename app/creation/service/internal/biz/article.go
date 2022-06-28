package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type ArticleRepo interface {
	CreateArticleDraft(ctx context.Context, uuid string) (int32, error)
}

type ArticleUseCase struct {
	repo ArticleRepo
	tm   Transaction
	log  *log.Helper
}

func NewArticleUseCase(repo ArticleRepo, tm Transaction, logger log.Logger) *ArticleUseCase {
	return &ArticleUseCase{
		repo: repo,
		tm:   tm,
		log:  log.NewHelper(log.With(logger, "module", "creation/biz/articleUseCase")),
	}
}

func (r *ArticleUseCase) CreateArticleDraft(ctx context.Context, uuid string) (int32, error) {
	id, err := r.repo.CreateArticleDraft(ctx, uuid)
	if err != nil {
		return 0, nil
	}
	return id, nil
}
