package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"net/url"
	"strconv"
)

type CommentRepo interface {
	ToReviewCreateComment(id int32, uuid string) error
	ToReviewCreateSubComment(id int32, uuid string) error
	CommentCreateReviewPass(ctx context.Context, id, creationId, creationType int32, uuid string) error
	CommentContentIrregular(ctx context.Context, review *TextReview, id int32, comment, kind, uuid string) error
	SubCommentCreateReviewPass(ctx context.Context, id, rootId, parentId int32, uuid string) error
	CreateCommentDbAndCache(ctx context.Context, id, createId, createType int32, uuid string) (string, error)
	CreateSubCommentDbAndCache(ctx context.Context, id, rootId, parentId int32, uuid string) (string, string, error)
	RemoveCommentDbAndCache(ctx context.Context, id int32, uuid string) error
	RemoveSubCommentDbAndCache(ctx context.Context, id int32, uuid string) error
	SetCommentAgreeDbAndCache(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error
	SetSubCommentAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error
	SetCommentCount(ctx context.Context, uuid string)
	SetSubCommentCount(ctx context.Context, uuid string)
	CancelCommentAgreeDbAndCache(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error
	CancelSubCommentAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error
	AddCommentContentReviewDbAndCache(ctx context.Context, commentId, result int32, uuid, jobId, label, comment, kind, section string) error
	GetCommentUser(ctx context.Context, uuid string) (int32, error)
}

type CommentUseCase struct {
	repo        CommentRepo
	messageRepo MessageRepo
	tm          Transaction
	jwt         Jwt
	log         *log.Helper
}

func NewCommentUseCase(repo CommentRepo, messageRepo MessageRepo, tm Transaction, jwt Jwt, logger log.Logger) *CommentUseCase {
	return &CommentUseCase{
		repo:        repo,
		messageRepo: messageRepo,
		tm:          tm,
		jwt:         jwt,
		log:         log.NewHelper(log.With(logger, "module", "message/biz/commentUseCase")),
	}
}

func (r *CommentUseCase) ToReviewCreateComment(id int32, uuid string) error {
	return r.repo.ToReviewCreateComment(id, uuid)
}

func (r *CommentUseCase) ToReviewCreateSubComment(id int32, uuid string) error {
	return r.repo.ToReviewCreateSubComment(id, uuid)
}

func (r *CommentUseCase) CommentCreateReview(ctx context.Context, tr *TextReview) error {
	var err error
	var token, id, creationId, creationType, comment, kind string
	var ok bool

	if token, ok = tr.CosHeaders["X-Cos-Meta-Token"]; !ok || token == "" {
		r.log.Info("token not exist，%v", tr)
		return nil
	}

	if id, ok = tr.CosHeaders["X-Cos-Meta-Id"]; !ok || id == "" {
		r.log.Info("id not exist，%v", tr)
		return nil
	}

	if creationId, ok = tr.CosHeaders["X-Cos-Meta-Creationid"]; !ok || creationId == "" {
		r.log.Info("creationId not exist，%v", tr)
		return nil
	}

	if creationType, ok = tr.CosHeaders["X-Cos-Meta-Creationtype"]; !ok || creationType == "" {
		r.log.Info("creationType not exist，%v", tr)
		return nil
	}

	if comment, ok = tr.CosHeaders["X-Cos-Meta-Comment"]; !ok || comment == "" {
		r.log.Info("title not exist，%v", tr)
		return nil
	}

	if kind, ok = tr.CosHeaders["X-Cos-Meta-Kind"]; !ok || kind == "" {
		r.log.Info("kind not exist，%v", tr)
		return nil
	}

	uuid, err := r.jwt.JwtCheck(token)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get uuid from token: %s", token))
	}

	aid, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: %v", tr))
	}

	cid, err := strconv.ParseInt(creationId, 10, 32)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: %v", tr))
	}

	cType, err := strconv.ParseInt(creationType, 10, 32)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: %v", tr))
	}

	comment, err = url.QueryUnescape(comment)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to decode string to UTF-8: %v", tr))
	}

	if tr.State != "Success" {
		r.log.Info("comment create review failed，%v", tr)
		return nil
	}

	if tr.Result == 0 {
		err = r.repo.CommentCreateReviewPass(ctx, int32(aid), int32(cid), int32(cType), uuid)
		if err != nil {
			return err
		}
	} else {
		err = r.repo.CommentContentIrregular(ctx, tr, int32(aid), comment, kind, uuid)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *CommentUseCase) SubCommentCreateReview(ctx context.Context, tr *TextReview) error {
	var err error
	var token, id, rootId, parentId, comment, kind string
	var ok bool

	if token, ok = tr.CosHeaders["X-Cos-Meta-Token"]; !ok || token == "" {
		r.log.Info("token not exist，%v", tr)
		return nil
	}

	if id, ok = tr.CosHeaders["X-Cos-Meta-Id"]; !ok || id == "" {
		r.log.Info("id not exist，%v", tr)
		return nil
	}

	if rootId, ok = tr.CosHeaders["X-Cos-Meta-Rootid"]; !ok || rootId == "" {
		r.log.Info("rootId not exist，%v", tr)
		return nil
	}

	if parentId, ok = tr.CosHeaders["X-Cos-Meta-Parentid"]; !ok {
		r.log.Info("parentId not exist，%v", tr)
		return nil
	}

	if comment, ok = tr.CosHeaders["X-Cos-Meta-Comment"]; !ok || comment == "" {
		r.log.Info("title not exist，%v", tr)
		return nil
	}

	if kind, ok = tr.CosHeaders["X-Cos-Meta-Kind"]; !ok || kind == "" {
		r.log.Info("kind not exist，%v", tr)
		return nil
	}

	if parentId == "" {
		parentId = "0"
	}

	uuid, err := r.jwt.JwtCheck(token)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get uuid from token: %s", token))
	}

	aid, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: %v", tr))
	}

	rid, err := strconv.ParseInt(rootId, 10, 32)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: %v", tr))
	}

	pid, err := strconv.ParseInt(parentId, 10, 32)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: %v", tr))
	}

	comment, err = url.QueryUnescape(comment)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to decode string to UTF-8: %v", tr))
	}

	if tr.State != "Success" {
		r.log.Info("comment create review failed，%v", tr)
		return nil
	}

	if tr.Result == 0 {
		err = r.repo.SubCommentCreateReviewPass(ctx, int32(aid), int32(rid), int32(pid), uuid)
		if err != nil {
			return err
		}
	} else {
		err = r.repo.CommentContentIrregular(ctx, tr, int32(aid), comment, kind, uuid)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *CommentUseCase) CreateCommentDbAndCache(ctx context.Context, id, createId, createType int32, uuid string) error {
	author, err := r.repo.CreateCommentDbAndCache(ctx, id, createId, createType, uuid)
	if err != nil {
		return err
	}
	if author != "" && author != uuid {
		r.repo.SetCommentCount(ctx, author)
	}
	return nil
}

