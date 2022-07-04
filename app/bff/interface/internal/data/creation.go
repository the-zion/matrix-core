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

func (r *articleRepo) GetLastArticleDraft(ctx context.Context, uuid string) (*biz.ArticleDraft, error) {
	reply, err := r.data.ac.GetLastArticleDraft(ctx, &creationV1.GetLastArticleDraftReq{
		Uuid: uuid,
	})
	if err != nil {
		return nil, err
	}
	return &biz.ArticleDraft{
		Id:     reply.Id,
		Status: reply.Status,
	}, nil
}
func (r *articleRepo) GetArticleList(ctx context.Context, page int32) ([]*biz.Article, error) {
	reply := make([]*biz.Article, 0)
	articleList, err := r.data.ac.GetArticleList(ctx, &creationV1.GetArticleListReq{
		Page: page,
	})
	if err != nil {
		return reply, err
	}
	for _, item := range articleList.Article {
		reply = append(reply, &biz.Article{
			Id:   item.Id,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (r *articleRepo) GetArticleDraftList(ctx context.Context, uuid string) ([]*biz.ArticleDraft, error) {
	reply := make([]*biz.ArticleDraft, 0)
	draftList, err := r.data.ac.GetArticleDraftList(ctx, &creationV1.GetArticleDraftListReq{
		Uuid: uuid,
	})
	if err != nil {
		return reply, err
	}
	for _, item := range draftList.Draft {
		reply = append(reply, &biz.ArticleDraft{
			Id: item.Id,
		})
	}
	return reply, nil
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

func (r *articleRepo) ArticleDraftMark(ctx context.Context, uuid string, id int32) error {
	_, err := r.data.ac.ArticleDraftMark(ctx, &creationV1.ArticleDraftMarkReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *articleRepo) SendArticle(ctx context.Context, uuid string, id int32) error {
	_, err := r.data.ac.SendArticle(ctx, &creationV1.SendArticleReq{
		Uuid: uuid,
		Id:   id,
	})
	if err != nil {
		return err
	}
	return nil
}
