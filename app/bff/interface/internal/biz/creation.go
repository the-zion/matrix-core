package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type CreationRepo interface {
	GetLeaderBoard(ctx context.Context) ([]*LeaderBoard, error)
	GetCollectArticle(ctx context.Context, id, page int32) ([]*Article, error)
	GetCollectArticleCount(ctx context.Context, id int32) (int32, error)
	CreateCollections(ctx context.Context, uuid, name, introduce string, auth int32) error
	EditCollections(ctx context.Context, id int32, uuid, name, introduce string, auth int32) error
	DeleteCollections(ctx context.Context, id int32, uuid string) error
	GetCollection(ctx context.Context, id int32, uuid string) (*Collections, error)
	GetCollections(ctx context.Context, uuid string, page int32) ([]*Collections, error)
	GetCollectionsByVisitor(ctx context.Context, uuid string, page int32) ([]*Collections, error)
	GetCollectionsCount(ctx context.Context, uuid string) (int32, error)
	GetCollectionsVisitorCount(ctx context.Context, uuid string) (int32, error)
}

type ArticleRepo interface {
	CreateArticleDraft(ctx context.Context, uuid string) (int32, error)
	GetLastArticleDraft(ctx context.Context, uuid string) (*ArticleDraft, error)
	GetArticleList(ctx context.Context, page int32) ([]*Article, error)
	GetArticleListHot(ctx context.Context, page int32) ([]*Article, error)
	GetUserArticleList(ctx context.Context, page int32, uuid string) ([]*Article, error)
	GetUserArticleListVisitor(ctx context.Context, page int32, uuid string) ([]*Article, error)
	GetArticleStatistic(ctx context.Context, id int32) (*ArticleStatistic, error)
	GetArticleListStatistic(ctx context.Context, ids []int32) ([]*ArticleStatistic, error)
	GetArticleDraftList(ctx context.Context, uuid string) ([]*ArticleDraft, error)
	ArticleDraftMark(ctx context.Context, id int32, uuid string) error
	SendArticle(ctx context.Context, id int32, uuid string) error
	SendArticleEdit(ctx context.Context, id int32, uuid string) error
	SetArticleAgree(ctx context.Context, id int32, uuid, userUuid string) error
	SetArticleView(ctx context.Context, id int32, uuid string) error
	SetArticleCollect(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error
	CancelArticleAgree(ctx context.Context, id int32, uuid, userUuid string) error
	CancelArticleCollect(ctx context.Context, id int32, uuid, userUuid string) error
	ArticleStatisticJudge(ctx context.Context, id int32, uuid string) (*ArticleStatisticJudge, error)
}

type CreationUseCase struct {
	repo CreationRepo
	log  *log.Helper
}

type ArticleUseCase struct {
	repo ArticleRepo
	log  *log.Helper
}

func NewCreationUseCase(repo CreationRepo, logger log.Logger) *CreationUseCase {
	return &CreationUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "bff/biz/CreationUseCase")),
	}
}

func NewArticleUseCase(repo ArticleRepo, logger log.Logger) *ArticleUseCase {
	return &ArticleUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "bff/biz/ArticleUseCase")),
	}
}

func (r *CreationUseCase) GetLeaderBoard(ctx context.Context) ([]*LeaderBoard, error) {
	return r.repo.GetLeaderBoard(ctx)
}

func (r *CreationUseCase) GetCollectArticle(ctx context.Context, id, page int32) ([]*Article, error) {
	return r.repo.GetCollectArticle(ctx, id, page)
}

func (r *CreationUseCase) GetCollectArticleCount(ctx context.Context, id int32) (int32, error) {
	return r.repo.GetCollectArticleCount(ctx, id)
}

func (r *CreationUseCase) GetCollection(ctx context.Context, id int32, uuid string) (*Collections, error) {
	return r.repo.GetCollection(ctx, id, uuid)
}

func (r *CreationUseCase) GetCollections(ctx context.Context, page int32) ([]*Collections, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetCollections(ctx, uuid, page)
}

func (r *CreationUseCase) GetCollectionsCount(ctx context.Context) (int32, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetCollectionsCount(ctx, uuid)
}

func (r *CreationUseCase) GetCollectionsByVisitor(ctx context.Context, page int32, uuid string) ([]*Collections, error) {
	return r.repo.GetCollectionsByVisitor(ctx, uuid, page)
}

func (r *CreationUseCase) GetCollectionsVisitorCount(ctx context.Context, uuid string) (int32, error) {
	return r.repo.GetCollectionsVisitorCount(ctx, uuid)
}

func (r *CreationUseCase) CreateCollections(ctx context.Context, name, introduce string, auth int32) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.CreateCollections(ctx, uuid, name, introduce, auth)
}

func (r *CreationUseCase) EditCollections(ctx context.Context, id int32, name, introduce string, auth int32) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.EditCollections(ctx, id, uuid, name, introduce, auth)
}

func (r *CreationUseCase) DeleteCollections(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.DeleteCollections(ctx, id, uuid)
}

func (r *ArticleUseCase) GetArticleList(ctx context.Context, page int32) ([]*Article, error) {
	return r.repo.GetArticleList(ctx, page)
}

func (r *ArticleUseCase) GetArticleListHot(ctx context.Context, page int32) ([]*Article, error) {
	return r.repo.GetArticleListHot(ctx, page)
}

func (r *ArticleUseCase) GetUserArticleList(ctx context.Context, page int32) ([]*Article, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetUserArticleList(ctx, page, uuid)
}

func (r *ArticleUseCase) GetUserArticleListVisitor(ctx context.Context, page int32, uuid string) ([]*Article, error) {
	return r.repo.GetUserArticleListVisitor(ctx, page, uuid)
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

func (r *ArticleUseCase) SendArticleEdit(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.SendArticleEdit(ctx, id, uuid)
}

func (r *ArticleUseCase) SetArticleAgree(ctx context.Context, id int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.SetArticleAgree(ctx, id, uuid, userUuid)
}

func (r *ArticleUseCase) SetArticleView(ctx context.Context, id int32, uuid string) error {
	return r.repo.SetArticleView(ctx, id, uuid)
}

func (r *ArticleUseCase) SetArticleCollect(ctx context.Context, id, collectionsId int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.SetArticleCollect(ctx, id, collectionsId, uuid, userUuid)
}

func (r *ArticleUseCase) CancelArticleAgree(ctx context.Context, id int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.CancelArticleAgree(ctx, id, uuid, userUuid)
}

func (r *ArticleUseCase) CancelArticleCollect(ctx context.Context, id int32, uuid string) error {
	userUuid := ctx.Value("uuid").(string)
	return r.repo.CancelArticleCollect(ctx, id, uuid, userUuid)
}

func (r *ArticleUseCase) ArticleStatisticJudge(ctx context.Context, id int32) (*ArticleStatisticJudge, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.ArticleStatisticJudge(ctx, id, uuid)
}
