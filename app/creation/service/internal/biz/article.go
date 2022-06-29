package biz

import (
	"context"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
)

type ArticleRepo interface {
	GetLastArticleDraft(ctx context.Context, uuid string) (*ArticleDraft, error)
	CreateArticleDraft(ctx context.Context, uuid string) (int32, error)
	CreateArticleFolder(ctx context.Context, id int32) error
	ArticleDraftMark(ctx context.Context, uuid string, id int32) error
	GetArticleDraftList(ctx context.Context, uuid string) ([]*ArticleDraft, error)
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

func (r *ArticleUseCase) GetLastArticleDraft(ctx context.Context, uuid string) (*ArticleDraft, error) {
	draft, err := r.repo.GetLastArticleDraft(ctx, uuid)
	if kerrors.IsNotFound(err) {
		return nil, v1.ErrorRecordNotFound("last draft not found: %s", err.Error())
	}
	if err != nil {
		return nil, v1.ErrorGetArticleDraftFailed("get last draft failed: %s", err.Error())
	}
	return draft, nil
}

func (r *ArticleUseCase) CreateArticleDraft(ctx context.Context, uuid string) (int32, error) {
	var id int32
	err := r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		id, err = r.repo.CreateArticleDraft(ctx, uuid)
		if err != nil {
			return v1.ErrorCreateArticleDraftFailed("create article draft failed: %s", err.Error())
		}

		err = r.repo.CreateArticleFolder(ctx, id)
		if err != nil {
			return v1.ErrorCreateArticleFolderFailed("create article folder failed: %s", err.Error())
		}

		return nil
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ArticleUseCase) ArticleDraftMark(ctx context.Context, uuid string, id int32) error {
	err := r.repo.ArticleDraftMark(ctx, uuid, id)
	if err != nil {
		return v1.ErrorDraftMarkFailed("mark draft failed: %s", err.Error())
	}
	return nil
}

func (r *ArticleUseCase) GetArticleDraftList(ctx context.Context, uuid string) ([]*ArticleDraft, error) {
	draftList, err := r.repo.GetArticleDraftList(ctx, uuid)
	if err != nil {
		return draftList, v1.ErrorGetDraftListFailed("get article draft list failed: %s", err.Error())
	}
	return draftList, nil
}
