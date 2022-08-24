package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strconv"
)

type CommentRepo interface {
	ToReviewCreateComment(id int32, uuid string) error
	ToReviewCreateSubComment(id int32, uuid string) error
	CommentCreateReviewPass(ctx context.Context, id, creationId, creationType int32, uuid string) error
	SubCommentCreateReviewPass(ctx context.Context, id, rootId, parentId int32, uuid string) error
	CreateCommentDbAndCache(ctx context.Context, id, createId, createType int32, uuid string) error
	CreateSubCommentDbAndCache(ctx context.Context, id, rootId, parentId int32, uuid string) error
	RemoveCommentDbAndCache(ctx context.Context, id, createId, createType int32, uuid string) error
	RemoveSubCommentDbAndCache(ctx context.Context, id, rootId int32, uuid string) error
	SetCommentAgreeDbAndCache(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error
	SetSubCommentAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error
	CancelCommentAgreeDbAndCache(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error
	CancelSubCommentAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error
	GetCommentUser(ctx context.Context, uuid string) (int32, error)
}

type CommentUseCase struct {
	repo CommentRepo
	jwt  Jwt
	log  *log.Helper
}

func NewCommentUseCase(repo CommentRepo, jwt Jwt, logger log.Logger) *CommentUseCase {
	return &CommentUseCase{
		repo: repo,
		jwt:  jwt,
		log:  log.NewHelper(log.With(logger, "module", "message/biz/commentUseCase")),
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
	var token, id, creationId, creationType string
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

	if tr.State != "Success" {
		r.log.Info("comment create review failed，%v", tr)
		return nil
	}

	if tr.Result == 0 {
		err = r.repo.CommentCreateReviewPass(ctx, int32(aid), int32(cid), int32(cType), uuid)
	} else {
		r.log.Info("comment create review not pass，%v", tr)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *CommentUseCase) SubCommentCreateReview(ctx context.Context, tr *TextReview) error {
	var err error
	var token, id, rootId, parentId string
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

	if tr.State != "Success" {
		r.log.Info("comment create review failed，%v", tr)
		return nil
	}

	if tr.Result == 0 {
		err = r.repo.SubCommentCreateReviewPass(ctx, int32(aid), int32(rid), int32(pid), uuid)
	} else {
		r.log.Info("sub comment create review not pass，%v", tr)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *CommentUseCase) CreateCommentDbAndCache(ctx context.Context, id, createId, createType int32, uuid string) error {
	return r.repo.CreateCommentDbAndCache(ctx, id, createId, createType, uuid)
}

func (r *CommentUseCase) CreateSubCommentDbAndCache(ctx context.Context, id, rootId, parentId int32, uuid string) error {
	return r.repo.CreateSubCommentDbAndCache(ctx, id, rootId, parentId, uuid)
}

func (r *CommentUseCase) RemoveCommentDbAndCache(ctx context.Context, id, createId, createType int32, uuid string) error {
	return r.repo.RemoveCommentDbAndCache(ctx, id, createId, createType, uuid)
}

func (r *CommentUseCase) RemoveSubCommentDbAndCache(ctx context.Context, id, rootId int32, uuid string) error {
	return r.repo.RemoveSubCommentDbAndCache(ctx, id, rootId, uuid)
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
