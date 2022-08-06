package service

import (
	"context"
	"github.com/the-zion/matrix-core/api/bff/interface/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *BffService) CreateCommentDraft(ctx context.Context, _ *emptypb.Empty) (*v1.CreateCommentDraftReply, error) {
	id, err := s.commc.CreateCommentDraft(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.CreateCommentDraftReply{
		Id: id,
	}, nil
}

func (s *BffService) GetLastCommentDraft(ctx context.Context, _ *emptypb.Empty) (*v1.GetLastCommentDraftReply, error) {
	draft, err := s.commc.GetLastCommentDraft(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetLastCommentDraftReply{
		Id:     draft.Id,
		Status: draft.Status,
	}, nil
}

func (s *BffService) GetCommentList(ctx context.Context, req *v1.GetCommentListReq) (*v1.GetCommentListReply, error) {
	reply := &v1.GetCommentListReply{Comment: make([]*v1.GetCommentListReply_Comment, 0)}
	commentList, err := s.commc.GetCommentList(ctx, req.Page, req.CreationId, req.CreationType)
	if err != nil {
		return nil, err
	}
	for _, item := range commentList {
		reply.Comment = append(reply.Comment, &v1.GetCommentListReply_Comment{
			Id:       item.Id,
			Uuid:     item.Uuid,
			Agree:    item.Agree,
			Comment:  item.Comment,
			Username: item.UserName,
		})
	}
	return reply, nil
}

func (s *BffService) GetCommentListHot(ctx context.Context, req *v1.GetCommentListReq) (*v1.GetCommentListReply, error) {
	reply := &v1.GetCommentListReply{Comment: make([]*v1.GetCommentListReply_Comment, 0)}
	commentList, err := s.commc.GetCommentListHot(ctx, req.Page, req.CreationId, req.CreationType)
	if err != nil {
		return nil, err
	}
	for _, item := range commentList {
		reply.Comment = append(reply.Comment, &v1.GetCommentListReply_Comment{
			Id:       item.Id,
			Uuid:     item.Uuid,
			Agree:    item.Agree,
			Comment:  item.Comment,
			Username: item.UserName,
		})
	}
	return reply, nil
}

func (s *BffService) SendComment(ctx context.Context, req *v1.SendCommentReq) (*emptypb.Empty, error) {
	err := s.commc.SendComment(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
