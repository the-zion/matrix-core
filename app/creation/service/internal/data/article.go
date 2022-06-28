package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/the-zion/matrix-core/app/creation/service/internal/biz"
)

var _ biz.ArticleRepo = (*articleRepo)(nil)

type articleRepo struct {
	data *Data
	log  *log.Helper
}

func NewArticleRepo(data *Data, logger log.Logger) biz.ArticleRepo {
	return &articleRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "creation/data/article")),
	}
}

func (r *articleRepo) CreateArticleDraft(ctx context.Context, uuid string) (int32, error) {
	return 0, nil
}
