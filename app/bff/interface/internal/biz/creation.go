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
	log  *log.Helper
}

func NewArticleUseCase(repo ArticleRepo, logger log.Logger) *ArticleUseCase {
	return &ArticleUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "bff/biz/ArticleUseCase")),
	}
}

func (r *ArticleUseCase) CreateArticleDraft(ctx context.Context) (int32, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.CreateArticleDraft(ctx, uuid)
}
