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
		Id:     draft.Id,
		Status: draft.Status,
	}, nil
}

func (s *CommentService) GetCommentList(ctx context.Context, req *v1.GetCommentListReq) (*v1.GetCommentListReply, error) {
	reply := &v1.GetCommentListReply{Comment: make([]*v1.GetCommentListReply_Comment, 0)}
	commentList, err := s.cc.GetCommentList(ctx, req.Page, req.CreationId, req.CreationType)
	if err != nil {
		return nil, err
	}
	for _, item := range commentList {
		reply.Comment = append(reply.Comment, &v1.GetCommentListReply_Comment{
			Id:   item.CommentId,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *CommentService) GetCommentListHot(ctx context.Context, req *v1.GetCommentListReq) (*v1.GetCommentListReply, error) {
	reply := &v1.GetCommentListReply{Comment: make([]*v1.GetCommentListReply_Comment, 0)}
	commentList, err := s.cc.GetCommentListHot(ctx, req.Page, req.CreationId, req.CreationType)
	if err != nil {
		return nil, err
	}
	for _, item := range commentList {
		reply.Comment = append(reply.Comment, &v1.GetCommentListReply_Comment{
			Id:   item.CommentId,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *CommentService) GetCommentListStatistic(ctx context.Context, req *v1.GetCommentListStatisticReq) (*v1.GetCommentListStatisticReply, error) {
	reply := &v1.GetCommentListStatisticReply{Count: make([]*v1.GetCommentListStatisticReply_Count, 0)}
	commentListStatistic, err := s.cc.GetCommentListStatistic(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	for _, item := range commentListStatistic {
		reply.Count = append(reply.Count, &v1.GetCommentListStatisticReply_Count{
			Id:      item.CommentId,
			Agree:   item.Agree,
			Comment: item.Comment,
		})
	}
	return reply, nil
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

func (s *CommentService) CreateComment(ctx context.Context, req *v1.CreateCommentReq) (*emptypb.Empty, error) {
	err := s.cc.CreateComment(ctx, req.Id, req.CreationId, req.CreationType, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommentService) CreateCommentDbCacheAndSearch(ctx context.Context, req *v1.CreateCommentDbCacheAndSearchReq) (*emptypb.Empty, error) {
	err := s.cc.CreateCommentDbCacheAndSearch(ctx, req.Id, req.CreationId, req.CreationType, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommentService) SendComment(ctx context.Context, req *v1.SendCommentReq) (*emptypb.Empty, error) {
	err := s.cc.SendComment(ctx, req.Id, req.Uuid, req.Ip)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
