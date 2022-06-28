package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	creationV1 "github.com/the-zion/matrix-core/api/creation/service/v1"
	"github.com/the-zion/matrix-core/app/bff/interface/internal/biz"
	"golang.org/x/sync/singleflight"
)

var _ biz.ArticleRepo = (*articleRepo)(nil)

type articleRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

func NewArticleRepo(data *Data, logger log.Logger) biz.ArticleRepo {
	return &articleRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "bff/data/article")),
		sg:   &singleflight.Group{},
	}
}

func (r *articleRepo) CreateArticleDraft(ctx context.Context, uuid string) (int32, error) {
	reply, err := r.data.ac.CreateArticleDraft(ctx, &creationV1.CreateArticleDraftReq{
		Uuid: uuid,
	})
	if err != nil {
		return 0, err
	}
	return reply.Id, nil
}
