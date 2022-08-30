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

func (s *BffService) GetUserCommentAgree(ctx context.Context, _ *emptypb.Empty) (*v1.GetUserCommentAgreeReply, error) {
	agreeMap, err := s.commc.GetUserCommentAgree(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserCommentAgreeReply{
		Agree: agreeMap,
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

func (s *BffService) GetSubCommentList(ctx context.Context, req *v1.GetSubCommentListReq) (*v1.GetSubCommentListReply, error) {
	reply := &v1.GetSubCommentListReply{Comment: make([]*v1.GetSubCommentListReply_Comment, 0)}
	commentList, err := s.commc.GetSubCommentList(ctx, req.Page, req.Id)
	if err != nil {
		return nil, err
	}
	for _, item := range commentList {
		reply.Comment = append(reply.Comment, &v1.GetSubCommentListReply_Comment{
			Id:        item.Id,
			Uuid:      item.Uuid,
			Reply:     item.Reply,
			Agree:     item.Agree,
			Username:  item.UserName,
			ReplyName: item.ReplyName,
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

func (s *BffService) GetUserCommentArticleReplyList(ctx context.Context, req *v1.GetUserCommentArticleReplyListReq) (*v1.GetUserCommentArticleReplyListReply, error) {
	reply := &v1.GetUserCommentArticleReplyListReply{List: make([]*v1.GetUserCommentArticleReplyListReply_List, 0)}
	commentList, err := s.commc.GetUserCommentArticleReplyList(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range commentList {
		reply.List = append(reply.List, &v1.GetUserCommentArticleReplyListReply_List{
			Id:             item.Id,
			CreationId:     item.CreationId,
			CreationAuthor: item.CreationAuthor,
		})
	}
	return reply, nil
}

func (s *BffService) GetUserSubCommentArticleReplyList(ctx context.Context, req *v1.GetUserSubCommentArticleReplyListReq) (*v1.GetUserSubCommentArticleReplyListReply, error) {
	reply := &v1.GetUserSubCommentArticleReplyListReply{List: make([]*v1.GetUserSubCommentArticleReplyListReply_List, 0)}
	commentList, err := s.commc.GetUserSubCommentArticleReplyList(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range commentList {
		reply.List = append(reply.List, &v1.GetUserSubCommentArticleReplyListReply_List{
			Id:             item.Id,
			CreationId:     item.CreationId,
			CreationAuthor: item.CreationAuthor,
		})
	}
	return reply, nil
}

func (s *BffService) GetUserCommentTalkReplyList(ctx context.Context, req *v1.GetUserCommentTalkReplyListReq) (*v1.GetUserCommentTalkReplyListReply, error) {
	reply := &v1.GetUserCommentTalkReplyListReply{List: make([]*v1.GetUserCommentTalkReplyListReply_List, 0)}
	commentList, err := s.commc.GetUserCommentTalkReplyList(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range commentList {
		reply.List = append(reply.List, &v1.GetUserCommentTalkReplyListReply_List{
			Id:             item.Id,
			CreationId:     item.CreationId,
			CreationAuthor: item.CreationAuthor,
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

func (s *BffService) SendSubComment(ctx context.Context, req *v1.SendSubCommentReq) (*emptypb.Empty, error) {
	err := s.commc.SendSubComment(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) RemoveComment(ctx context.Context, req *v1.RemoveCommentReq) (*emptypb.Empty, error) {
	err := s.commc.RemoveComment(ctx, req.Id, req.CreationId, req.CreationType, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) RemoveSubComment(ctx context.Context, req *v1.RemoveSubCommentReq) (*emptypb.Empty, error) {
	err := s.commc.RemoveSubComment(ctx, req.Id, req.RootId, req.Uuid, req.Reply)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SetCommentAgree(ctx context.Context, req *v1.SetCommentAgreeReq) (*emptypb.Empty, error) {
	err := s.commc.SetCommentAgree(ctx, req.Id, req.CreationId, req.CreationType, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SetSubCommentAgree(ctx context.Context, req *v1.SetSubCommentAgreeReq) (*emptypb.Empty, error) {
	err := s.commc.SetSubCommentAgree(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) CancelCommentAgree(ctx context.Context, req *v1.CancelCommentAgreeReq) (*emptypb.Empty, error) {
	err := s.commc.CancelCommentAgree(ctx, req.Id, req.CreationId, req.CreationType, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) CancelSubCommentAgree(ctx context.Context, req *v1.CancelSubCommentAgreeReq) (*emptypb.Empty, error) {
	err := s.commc.CancelSubCommentAgree(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
