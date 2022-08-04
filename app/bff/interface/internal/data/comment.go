package data

import (
	"context"
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
