package biz

import (
	"context"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/comment/service/v1"
	"golang.org/x/sync/errgroup"
)

type CommentRepo interface {
	GetLastCommentDraft(ctx context.Context, uuid string) (*CommentDraft, error)
	GetUserCommentAgree(ctx context.Context, uuid string) (map[int32]bool, error)
	GetCommentList(ctx context.Context, page, creationId, creationType int32) ([]*Comment, error)
	GetCommentListHot(ctx context.Context, page, creationId, creationType int32) ([]*Comment, error)
	GetCommentListStatistic(ctx context.Context, ids []int32) ([]*CommentStatistic, error)
	CreateCommentDraft(ctx context.Context, uuid string) (int32, error)
	CreateCommentFolder(ctx context.Context, id int32, uuid string) error
	CreateComment(ctx context.Context, id, creationId, creationType int32, uuid string) error
	CreateCommentCache(ctx context.Context, id, creationId, creationType int32, uuid string) error
	SendComment(ctx context.Context, id int32, uuid string) (*CommentDraft, error)
	SendCommentAgreeToMq(ctx context.Context, id, creationId, creationType int32, uuid, userUuid, mode string) error
	SendCommentToMq(ctx context.Context, comment *Comment, mode string) error
	SendReviewToMq(ctx context.Context, review *CommentReview) error
	SendCommentStatisticToMq(ctx context.Context, uuid, mode string) error
	SetRecord(ctx context.Context, id int32, uuid, ip string) error
	SetCommentAgree(ctx context.Context, id int32, uuid string) error
	SetUserCommentAgree(ctx context.Context, id int32, userUuid string) error
	SetUserCommentAgreeToCache(ctx context.Context, id int32, userUuid string) error
	SetCommentAgreeToCache(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error
	CancelUserCommentAgreeFromCache(ctx context.Context, id int32, userUuid string) error
	CancelCommentAgree(ctx context.Context, id int32, uuid string) error
	CancelUserCommentAgree(ctx context.Context, id int32, userUuid string) error
	CancelCommentAgreeFromCache(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error
	DeleteCommentDraft(ctx context.Context, id int32, uuid string) error
	RemoveComment(ctx context.Context, id int32, uuid string) error
	RemoveCommentAgree(ctx context.Context, id int32, uuid string) error
	//RemoveSubComment(ctx context.Context, id int32, uuid string) error
	RemoveCommentCache(ctx context.Context, id, creationId, creationType int32, uuid string) error
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

func (r *CommentUseCase) GetUserCommentAgree(ctx context.Context, uuid string) (map[int32]bool, error) {
	agreeMap, err := r.repo.GetUserCommentAgree(ctx, uuid)
	if err != nil {
		return nil, v1.ErrorGetUserCommentAgreeFailed("get user comment agree failed: %s", err.Error())
	}
	return agreeMap, nil
}

func (r *CommentUseCase) GetCommentList(ctx context.Context, page, creationId, creationType int32) ([]*Comment, error) {
	commentList, err := r.repo.GetCommentList(ctx, page, creationId, creationType)
	if err != nil {
		return nil, v1.ErrorGetCommentListFailed("get comment list failed: %s", err.Error())
	}
	return commentList, nil
}

func (r *CommentUseCase) GetCommentListHot(ctx context.Context, page, creationId, creationType int32) ([]*Comment, error) {
	commentList, err := r.repo.GetCommentListHot(ctx, page, creationId, creationType)
	if err != nil {
		return nil, v1.ErrorGetCommentListFailed("get comment hot list failed: %s", err.Error())
	}
	return commentList, nil
}

func (r *CommentUseCase) GetCommentListStatistic(ctx context.Context, ids []int32) ([]*CommentStatistic, error) {
	commentListStatistic, err := r.repo.GetCommentListStatistic(ctx, ids)
	if err != nil {
		return nil, v1.ErrorGetCommentStatisticFailed("get comment statistic list failed: %s", err.Error())
	}
	return commentListStatistic, nil
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

func (r *CommentUseCase) CreateComment(ctx context.Context, id, creationTd, creationType int32, uuid string) error {
	err := r.repo.SendCommentToMq(ctx, &Comment{
		CommentId:    id,
		Uuid:         uuid,
		CreationId:   creationTd,
		CreationType: creationType,
	}, "create_comment_db_cache_and_search")
	if err != nil {
		return v1.ErrorCreateCommentFailed("create comment to mq failed: %s", err.Error())
	}
	return nil
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

func (r *CommentUseCase) RemoveComment(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error {
	if uuid != userUuid {
		return v1.ErrorRemoveCommentFailed("remove comment failed: no auth")
	}
	err := r.repo.RemoveCommentCache(ctx, id, creationId, creationType, uuid)
	if err != nil {
		return v1.ErrorRemoveCommentFailed("remove comment failed: %s", err.Error())
	}

	err = r.repo.SendCommentToMq(ctx, &Comment{
		CommentId:    id,
		CreationId:   creationId,
		CreationType: creationType,
		Uuid:         uuid,
	}, "remove_comment_db_and_cache")
	if err != nil {
		return v1.ErrorCancelAgreeFailed("remove comment failed failed: %s", err.Error())
	}
	return nil
}

func (r *CommentUseCase) SetCommentAgree(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error {
	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		err := r.repo.SetUserCommentAgreeToCache(ctx, id, userUuid)
		if err != nil {
			return v1.ErrorSetAgreeFailed("set comment agree to cache failed: %s", err.Error())
		}
		return nil
	})
	g.Go(func() error {
		err := r.repo.SendCommentAgreeToMq(ctx, id, creationId, creationType, uuid, userUuid, "set_comment_agree_db_and_cache")
		if err != nil {
			return v1.ErrorSetAgreeFailed("set comment agree to mq failed: %s", err.Error())
		}
		return nil
	})
	return g.Wait()
}

func (r *CommentUseCase) SetCommentAgreeDbAndCache(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetCommentAgree(ctx, id, uuid)
		if err != nil {
			return v1.ErrorSetAgreeFailed("set comment agree failed: %s", err.Error())
		}
		err = r.repo.SetUserCommentAgree(ctx, id, userUuid)
		if err != nil {
			return v1.ErrorSetAgreeFailed("set comment agree failed: %s", err.Error())
		}
		err = r.repo.SetCommentAgreeToCache(ctx, id, creationId, creationType, uuid, userUuid)
		if err != nil {
			return v1.ErrorSetAgreeFailed("set comment agree to cache failed: %s", err.Error())
		}
		err = r.repo.SendCommentStatisticToMq(ctx, uuid, "agree")
		if err != nil {
			return v1.ErrorSetAgreeFailed("set comment agree to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *CommentUseCase) CancelCommentAgree(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error {
	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		err := r.repo.CancelUserCommentAgreeFromCache(ctx, id, userUuid)
		if err != nil {
			return v1.ErrorCancelAgreeFailed("cancel comment agree from cache failed: %s", err.Error())
		}
		return nil
	})
	g.Go(func() error {
		err := r.repo.SendCommentAgreeToMq(ctx, id, creationId, creationType, uuid, userUuid, "cancel_comment_agree_db_and_cache")
		if err != nil {
			return v1.ErrorCancelAgreeFailed("cancel comment agree from mq failed: %s", err.Error())
		}
		return nil
	})
	return g.Wait()
}

func (r *CommentUseCase) CancelCommentAgreeDbAndCache(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.CancelCommentAgree(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCancelAgreeFailed("cancel comment agree failed: %s", err.Error())
		}
		err = r.repo.CancelUserCommentAgree(ctx, id, userUuid)
		if err != nil {
			return v1.ErrorCancelAgreeFailed("cancel comment agree failed: %s", err.Error())
		}
		err = r.repo.CancelCommentAgreeFromCache(ctx, id, creationId, creationType, uuid, userUuid)
		if err != nil {
			return v1.ErrorCancelAgreeFailed("cancel comment agree from cache failed: %s", err.Error())
		}
		err = r.repo.SendCommentStatisticToMq(ctx, uuid, "agree_cancel")
		if err != nil {
			return v1.ErrorCancelAgreeFailed("cancel comment agree from mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *CommentUseCase) CreateCommentDbAndCache(ctx context.Context, id, creationId, creationType int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		err = r.repo.DeleteCommentDraft(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCreateCommentFailed("delete comment draft failed: %s", err.Error())
		}

		err = r.repo.CreateComment(ctx, id, creationId, creationType, uuid)
		if err != nil {
			return v1.ErrorCreateCommentFailed("create comment failed: %s", err.Error())
		}

		err = r.repo.CreateCommentCache(ctx, id, creationId, creationType, uuid)
		if err != nil {
			return v1.ErrorCreateCommentFailed("create comment cache failed: %s", err.Error())
		}
		return nil
	})
}

func (r *CommentUseCase) RemoveCommentDbAndCache(ctx context.Context, id, creationId, creationType int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		err = r.repo.RemoveComment(ctx, id, uuid)
		if err != nil {
			return v1.ErrorRemoveCommentFailed("remove comment failed: %s", err.Error())
		}

		err = r.repo.RemoveCommentAgree(ctx, id, uuid)
		if err != nil {
			return v1.ErrorRemoveCommentFailed("remove comment agree record failed: %s", err.Error())
		}

		err = r.repo.RemoveCommentCache(ctx, id, creationId, creationType, uuid)
		if err != nil {
			return v1.ErrorRemoveCommentFailed("remove comment cache failed: %s", err.Error())
		}
		return nil
	})
}
