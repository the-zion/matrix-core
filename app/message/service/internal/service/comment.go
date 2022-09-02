package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/message/service/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *MessageService) ToReviewCreateComment(id int32, uuid string) error {
	return s.commc.ToReviewCreateComment(id, uuid)
}

func (s *MessageService) ToReviewCreateSubComment(id int32, uuid string) error {
	return s.commc.ToReviewCreateSubComment(id, uuid)
}

func (s *MessageService) CommentCreateReview(ctx context.Context, req *v1.TextReviewReq) (*emptypb.Empty, error) {
	tr := s.TextReview(req)
	err := s.commc.CommentCreateReview(ctx, tr)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *MessageService) SubCommentCreateReview(ctx context.Context, req *v1.TextReviewReq) (*emptypb.Empty, error) {
	tr := s.TextReview(req)
	err := s.commc.SubCommentCreateReview(ctx, tr)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *MessageService) CreateCommentDbAndCache(ctx context.Context, id, createId, createType int32, uuid string) error {
	return s.commc.CreateCommentDbAndCache(ctx, id, createId, createType, uuid)
}

func (s *MessageService) CreateSubCommentDbAndCache(ctx context.Context, id, rootId, parentId int32, uuid string) error {
	return s.commc.CreateSubCommentDbAndCache(ctx, id, rootId, parentId, uuid)
}

func (s *MessageService) RemoveCommentDbAndCache(ctx context.Context, id int32, uuid string) error {
	return s.commc.RemoveCommentDbAndCache(ctx, id, uuid)
}

func (s *MessageService) RemoveSubCommentDbAndCache(ctx context.Context, id int32, uuid string) error {
	return s.commc.RemoveSubCommentDbAndCache(ctx, id, uuid)
}

func (s *MessageService) SetCommentAgreeDbAndCache(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error {
	return s.commc.SetCommentAgreeDbAndCache(ctx, id, creationId, creationType, uuid, userUuid)
}

func (s *MessageService) SetSubCommentAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return s.commc.SetSubCommentAgreeDbAndCache(ctx, id, uuid, userUuid)
}

func (s *MessageService) CancelCommentAgreeDbAndCache(ctx context.Context, id, creationId, creationType int32, uuid, userUuid string) error {
	return s.commc.CancelCommentAgreeDbAndCache(ctx, id, creationId, creationType, uuid, userUuid)
}

func (s *MessageService) CancelSubCommentAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return s.commc.CancelSubCommentAgreeDbAndCache(ctx, id, uuid, userUuid)
}
