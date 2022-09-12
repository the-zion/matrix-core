package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/tencentyun/cos-go-sdk-v5"
	commentV1 "github.com/the-zion/matrix-core/api/comment/service/v1"
	"github.com/the-zion/matrix-core/app/message/service/internal/biz"
	"strconv"
)

type commentRepo struct {
	data *Data
	log  *log.Helper
}

func NewCommentRepo(data *Data, logger log.Logger) biz.CommentRepo {
	return &commentRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "message/data/comment")),
	}
}

func (r *commentRepo) ToReviewCreateComment(id int32, uuid string) error {
	opt := &cos.PutTextAuditingJobOptions{
		InputObject: "comment/" + uuid + "/" + strconv.Itoa(int(id)) + "/content",
		Conf: &cos.TextAuditingJobConf{
			CallbackVersion: "Detail",
			Callback:        r.data.cosCommentCli.callback["comment_create"],
		},
	}

	_, _, err := r.data.cosCommentCli.cos.CI.PutTextAuditingJob(context.Background(), opt)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send comment create review request to cos: id(%v) uuid(%s)", id, uuid))
	}
	return nil
}

func (r *commentRepo) ToReviewCreateSubComment(id int32, uuid string) error {
	opt := &cos.PutTextAuditingJobOptions{
		InputObject: "comment/" + uuid + "/" + strconv.Itoa(int(id)) + "/content",
		Conf: &cos.TextAuditingJobConf{
			CallbackVersion: "Detail",
			Callback:        r.data.cosCommentCli.callback["sub_comment_create"],
		},
	}

	_, _, err := r.data.cosCommentCli.cos.CI.PutTextAuditingJob(context.Background(), opt)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send sub comment create review request to cos: id(%v) uuid(%s)", id, uuid))
	}
	return nil
}

func (r *commentRepo) CommentCreateReviewPass(ctx context.Context, id, creationId, creationType int32, uuid string) error {
	_, err := r.data.commc.CreateComment(ctx, &commentV1.CreateCommentReq{
		Id:           id,
		CreationId:   creationId,
		CreationType: creationType,
		Uuid:         uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepo) CommentContentIrregular(ctx context.Context, review *biz.TextReview, id int32, comment, kind, uuid string) error {
	_, err := r.data.commc.CommentContentIrregular(ctx, &commentV1.CommentContentIrregularReq{
		Id:      id,
		Uuid:    uuid,
		Comment: comment,
		Kind:    kind,
		JobId:   review.JobId,
		Label:   review.Label,
		Result:  review.Result,
		Section: review.Section,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepo) SubCommentCreateReviewPass(ctx context.Context, id, rootId, parentId int32, uuid string) error {
	_, err := r.data.commc.CreateSubComment(ctx, &commentV1.CreateSubCommentReq{
		Id:       id,
		RootId:   rootId,
		ParentId: parentId,
		Uuid:     uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepo) CreateCommentDbAndCache(ctx context.Context, id, createId, createType int32, uuid string) error {
	_, err := r.data.commc.CreateCommentDbAndCache(ctx, &commentV1.CreateCommentDbAndCacheReq{
		Id:           id,
		CreationId:   createId,
		CreationType: createType,
		Uuid:         uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepo) CreateSubCommentDbAndCache(ctx context.Context, id, rootId, parentId int32, uuid string) error {
	_, err := r.data.commc.CreateSubCommentDbAndCache(ctx, &commentV1.CreateSubCommentDbAndCacheReq{
		Id:       id,
		RootId:   rootId,
		ParentId: parentId,
		Uuid:     uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepo) RemoveCommentDbAndCache(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.commc.RemoveCommentDbAndCache(ctx, &commentV1.RemoveCommentDbAndCacheReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepo) RemoveSubCommentDbAndCache(ctx context.Context, id int32, uuid string) error {
	_, err := r.data.commc.RemoveSubCommentDbAndCache(ctx, &commentV1.RemoveSubCommentDbAndCacheReq{
		Id:   id,
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepo) SetCommentAgreeDbAndCache(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error {
	_, err := r.data.commc.SetCommentAgreeDbAndCache(ctx, &commentV1.SetCommentAgreeReq{
		Id:           id,
		Uuid:         uuid,
		UserUuid:     userUuid,
		CreationId:   creationId,
		CreationType: creationType,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepo) SetSubCommentAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.commc.SetSubCommentAgreeDbAndCache(ctx, &commentV1.SetSubCommentAgreeReq{
		Id:       id,
		Uuid:     uuid,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepo) CancelCommentAgreeDbAndCache(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error {
	_, err := r.data.commc.CancelCommentAgreeDbAndCache(ctx, &commentV1.CancelCommentAgreeReq{
		Id:           id,
		Uuid:         uuid,
		UserUuid:     userUuid,
		CreationId:   creationId,
		CreationType: creationType,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepo) CancelSubCommentAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	_, err := r.data.commc.CancelSubCommentAgreeDbAndCache(ctx, &commentV1.CancelSubCommentAgreeReq{
		Id:       id,
		Uuid:     uuid,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepo) AddCommentContentReviewDbAndCache(ctx context.Context, commentId, result int32, uuid, jobId, label, comment, kind string, section string) error {
	_, err := r.data.commc.AddCommentContentReviewDbAndCache(ctx, &commentV1.AddCommentContentReviewDbAndCacheReq{
		CommentId: commentId,
		Kind:      kind,
		Uuid:      uuid,
		JobId:     jobId,
		Label:     label,
		Result:    result,
		Comment:   comment,
		Section:   section,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepo) GetCommentUser(ctx context.Context, uuid string) (int32, error) {
	reply, err := r.data.commc.GetCommentUser(ctx, &commentV1.GetCommentUserReq{
		Uuid: uuid,
	})
	if err != nil {
		return 0, err
	}
	return reply.Comment, nil
}
