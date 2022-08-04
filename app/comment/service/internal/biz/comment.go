package biz

import (
	"context"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/comment/service/v1"
)

type CommentRepo interface {
	GetLastCommentDraft(ctx context.Context, uuid string) (*CommentDraft, error)
	CreateCommentDraft(ctx context.Context, uuid string) (int32, error)
	CreateCommentFolder(ctx context.Context, id int32, uuid string) error
	SendComment(ctx context.Context, id int32, uuid string) (*CommentDraft, error)
	SetRecord(ctx context.Context, id int32, uuid, ip string) error
	SendReviewToMq(ctx context.Context, review *CommentReview) error
}

type CommentUseCase struct {
	repo CommentRepo
	tm   Transaction
	log  *log.Helper
}

func NewCommentUseCase(repo CommentRepo, tm Transaction, logger log.Logger) *CommentUseCase {
	return &CommentUseCase{
		repo: repo,
		tm:   tm,
		log:  log.NewHelper(log.With(logger, "module", "comment/biz/CommentUseCase")),
	}
}

func (r *CommentUseCase) GetLastCommentDraft(ctx context.Context, uuid string) (*CommentDraft, error) {
	draft, err := r.repo.GetLastCommentDraft(ctx, uuid)
	if kerrors.IsNotFound(err) {
		return nil, v1.ErrorRecordNotFound("last draft not found: %s", err.Error())
	}
	if err != nil {
		return nil, v1.ErrorGetCommentDraftFailed("get last draft failed: %s", err.Error())
	}
	return draft, nil
}

func (r *CommentUseCase) CreateCommentDraft(ctx context.Context, uuid string) (int32, error) {
	var id int32
	err := r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		id, err = r.repo.CreateCommentDraft(ctx, uuid)
		if err != nil {
			return v1.ErrorCreateDraftFailed("create comment draft failed: %s", err.Error())
		}

		err = r.repo.CreateCommentFolder(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCreateDraftFailed("create comment folder failed: %s", err.Error())
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *CommentUseCase) SendComment(ctx context.Context, id int32, uuid, ip string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		draft, err := r.repo.SendComment(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCreateDraftFailed("send comment failed: %s", err.Error())
		}

		err = r.repo.SetRecord(ctx, id, uuid, ip)
		if err != nil {
			return v1.ErrorSetRecordFailed("set record failed: %s", err.Error())
		}

		err = r.repo.SendReviewToMq(ctx, &CommentReview{
			Uuid: draft.Uuid,
			Id:   draft.Id,
		})
		if err != nil {
			return v1.ErrorCreateDraftFailed("send create review to mq failed: %s", err.Error())
		}
		return nil
	})
}
