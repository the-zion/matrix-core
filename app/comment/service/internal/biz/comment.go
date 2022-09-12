package biz

import (
	"context"
	"fmt"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	v1 "github.com/the-zion/matrix-core/api/comment/service/v1"
	"golang.org/x/sync/errgroup"
)

type CommentRepo interface {
	GetLastCommentDraft(ctx context.Context, uuid string) (*CommentDraft, error)
	GetUserCommentAgree(ctx context.Context, uuid string) (map[int32]bool, error)
	GetCommentUser(ctx context.Context, uuid string) (*CommentUser, error)
	GetCommentList(ctx context.Context, page, creationId, creationType int32) ([]*Comment, error)
	GetSubCommentList(ctx context.Context, page, id int32) ([]*SubComment, error)
	GetCommentListHot(ctx context.Context, page, creationId, creationType int32) ([]*Comment, error)
	GetCommentListStatistic(ctx context.Context, ids []int32) ([]*CommentStatistic, error)
	GetSubCommentListStatistic(ctx context.Context, ids []int32) ([]*CommentStatistic, error)
	GetUserCommentArticleReplyList(ctx context.Context, page int32, uuid string) ([]*Comment, error)
	GetUserSubCommentArticleReplyList(ctx context.Context, page int32, uuid string) ([]*SubComment, error)
	GetUserCommentTalkReplyList(ctx context.Context, page int32, uuid string) ([]*Comment, error)
	GetUserSubCommentTalkReplyList(ctx context.Context, page int32, uuid string) ([]*SubComment, error)
	GetUserCommentArticleRepliedList(ctx context.Context, page int32, uuid string) ([]*Comment, error)
	GetUserSubCommentArticleRepliedList(ctx context.Context, page int32, uuid string) ([]*SubComment, error)
	GetUserCommentTalkRepliedList(ctx context.Context, page int32, uuid string) ([]*Comment, error)
	GetUserSubCommentTalkRepliedList(ctx context.Context, page int32, uuid string) ([]*SubComment, error)
	GetCommentContentReview(ctx context.Context, page int32, uuid string) ([]*TextReview, error)
	GetRootComment(ctx context.Context, id int32) (*Comment, error)
	GetSubComment(ctx context.Context, id int32) (*SubComment, error)
	GetParentCommentUserId(ctx context.Context, id, rootId int32) (string, error)
	GetSubCommentReply(ctx context.Context, id int32, uuid string) (string, error)
	GetArticleAuthor(ctx context.Context, id int32) (string, error)
	GetTalkAuthor(ctx context.Context, id int32) (string, error)

	CreateCommentDraft(ctx context.Context, uuid string) (int32, error)
	CreateCommentFolder(ctx context.Context, id int32, uuid string) error
	CreateComment(ctx context.Context, id, creationId, creationType int32, creationAuthor, uuid string) error
	CreateSubComment(ctx context.Context, id, rootId, parentId, creationId, creationType int32, creationAuthor, rootUser, uuid, reply string) error
	CreateCommentCache(ctx context.Context, id, creationId, creationType int32, creationAuthor, uuid string) error
	CreateSubCommentCache(ctx context.Context, id, rootId, parentId, creationId, creationType int32, creationAuthor, rootUser, uuid, reply string) error

	AddArticleCommentUser(ctx context.Context, creationAuthor, uuid string) error
	AddArticleSubCommentUser(ctx context.Context, parentId int32, rootUser, parentUser, uuid string) error
	AddTalkCommentUser(ctx context.Context, creationAuthor, uuid string) error
	AddTalkSubCommentUser(ctx context.Context, parentId int32, rootUser, parentUser, uuid string) error

	SendComment(ctx context.Context, id int32, uuid string) (*CommentDraft, error)
	SendCommentAgreeToMq(ctx context.Context, id, creationId, creationType int32, uuid, userUuid, mode string) error
	SendSubCommentAgreeToMq(ctx context.Context, id int32, uuid, userUuid, mode string) error
	SendCommentToMq(ctx context.Context, comment *Comment, mode string) error
	SendCommentContentIrregularToMq(ctx context.Context, review *TextReview) error
	SendSubCommentToMq(ctx context.Context, comment *SubComment, mode string) error
	SendReviewToMq(ctx context.Context, review *CommentReview) error
	SendCommentStatisticToMq(ctx context.Context, uuid, userUuid, mode string) error
	SendScoreToMq(ctx context.Context, score int32, uuid, mode string) error

	SetRecord(ctx context.Context, id int32, uuid, ip string) error
	SetCommentAgree(ctx context.Context, id int32, uuid string) error
	SetSubCommentAgree(ctx context.Context, id int32, uuid string) error
	SetCommentComment(ctx context.Context, id int32) error
	SetUserCommentAgree(ctx context.Context, id int32, userUuid string) error
	SetUserCommentAgreeToCache(ctx context.Context, id int32, userUuid string) error
	SetCommentAgreeToCache(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error
	SetSubCommentAgreeToCache(ctx context.Context, id int32, uuid, userUuid string) error
	SetCommentContentIrregular(ctx context.Context, review *TextReview) (*TextReview, error)
	SetCommentContentIrregularToCache(ctx context.Context, review *TextReview) error

	CancelUserCommentAgreeFromCache(ctx context.Context, id int32, userUuid string) error
	CancelCommentAgree(ctx context.Context, id int32, uuid string) error
	CancelSubCommentAgree(ctx context.Context, id int32, uuid string) error
	CancelUserCommentAgree(ctx context.Context, id int32, userUuid string) error
	CancelCommentAgreeFromCache(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error
	CancelSubCommentAgreeFromCache(ctx context.Context, id int32, uuid, userUuid string) error

	DeleteCommentDraft(ctx context.Context, id int32, uuid string) error

	ReduceArticleCommentUser(ctx context.Context, creationAuthor, uuid string) error
	ReduceTalkCommentUser(ctx context.Context, creationAuthor, uuid string) error
	ReduceCreationComment(ctx context.Context, createId, createType int32, uuid string) error
	ReduceArticleSubCommentUser(ctx context.Context, parentId int32, rootUser, parentUser, uuid string) error
	ReduceTalkSubCommentUser(ctx context.Context, parentId int32, rootUser, parentUser, uuid string) error

	RemoveComment(ctx context.Context, id int32, uuid string) error
	RemoveSubComment(ctx context.Context, id, rootId int32, uuid string) error
	RemoveCommentAgree(ctx context.Context, id int32, uuid string) error
	RemoveCommentCache(ctx context.Context, id, creationId, creationType int32, creationAuthor, uuid string) error
	RemoveSubCommentCache(ctx context.Context, id, rootId, parentId, creationId, creationType int32, creationAuthor, rootUser, uuid, reply string) error
}

type CommentUseCase struct {
	repo CommentRepo
	tm   Transaction
	re   Recovery
	log  *log.Helper
}

func NewCommentUseCase(repo CommentRepo, re Recovery, tm Transaction, logger log.Logger) *CommentUseCase {
	return &CommentUseCase{
		repo: repo,
		tm:   tm,
		re:   re,
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

func (r *CommentUseCase) GetCommentUser(ctx context.Context, uuid string) (*CommentUser, error) {
	commentUser, err := r.repo.GetCommentUser(ctx, uuid)
	if err != nil {
		return nil, v1.ErrorGetCommentUserFailed("get comment user failed: %s", err.Error())
	}
	return commentUser, nil
}

func (r *CommentUseCase) GetCommentList(ctx context.Context, page, creationId, creationType int32) ([]*Comment, error) {
	commentList, err := r.repo.GetCommentList(ctx, page, creationId, creationType)
	if err != nil {
		return nil, v1.ErrorGetCommentListFailed("get comment list failed: %s", err.Error())
	}
	return commentList, nil
}

func (r *CommentUseCase) GetSubCommentList(ctx context.Context, page, id int32) ([]*SubComment, error) {
	subCommentList, err := r.repo.GetSubCommentList(ctx, page, id)
	if err != nil {
		return nil, v1.ErrorGetCommentListFailed("get sub comment list failed: %s", err.Error())
	}
	return subCommentList, nil
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

func (r *CommentUseCase) GetSubCommentListStatistic(ctx context.Context, ids []int32) ([]*CommentStatistic, error) {
	commentListStatistic, err := r.repo.GetSubCommentListStatistic(ctx, ids)
	if err != nil {
		return nil, v1.ErrorGetCommentStatisticFailed("get sub comment statistic list failed: %s", err.Error())
	}
	return commentListStatistic, nil
}

func (r *CommentUseCase) GetUserCommentArticleReplyList(ctx context.Context, page int32, uuid string) ([]*Comment, error) {
	commentList, err := r.repo.GetUserCommentArticleReplyList(ctx, page, uuid)
	if err != nil {
		return nil, v1.ErrorGetUserCommentCreationReplyListFailed("get user comment article reply list failed: %s", err.Error())
	}
	return commentList, nil
}

func (r *CommentUseCase) GetUserSubCommentArticleReplyList(ctx context.Context, page int32, uuid string) ([]*SubComment, error) {
	commentList, err := r.repo.GetUserSubCommentArticleReplyList(ctx, page, uuid)
	if err != nil {
		return nil, v1.ErrorGetUserSubCommentCreationReplyListFailed("get user sub comment article reply list failed: %s", err.Error())
	}
	return commentList, nil
}

func (r *CommentUseCase) GetUserCommentTalkReplyList(ctx context.Context, page int32, uuid string) ([]*Comment, error) {
	commentList, err := r.repo.GetUserCommentTalkReplyList(ctx, page, uuid)
	if err != nil {
		return nil, v1.ErrorGetUserCommentCreationReplyListFailed("get user comment talk reply list failed: %s", err.Error())
	}
	return commentList, nil
}

func (r *CommentUseCase) GetUserSubCommentTalkReplyList(ctx context.Context, page int32, uuid string) ([]*SubComment, error) {
	commentList, err := r.repo.GetUserSubCommentTalkReplyList(ctx, page, uuid)
	if err != nil {
		return nil, v1.ErrorGetUserSubCommentCreationReplyListFailed("get user sub comment talk reply list failed: %s", err.Error())
	}
	return commentList, nil
}

func (r *CommentUseCase) GetUserCommentArticleRepliedList(ctx context.Context, page int32, uuid string) ([]*Comment, error) {
	commentList, err := r.repo.GetUserCommentArticleRepliedList(ctx, page, uuid)
	if err != nil {
		return nil, v1.ErrorGetUserCommentCreationRepliedListFailed("get user comment article replied list failed: %s", err.Error())
	}
	return commentList, nil
}

func (r *CommentUseCase) GetUserSubCommentArticleRepliedList(ctx context.Context, page int32, uuid string) ([]*SubComment, error) {
	commentList, err := r.repo.GetUserSubCommentArticleRepliedList(ctx, page, uuid)
	if err != nil {
		return nil, v1.ErrorGetUserSubCommentCreationRepliedListFailed("get user sub comment article replied list failed: %s", err.Error())
	}
	return commentList, nil
}

func (r *CommentUseCase) GetUserCommentTalkRepliedList(ctx context.Context, page int32, uuid string) ([]*Comment, error) {
	commentList, err := r.repo.GetUserCommentTalkRepliedList(ctx, page, uuid)
	if err != nil {
		return nil, v1.ErrorGetUserCommentCreationRepliedListFailed("get user comment talk replied list failed: %s", err.Error())
	}
	return commentList, nil
}

func (r *CommentUseCase) GetUserSubCommentTalkRepliedList(ctx context.Context, page int32, uuid string) ([]*SubComment, error) {
	commentList, err := r.repo.GetUserSubCommentTalkRepliedList(ctx, page, uuid)
	if err != nil {
		return nil, v1.ErrorGetUserSubCommentCreationRepliedListFailed("get user sub comment talk replied list failed: %s", err.Error())
	}
	return commentList, nil
}

func (r *CommentUseCase) GetCommentContentReview(ctx context.Context, page int32, uuid string) ([]*TextReview, error) {
	reviewList, err := r.repo.GetCommentContentReview(ctx, page, uuid)
	if err != nil {
		return nil, v1.ErrorGetContentReviewFailed("get comment content review failed: %s", err.Error())
	}
	return reviewList, nil
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
	}, "create_comment_db_and_cache")
	if err != nil {
		return v1.ErrorCreateCommentFailed("create comment to mq failed: %s", err.Error())
	}
	return nil
}

func (r *CommentUseCase) CommentContentIrregular(ctx context.Context, review *TextReview) error {
	err := r.repo.SendCommentContentIrregularToMq(ctx, review)
	if err != nil {
		return v1.ErrorSetContentIrregularFailed("set comment content irregular to mq failed: %s", err.Error())
	}
	return nil
}

func (r *CommentUseCase) CreateSubComment(ctx context.Context, id, rootId, parentId int32, uuid string) error {
	err := r.repo.SendSubCommentToMq(ctx, &SubComment{
		CommentId: id,
		Uuid:      uuid,
		RootId:    rootId,
		ParentId:  parentId,
	}, "create_sub_comment_db_and_cache")
	if err != nil {
		return v1.ErrorCreateCommentFailed("create sub comment to mq failed: %s", err.Error())
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
			Mode: "comment_review",
		})
		if err != nil {
			return v1.ErrorCreateDraftFailed("send create review to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *CommentUseCase) SendSubComment(ctx context.Context, id int32, uuid, ip string) error {
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
			Mode: "sub_comment_review",
		})
		if err != nil {
			return v1.ErrorCreateDraftFailed("send create review to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *CommentUseCase) RemoveComment(ctx context.Context, id int32, uuid string) error {
	err := r.repo.SendCommentToMq(ctx, &Comment{
		CommentId: id,
		Uuid:      uuid,
	}, "remove_comment_db_and_cache")
	if err != nil {
		return v1.ErrorCancelAgreeFailed("remove comment failed failed: %s", err.Error())
	}
	return nil
}

func (r *CommentUseCase) RemoveSubComment(ctx context.Context, id int32, uuid string) error {
	err := r.repo.SendSubCommentToMq(ctx, &SubComment{
		CommentId: id,
		Uuid:      uuid,
	}, "remove_sub_comment_db_and_cache")
	if err != nil {
		return v1.ErrorCancelAgreeFailed("remove comment failed failed: %s", err.Error())
	}
	return nil
}

func (r *CommentUseCase) SetCommentAgree(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error {
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		err := r.repo.SetUserCommentAgreeToCache(ctx, id, userUuid)
		if err != nil {
			return v1.ErrorSetAgreeFailed("set comment agree to cache failed: %s", err.Error())
		}
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		err := r.repo.SendCommentAgreeToMq(ctx, id, creationId, creationType, uuid, userUuid, "set_comment_agree_db_and_cache")
		if err != nil {
			return v1.ErrorSetAgreeFailed("set comment agree to mq failed: %s", err.Error())
		}
		return nil
	}))
	return g.Wait()
}

func (r *CommentUseCase) SetSubCommentAgree(ctx context.Context, id int32, uuid, userUuid string) error {
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		err := r.repo.SetUserCommentAgreeToCache(ctx, id, userUuid)
		if err != nil {
			return v1.ErrorSetAgreeFailed("set comment agree to cache failed: %s", err.Error())
		}
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		err := r.repo.SendSubCommentAgreeToMq(ctx, id, uuid, userUuid, "set_sub_comment_agree_db_and_cache")
		if err != nil {
			return v1.ErrorSetAgreeFailed("set comment agree to mq failed: %s", err.Error())
		}
		return nil
	}))
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
		err = r.repo.SendCommentStatisticToMq(ctx, uuid, userUuid, "agree")
		if err != nil {
			return v1.ErrorSetAgreeFailed("set comment agree to mq failed: %s", err.Error())
		}
		err = r.repo.SendScoreToMq(ctx, 2, uuid, "add_score")
		if err != nil {
			return v1.ErrorCreateCommentFailed("send 2 score to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *CommentUseCase) SetSubCommentAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetSubCommentAgree(ctx, id, uuid)
		if err != nil {
			return v1.ErrorSetAgreeFailed("set comment agree failed: %s", err.Error())
		}
		err = r.repo.SetUserCommentAgree(ctx, id, userUuid)
		if err != nil {
			return v1.ErrorSetAgreeFailed("set comment agree failed: %s", err.Error())
		}
		err = r.repo.SetSubCommentAgreeToCache(ctx, id, uuid, userUuid)
		if err != nil {
			return v1.ErrorSetAgreeFailed("set sub comment agree to cache failed: %s", err.Error())
		}
		err = r.repo.SendCommentStatisticToMq(ctx, uuid, userUuid, "agree")
		if err != nil {
			return v1.ErrorSetAgreeFailed("set comment agree to mq failed: %s", err.Error())
		}
		err = r.repo.SendScoreToMq(ctx, 2, uuid, "add_score")
		if err != nil {
			return v1.ErrorCreateCommentFailed("send 2 score to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *CommentUseCase) CancelCommentAgree(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error {
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		err := r.repo.CancelUserCommentAgreeFromCache(ctx, id, userUuid)
		if err != nil {
			return v1.ErrorCancelAgreeFailed("cancel comment agree from cache failed: %s", err.Error())
		}
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		err := r.repo.SendCommentAgreeToMq(ctx, id, creationId, creationType, uuid, userUuid, "cancel_comment_agree_db_and_cache")
		if err != nil {
			return v1.ErrorCancelAgreeFailed("cancel comment agree from mq failed: %s", err.Error())
		}
		return nil
	}))
	return g.Wait()
}

