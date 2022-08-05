package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	commentV1 "github.com/the-zion/matrix-core/api/comment/service/v1"
	"github.com/the-zion/matrix-core/app/bff/interface/internal/biz"
	"golang.org/x/sync/singleflight"
)

var _ biz.CommentRepo = (*commentRepo)(nil)

type commentRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

func NewCommentRepo(data *Data, logger log.Logger) biz.CommentRepo {
	return &commentRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "bff/data/comment")),
		sg:   &singleflight.Group{},
	}
}

func (r *commentRepo) GetLastCommentDraft(ctx context.Context, uuid string) (*biz.CommentDraft, error) {
	reply, err := r.data.commc.GetLastCommentDraft(ctx, &commentV1.GetLastCommentDraftReq{
		Uuid: uuid,
	})
	if err != nil {
		return nil, err
	}
	return &biz.CommentDraft{
		Id: reply.Id,
	}, nil
}

func (r *commentRepo) GetCommentList(ctx context.Context, page, creationId, creationType int32) ([]*biz.Comment, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("comment_page_%v_%v_%v", creationId, creationType, page), func() (interface{}, error) {
		reply := make([]*biz.Comment, 0)
		commentList, err := r.data.commc.GetCommentList(ctx, &commentV1.GetCommentListReq{
			Page:         page,
			CreationId:   creationId,
			CreationType: creationType,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range commentList.Comment {
			reply = append(reply, &biz.Comment{
				Id:   item.Id,
				Uuid: item.Uuid,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Comment), nil
}

func (r *commentRepo) CreateCommentDraft(ctx context.Context, uuid string) (int32, error) {
	reply, err := r.data.commc.CreateCommentDraft(ctx, &commentV1.CreateCommentDraftReq{
		Uuid: uuid,
	})
	if err != nil {
		return 0, err
	}
	return reply.Id, nil
}

func (r *commentRepo) SendComment(ctx context.Context, id int32, uuid, ip string) error {
	_, err := r.data.commc.SendComment(ctx, &commentV1.SendCommentReq{
		Id:   id,
		Uuid: uuid,
		Ip:   ip,
	})
	if err != nil {
		return err
	}
	return nil
}
