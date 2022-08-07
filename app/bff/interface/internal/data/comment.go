package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	commentV1 "github.com/the-zion/matrix-core/api/comment/service/v1"
	userV1 "github.com/the-zion/matrix-core/api/user/service/v1"
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
		Id:     reply.Id,
		Status: reply.Status,
	}, nil
}

func (r *commentRepo) GetUserCommentAgree(ctx context.Context, uuid string) (map[int32]bool, error) {
	reply, err := r.data.commc.GetUserCommentAgree(ctx, &commentV1.GetUserCommentAgreeReq{
		Uuid: uuid,
	})
	if err != nil {
		return nil, err
	}
	return reply.Agree, nil
}

func (r *commentRepo) GetCommentList(ctx context.Context, page, creationId, creationType int32) ([]*biz.Comment, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("comment_%v_%v_%v", creationId, creationType, page), func() (interface{}, error) {
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

func (r *commentRepo) GetCommentListHot(ctx context.Context, page, creationId, creationType int32) ([]*biz.Comment, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("comment_hot_%v_%v_%v", creationId, creationType, page), func() (interface{}, error) {
		reply := make([]*biz.Comment, 0)
		commentList, err := r.data.commc.GetCommentListHot(ctx, &commentV1.GetCommentListReq{
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

func (r *commentRepo) GetCommentListStatistic(ctx context.Context, page, creationId, creationType int32, key string, commentList []*biz.Comment) ([]*biz.CommentStatistic, error) {
	ids := make([]int32, 0)
	for _, item := range commentList {
		ids = append(ids, item.Id)
	}
	result, err, _ := r.sg.Do(fmt.Sprintf("%s_%v_%v_%v", key, creationId, creationType, page), func() (interface{}, error) {
		reply := make([]*biz.CommentStatistic, 0)
		commentListStatistic, err := r.data.commc.GetCommentListStatistic(ctx, &commentV1.GetCommentListStatisticReq{
			Ids: ids,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range commentListStatistic.Count {
			reply = append(reply, &biz.CommentStatistic{
				Id:      item.Id,
				Agree:   item.Agree,
				Comment: item.Comment,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.CommentStatistic), nil
}

func (r *commentRepo) GetUserProfileList(ctx context.Context, page, creationId, creationType int32, key string, commentList []*biz.Comment) ([]*biz.UserProfile, error) {
	uuids := make([]string, 0)
	for _, item := range commentList {
		uuids = append(uuids, item.Uuid)
	}
	result, err, _ := r.sg.Do(fmt.Sprintf("%s_%v_%v_%v", key, creationId, creationType, page), func() (interface{}, error) {
		reply := make([]*biz.UserProfile, 0)
		userProfileList, err := r.data.uc.GetProfileList(ctx, &userV1.GetProfileListReq{
			Uuids: uuids,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range userProfileList.Profile {
			reply = append(reply, &biz.UserProfile{
				Uuid:     item.Uuid,
				Username: item.Username,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.UserProfile), nil
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

func (r *commentRepo) RemoveComment(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error {
	_, err := r.data.commc.RemoveComment(ctx, &commentV1.RemoveCommentReq{
		Id:           id,
		CreationId:   creationId,
		CreationType: creationType,
		Uuid:         uuid,
		UserUuid:     userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepo) SetCommentAgree(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error {
	_, err := r.data.commc.SetCommentAgree(ctx, &commentV1.SetCommentAgreeReq{
		Uuid:         uuid,
		Id:           id,
		UserUuid:     userUuid,
		CreationId:   creationId,
		CreationType: creationType,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepo) CancelCommentAgree(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error {
	_, err := r.data.commc.CancelCommentAgree(ctx, &commentV1.CancelCommentAgreeReq{
		Uuid:         uuid,
		Id:           id,
		UserUuid:     userUuid,
		CreationId:   creationId,
		CreationType: creationType,
	})
	if err != nil {
		return err
	}
	return nil
}