func (r *CommentUseCase) CancelSubCommentAgree(ctx context.Context, id int32, uuid, userUuid string) error {
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		err := r.repo.CancelUserCommentAgreeFromCache(ctx, id, userUuid)
		if err != nil {
			return v1.ErrorCancelAgreeFailed("cancel sub comment agree from cache failed: %s", err.Error())
		}
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		err := r.repo.SendSubCommentAgreeToMq(ctx, id, uuid, userUuid, "cancel_sub_comment_agree_db_and_cache")
		if err != nil {
			return v1.ErrorCancelAgreeFailed("cancel sub comment agree from mq failed: %s", err.Error())
		}
		return nil
	}))
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
		err = r.repo.SendCommentStatisticToMq(ctx, uuid, userUuid, "agree_cancel")
		if err != nil {
			return v1.ErrorCancelAgreeFailed("cancel comment agree from mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *CommentUseCase) CancelSubCommentAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.CancelSubCommentAgree(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCancelAgreeFailed("cancel comment agree failed: %s", err.Error())
		}
		err = r.repo.CancelUserCommentAgree(ctx, id, userUuid)
		if err != nil {
			return v1.ErrorCancelAgreeFailed("cancel comment agree failed: %s", err.Error())
		}
		err = r.repo.CancelSubCommentAgreeFromCache(ctx, id, uuid, userUuid)
		if err != nil {
			return v1.ErrorCancelAgreeFailed("cancel comment agree from cache failed: %s", err.Error())
		}
		err = r.repo.SendCommentStatisticToMq(ctx, uuid, userUuid, "agree_cancel")
		if err != nil {
			return v1.ErrorCancelAgreeFailed("cancel comment agree from mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *CommentUseCase) AddCommentContentReviewDbAndCache(ctx context.Context, review *TextReview) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		review, err := r.repo.SetCommentContentIrregular(ctx, review)
		if err != nil {
			return v1.ErrorSetContentIrregularFailed("set comment content irregular failed: %s", err.Error())
		}

		err = r.repo.SetCommentContentIrregularToCache(ctx, review)
		if err != nil {
			return v1.ErrorSetContentIrregularFailed("set comment content irregular to cache failed: %s", err.Error())
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

		creationAuthor, err := r.getCreationAuthor(ctx, creationId, creationType)
		if err != nil {
			return v1.ErrorCreateCommentFailed("get creation author failed: %s", err.Error())
		}

		err = r.repo.CreateComment(ctx, id, creationId, creationType, creationAuthor, uuid)
		if err != nil {
			return v1.ErrorCreateCommentFailed("create comment failed: %s", err.Error())
		}

		err = r.addCommentUser(ctx, creationType, creationAuthor, uuid)
		if err != nil {
			return v1.ErrorCreateCommentFailed("create comment user failed: %s", err.Error())
		}

		err = r.repo.CreateCommentCache(ctx, id, creationId, creationType, creationAuthor, uuid)
		if err != nil {
			return v1.ErrorCreateCommentFailed("create comment cache failed: %s", err.Error())
		}

		err = r.repo.SendScoreToMq(ctx, 5, uuid, "add_score")
		if err != nil {
			return v1.ErrorCreateCommentFailed("send 5 score to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *CommentUseCase) getCreationAuthor(ctx context.Context, creationId, creationType int32) (string, error) {
	if creationType == 1 {
		return r.repo.GetArticleAuthor(ctx, creationId)
	} else if creationType == 3 {
		return r.repo.GetTalkAuthor(ctx, creationId)
	} else {
		return "", errors.Errorf(fmt.Sprintf("fail to get creation author: reason(%v), creationId(%v), creationType(%v)", "creationType not correct ", creationId, creationType))
	}
}

func (r *CommentUseCase) addCommentUser(ctx context.Context, creationType int32, creationAuthor, uuid string) error {
	switch creationType {
	case 1:
		return r.repo.AddArticleCommentUser(ctx, creationAuthor, uuid)
	case 3:
		return r.repo.AddTalkCommentUser(ctx, creationAuthor, uuid)
	}
	return nil
}

func (r *CommentUseCase) CreateSubCommentDbAndCache(ctx context.Context, id, rootId, parentId int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		var parentUserId string
		if parentId != 0 {
			parentUserId, err = r.repo.GetParentCommentUserId(ctx, parentId, rootId)
			if err != nil {
				return v1.ErrorCreateCommentFailed("get sub comment parent id failed: %s", err.Error())
			}
		}
		err = r.repo.DeleteCommentDraft(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCreateCommentFailed("delete sub comment draft failed: %s", err.Error())
		}

		rootComment, err := r.repo.GetRootComment(ctx, rootId)
		if err != nil {
			return v1.ErrorCreateCommentFailed("get sub comment parent id failed: %s", err.Error())
		}

		err = r.repo.CreateSubComment(ctx, id, rootId, parentId, rootComment.CreationId, rootComment.CreationType, rootComment.CreationAuthor, rootComment.Uuid, uuid, parentUserId)
		if err != nil {
			return v1.ErrorCreateCommentFailed("create sub comment failed: %s", err.Error())
		}

		err = r.repo.SetCommentComment(ctx, rootId)
		if err != nil {
			return v1.ErrorCreateCommentFailed("create sub comment failed: %s", err.Error())
		}

		err = r.addSubCommentUser(ctx, rootComment.CreationType, parentId, rootComment.Uuid, parentUserId, uuid)
		if err != nil {
			return v1.ErrorCreateCommentFailed("create comment user failed: %s", err.Error())
		}

		err = r.repo.CreateSubCommentCache(ctx, id, rootId, parentId, rootComment.CreationId, rootComment.CreationType, rootComment.CreationAuthor, rootComment.Uuid, uuid, parentUserId)
		if err != nil {
			return v1.ErrorCreateCommentFailed("create sub comment cache failed: %s", err.Error())
		}

		err = r.repo.SendScoreToMq(ctx, 5, uuid, "add_score")
		if err != nil {
			return v1.ErrorCreateCommentFailed("send 5 score to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *CommentUseCase) addSubCommentUser(ctx context.Context, creationType, parentId int32, rootUser, parentUser, uuid string) error {
	switch creationType {
	case 1:
		return r.repo.AddArticleSubCommentUser(ctx, parentId, rootUser, parentUser, uuid)
	case 3:
		return r.repo.AddTalkSubCommentUser(ctx, parentId, rootUser, parentUser, uuid)
	}
	return nil
}

func (r *CommentUseCase) RemoveCommentDbAndCache(ctx context.Context, id int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		rootComment, err := r.repo.GetRootComment(ctx, id)
		if err != nil {
			return v1.ErrorRemoveCommentFailed("get comment info failed: %s", err.Error())
		}
		err = r.repo.RemoveComment(ctx, id, uuid)
		if err != nil {
			return v1.ErrorRemoveCommentFailed("remove comment failed: %s", err.Error())
		}

		err = r.repo.RemoveCommentAgree(ctx, id, rootComment.Uuid)
		if err != nil {
			return v1.ErrorRemoveCommentFailed("remove comment agree record failed: %s", err.Error())
		}

		err = r.reduceCommentUser(ctx, rootComment.CreationType, rootComment.CreationAuthor, rootComment.Uuid)
		if err != nil {
			return v1.ErrorRemoveCommentFailed("reduce comment user failed: %s", err.Error())
		}

		err = r.repo.RemoveCommentCache(ctx, id, rootComment.CreationId, rootComment.CreationType, rootComment.CreationAuthor, rootComment.Uuid)
		if err != nil {
			return v1.ErrorRemoveCommentFailed("remove comment cache failed: %s", err.Error())
		}

		err = r.repo.ReduceCreationComment(ctx, rootComment.CreationId, rootComment.CreationType, rootComment.CreationAuthor)
		if err != nil {
			return v1.ErrorRemoveCommentFailed("reduce creation comment failed: %s", err.Error())
		}

		return nil
	})
}

func (r *CommentUseCase) reduceCommentUser(ctx context.Context, creationType int32, creationAuthor, uuid string) error {
	switch creationType {
	case 1:
		return r.repo.ReduceArticleCommentUser(ctx, creationAuthor, uuid)
	case 3:
		return r.repo.ReduceTalkCommentUser(ctx, creationAuthor, uuid)
	}
	return nil
}

func (r *CommentUseCase) RemoveSubCommentDbAndCache(ctx context.Context, id int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		subComment, err := r.repo.GetSubComment(ctx, id)
		if err != nil {
			return v1.ErrorRemoveCommentFailed("get sub comment failed: %s", err.Error())
		}

		err = r.repo.RemoveSubComment(ctx, id, subComment.RootId, uuid)
		if err != nil {
			return v1.ErrorRemoveCommentFailed("remove sub comment failed: %s", err.Error())
		}

		err = r.repo.RemoveCommentAgree(ctx, id, uuid)
		if err != nil {
			return v1.ErrorRemoveCommentFailed("remove sub comment agree record failed: %s", err.Error())
		}

		err = r.reduceSubCommentUser(ctx, subComment.CreationType, subComment.ParentId, subComment.Uuid, subComment.Reply, uuid)
		if err != nil {
			return v1.ErrorCreateCommentFailed("create comment user failed: %s", err.Error())
		}

		err = r.repo.RemoveSubCommentCache(ctx, id, subComment.RootId, subComment.ParentId, subComment.CreationId, subComment.CreationType, subComment.CreationAuthor, subComment.RootUser, subComment.Uuid, subComment.Reply)
		if err != nil {
			return v1.ErrorRemoveCommentFailed("remove sub comment cache failed: %s", err.Error())
		}
		return nil
	})
}

func (r *CommentUseCase) reduceSubCommentUser(ctx context.Context, creationType, parentId int32, rootUser, parentUser, uuid string) error {
	switch creationType {
	case 1:
		return r.repo.ReduceArticleSubCommentUser(ctx, parentId, rootUser, parentUser, uuid)
	case 3:
		return r.repo.ReduceTalkSubCommentUser(ctx, parentId, rootUser, parentUser, uuid)
	}
	return nil
}
