package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
)

type CreationRepo interface {
	GetLeaderBoard(ctx context.Context) ([]*LeaderBoard, error)
	GetCollectArticle(ctx context.Context, id, page int32) ([]*Article, error)
	GetCollectArticleCount(ctx context.Context, id int32) (int32, error)
	GetCollectCount(ctx context.Context, id int32) (int64, error)
	GetCollection(ctx context.Context, id int32, uuid string) (*Collections, error)
	GetCollections(ctx context.Context, uuid string, page int32) ([]*Collections, error)
	GetCollectionsByVisitor(ctx context.Context, uuid string, page int32) ([]*Collections, error)
	GetCollectionsCount(ctx context.Context, uuid string) (int32, error)
	GetCollectionsVisitorCount(ctx context.Context, uuid string) (int32, error)
	CreateCollections(ctx context.Context, uuid, name, introduce string, auth int32) error
	EditCollections(ctx context.Context, id int32, uuid, name, introduce string, auth int32) error
	DeleteCollections(ctx context.Context, id int32, uuid string) error
}

type CreationUseCase struct {
	repo CreationRepo
	log  *log.Helper
}

func NewCreationUseCase(repo CreationRepo, logger log.Logger) *CreationUseCase {
	return &CreationUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "creation/biz/creationUseCase")),
	}
}

func (r *CreationUseCase) GetLeaderBoard(ctx context.Context) ([]*LeaderBoard, error) {
	boardList, err := r.repo.GetLeaderBoard(ctx)
	if err != nil {
		return nil, v1.ErrorGetLeaderBoardFailed("get leader board failed: %s", err.Error())
	}
	return boardList, nil
}

func (r *CreationUseCase) GetCollectArticle(ctx context.Context, id, page int32) ([]*Article, error) {
	articleList, err := r.repo.GetCollectArticle(ctx, id, page)
	if err != nil {
		return nil, v1.ErrorGetCollectArticleFailed("get collect article list failed: %s", err.Error())
	}
	return articleList, nil
}

func (r *CreationUseCase) GetCollectArticleCount(ctx context.Context, id int32) (int32, error) {
	count, err := r.repo.GetCollectArticleCount(ctx, id)
	if err != nil {
		return 0, v1.ErrorGetCollectArticleCountFailed("get collect article count failed: %s", err.Error())
	}
	return count, nil
}

func (r *CreationUseCase) GetCollection(ctx context.Context, id int32, uuid string) (*Collections, error) {
	collection, err := r.repo.GetCollection(ctx, id, uuid)
	if err != nil {
		return nil, v1.ErrorGetCollectionFailed("get collection failed: %s", err.Error())
	}
	return collection, nil
}

func (r *CreationUseCase) GetCollections(ctx context.Context, uuid string, page int32) ([]*Collections, error) {
	collections, err := r.repo.GetCollections(ctx, uuid, page)
	if err != nil {
		return nil, v1.ErrorGetCollectionsFailed("get collections failed: %s", err.Error())
	}
	return collections, nil
}
func (r *CreationUseCase) GetCollectionsByVisitor(ctx context.Context, uuid string, page int32) ([]*Collections, error) {
	collections, err := r.repo.GetCollectionsByVisitor(ctx, uuid, page)
	if err != nil {
		return nil, v1.ErrorGetCollectionsFailed("get collections failed: %s", err.Error())
	}
	return collections, nil
}

func (r *CreationUseCase) GetCollectionsCount(ctx context.Context, uuid string) (int32, error) {
	count, err := r.repo.GetCollectionsCount(ctx, uuid)
	if err != nil {
		return 0, v1.ErrorGetCollectionsCountFailed("get collections count failed: %s", err.Error())
	}
	return count, nil
}

func (r *CreationUseCase) GetCollectionsVisitorCount(ctx context.Context, uuid string) (int32, error) {
	count, err := r.repo.GetCollectionsVisitorCount(ctx, uuid)
	if err != nil {
		return 0, v1.ErrorGetCollectionsCountFailed("get collections count failed: %s", err.Error())
	}
	return count, nil
}

func (r *CreationUseCase) CreateCollections(ctx context.Context, uuid, name, introduce string, auth int32) error {
	err := r.repo.CreateCollections(ctx, uuid, name, introduce, auth)
	if err != nil {
		return v1.ErrorCreateCollectionsFailed("create collections failed: %s", err.Error())
	}
	return nil
}

func (r *CreationUseCase) EditCollections(ctx context.Context, id int32, uuid, name, introduce string, auth int32) error {
	err := r.repo.EditCollections(ctx, id, uuid, name, introduce, auth)
	if err != nil {
		return v1.ErrorEditCollectionsFailed("edit collections failed: %s", err.Error())
	}
	return nil
}

func (r *CreationUseCase) DeleteCollections(ctx context.Context, id int32, uuid string) error {
	count, err := r.repo.GetCollectCount(ctx, id)
	if err != nil {
		return v1.ErrorGetCollectCountFailed("get collection count failed: %s", err.Error())
	}

	if count == 1 {
		return v1.ErrorCollectionsIsNotEmpty("collections is not empty")
	}

	err = r.repo.DeleteCollections(ctx, id, uuid)
	if err != nil {
		return v1.ErrorDeleteCollectionsFailed("delete collections failed: %s", err.Error())
	}
	return nil
}
