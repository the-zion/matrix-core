package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type CreationRepo interface {
	GetLeaderBoard(ctx context.Context) ([]*LeaderBoard, error)
	CreateCollections(ctx context.Context, uuid, name, introduce string, auth int32) error
}

type ArticleRepo interface {
	CreateArticleDraft(ctx context.Context, uuid string) (int32, error)
	GetLastArticleDraft(ctx context.Context, uuid string) (*ArticleDraft, error)
	GetArticleList(ctx context.Context, page int32) ([]*Article, error)
	GetArticleListHot(ctx context.Context, page int32) ([]*Article, error)
	GetArticleStatistic(ctx context.Context, id int32) (*ArticleStatistic, error)
	GetArticleListStatistic(ctx context.Context, ids []int32) ([]*ArticleStatistic, error)
	GetArticleDraftList(ctx context.Context, uuid string) ([]*ArticleDraft, error)
	ArticleDraftMark(ctx context.Context, id int32, uuid string) error
	SendArticle(ctx context.Context, id int32, uuid string) error
	SetArticleAgree(ctx context.Context, id int32, uuid string) error
	SetArticleView(ctx context.Context, id int32, uuid string) error
	SetArticleCollect(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error
}

type ArticleUseCase struct {
	repo ArticleRepo
	log  *log.Helper
}

type CreationUseCase struct {
	repo CreationRepo
	log  *log.Helper
}

func NewArticleUseCase(repo ArticleRepo, logger log.Logger) *ArticleUseCase {
	return &ArticleUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "bff/biz/ArticleUseCase")),
	}
}

func NewCreationUseCase(repo CreationRepo, logger log.Logger) *CreationUseCase {
	return &CreationUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "bff/biz/CreationUseCase")),
	}
}

func (r *CreationUseCase) GetLeaderBoard(ctx context.Context) ([]*LeaderBoard, error) {
	return r.repo.GetLeaderBoard(ctx)
}

func (r *ArticleUseCase) GetArticleList(ctx context.Context, page int32) ([]*Article, error) {
	return r.repo.GetArticleList(ctx, page)
}

func (r *ArticleUseCase) GetArticleListHot(ctx context.Context, page int32) ([]*Article, error) {
	return r.repo.GetArticleListHot(ctx, page)
}

func (r *ArticleUseCase) GetArticleStatistic(ctx context.Context, id int32) (*ArticleStatistic, error) {
	return r.repo.GetArticleStatistic(ctx, id)
}

func (r *ArticleUseCase) GetArticleListStatistic(ctx context.Context, ids []int32) ([]*ArticleStatistic, error) {
	return r.repo.GetArticleListStatistic(ctx, ids)
}

func (r *ArticleUseCase) GetLastArticleDraft(ctx context.Context) (*ArticleDraft, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetLastArticleDraft(ctx, uuid)
}

func (r *ArticleUseCase) CreateArticleDraft(ctx context.Context) (int32, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.CreateArticleDraft(ctx, uuid)
}

func (r *ArticleUseCase) ArticleDraftMark(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.ArticleDraftMark(ctx, id, uuid)
}

func (r *ArticleUseCase) GetArticleDraftList(ctx context.Context) ([]*ArticleDraft, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetArticleDraftList(ctx, uuid)
}

func (r *ArticleUseCase) SendArticle(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.SendArticle(ctx, id, uuid)
}

func (r *ArticleUseCase) SetArticleAgree(ctx context.Context, id int32, uuid string) error {
	return r.repo.SetArticleAgree(ctx, id, uuid)
}

func (r *ArticleUseCase) SetArticleView(ctx context.Context, id int32, uuid string) error {
	return r.repo.SetArticleView(ctx, id, uuid)
}

func (r *ArticleUseCase) SetArticleCollect(ctx context.Context, id, collectionsId int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.SetArticleCollect(ctx, id, collectionsId, uuid, userUuid)
}

func (r *CreationUseCase) CreateCollections(ctx context.Context, name, introduce string, auth int32) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.CreateCollections(ctx, uuid, name, introduce, auth)
}
