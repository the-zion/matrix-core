package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
)

type CreationRepo interface {
	GetLeaderBoard(ctx context.Context) ([]*LeaderBoard, error)
	GetCollections(ctx context.Context, uuid string, page int32) ([]*Collections, error)
	GetCollectionsByVisitor(ctx context.Context, uuid string, page int32) ([]*Collections, error)
	GetCollectionsCount(ctx context.Context, uuid string) (int32, error)
	GetCollectionsVisitorCount(ctx context.Context, uuid string) (int32, error)
	CreateCollections(ctx context.Context, uuid, name, introduce string, auth int32) error
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