func (r *CommentUseCase) CreateSubCommentDbAndCache(ctx context.Context, id, rootId, parentId int32, uuid string) error {
	root, parent, err := r.repo.CreateSubCommentDbAndCache(ctx, id, rootId, parentId, uuid)
	if err != nil {
		return err
	}
	if root != "" && root != uuid {
		r.repo.SetSubCommentCount(ctx, root)
	}

	if parent != "" && root != parent && parent != uuid {
		r.repo.SetSubCommentCount(ctx, parent)
	}

	return nil
}

func (r *CommentUseCase) RemoveCommentDbAndCache(ctx context.Context, id int32, uuid string) error {
	return r.repo.RemoveCommentDbAndCache(ctx, id, uuid)
}

func (r *CommentUseCase) RemoveSubCommentDbAndCache(ctx context.Context, id int32, uuid string) error {
	return r.repo.RemoveSubCommentDbAndCache(ctx, id, uuid)
}

func (r *CommentUseCase) SetCommentAgreeDbAndCache(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error {
	return r.repo.SetCommentAgreeDbAndCache(ctx, id, creationId, creationType, uuid, userUuid)
}

func (r *CommentUseCase) SetSubCommentAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.repo.SetSubCommentAgreeDbAndCache(ctx, id, uuid, userUuid)
}

func (r *CommentUseCase) CancelCommentAgreeDbAndCache(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error {
	return r.repo.CancelCommentAgreeDbAndCache(ctx, id, creationId, creationType, uuid, userUuid)
}

func (r *CommentUseCase) CancelSubCommentAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.repo.CancelSubCommentAgreeDbAndCache(ctx, id, uuid, userUuid)
}

func (r *CommentUseCase) AddCommentContentReviewDbAndCache(ctx context.Context, commentId, result int32, uuid, jobId, label, comment, kind, section string) error {
	err := r.repo.AddCommentContentReviewDbAndCache(ctx, commentId, result, uuid, jobId, label, comment, kind, section)
	if err != nil {
		return err
	}
	err = r.tm.ExecTx(ctx, func(ctx context.Context) error {
		notification, err := r.messageRepo.AddMailBoxSystemNotification(ctx, commentId, "comment-create", "", uuid, label, result, section, "", "", comment)
		if err != nil {
			return err
		}

		err = r.messageRepo.AddMailBoxSystemNotificationToCache(ctx, notification)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to add mail box system notification for comment: error(%v), uuid(%s)", err, uuid)
	}
	return nil
}
