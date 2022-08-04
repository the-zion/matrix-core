package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/comment/service/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *CommentService) GetHealth(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *CommentService) GetLastCommentDraft(ctx context.Context, req *v1.GetLastCommentDraftReq) (*v1.GetLastCommentDraftReply, error) {
	draft, err := s.cc.GetLastCommentDraft(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetLastCommentDraftReply{
		Id: draft.Id,
	}, nil
}

func (s *CommentService) CreateCommentDraft(ctx context.Context, req *v1.CreateCommentDraftReq) (*v1.CreateCommentDraftReply, error) {
	id, err := s.cc.CreateCommentDraft(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.CreateCommentDraftReply{
		Id: id,
	}, nil
}

func (s *CommentService) SendComment(ctx context.Context, req *v1.SendCommentReq) (*emptypb.Empty, error) {
	err := s.cc.SendComment(ctx, req.Id, req.Uuid, req.Ip)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
