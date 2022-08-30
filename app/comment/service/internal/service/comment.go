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

func (s *CommentService) GetUserCommentAgree(ctx context.Context, req *v1.GetUserCommentAgreeReq) (*v1.GetUserCommentAgreeReply, error) {
	agreeMap, err := s.cc.GetUserCommentAgree(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserCommentAgreeReply{
		Agree: agreeMap,
	}, nil
}

func (s *CommentService) GetCommentUser(ctx context.Context, req *v1.GetCommentUserReq) (*v1.GetCommentUserReply, error) {
	commentUser, err := s.cc.GetCommentUser(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetCommentUserReply{
		Comment: commentUser.Comment,
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

func (s *CommentService) GetSubCommentList(ctx context.Context, req *v1.GetSubCommentListReq) (*v1.GetSubCommentListReply, error) {
	reply := &v1.GetSubCommentListReply{Comment: make([]*v1.GetSubCommentListReply_Comment, 0)}
	subCommentList, err := s.cc.GetSubCommentList(ctx, req.Page, req.Id)
	if err != nil {
		return nil, err
	}
	for _, item := range subCommentList {
		reply.Comment = append(reply.Comment, &v1.GetSubCommentListReply_Comment{
			Id:    item.CommentId,
			Uuid:  item.Uuid,
			Reply: item.Reply,
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

func (s *CommentService) GetSubCommentListStatistic(ctx context.Context, req *v1.GetCommentListStatisticReq) (*v1.GetCommentListStatisticReply, error) {
	reply := &v1.GetCommentListStatisticReply{Count: make([]*v1.GetCommentListStatisticReply_Count, 0)}
	commentListStatistic, err := s.cc.GetSubCommentListStatistic(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	for _, item := range commentListStatistic {
		reply.Count = append(reply.Count, &v1.GetCommentListStatisticReply_Count{
			Id:    item.CommentId,
			Agree: item.Agree,
		})
	}
	return reply, nil
}

func (s *CommentService) GetUserCommentArticleReplyList(ctx context.Context, req *v1.GetUserCommentArticleReplyListReq) (*v1.GetUserCommentArticleReplyListReply, error) {
	reply := &v1.GetUserCommentArticleReplyListReply{List: make([]*v1.GetUserCommentArticleReplyListReply_List, 0)}
	commentList, err := s.cc.GetUserCommentArticleReplyList(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range commentList {
		reply.List = append(reply.List, &v1.GetUserCommentArticleReplyListReply_List{
			Id:             item.CommentId,
			CreationId:     item.CreationId,
			CreationAuthor: item.CreationAuthor,
		})
	}
	return reply, nil
}

func (s *CommentService) GetUserCommentTalkReplyList(ctx context.Context, req *v1.GetUserCommentTalkReplyListReq) (*v1.GetUserCommentTalkReplyListReply, error) {
	reply := &v1.GetUserCommentTalkReplyListReply{List: make([]*v1.GetUserCommentTalkReplyListReply_List, 0)}
	commentList, err := s.cc.GetUserCommentTalkReplyList(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range commentList {
		reply.List = append(reply.List, &v1.GetUserCommentTalkReplyListReply_List{
			Id:             item.CommentId,
			CreationId:     item.CreationId,
			CreationAuthor: item.CreationAuthor,
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

func (s *CommentService) CreateSubComment(ctx context.Context, req *v1.CreateSubCommentReq) (*emptypb.Empty, error) {
	err := s.cc.CreateSubComment(ctx, req.Id, req.RootId, req.ParentId, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommentService) CreateCommentDbAndCache(ctx context.Context, req *v1.CreateCommentDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.cc.CreateCommentDbAndCache(ctx, req.Id, req.CreationId, req.CreationType, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommentService) CreateSubCommentDbAndCache(ctx context.Context, req *v1.CreateSubCommentDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.cc.CreateSubCommentDbAndCache(ctx, req.Id, req.RootId, req.ParentId, req.Uuid)
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

func (s *CommentService) SendSubComment(ctx context.Context, req *v1.SendSubCommentReq) (*emptypb.Empty, error) {
	err := s.cc.SendSubComment(ctx, req.Id, req.Uuid, req.Ip)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommentService) RemoveComment(ctx context.Context, req *v1.RemoveCommentReq) (*emptypb.Empty, error) {
	err := s.cc.RemoveComment(ctx, req.Id, req.CreationId, req.CreationType, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommentService) RemoveSubComment(ctx context.Context, req *v1.RemoveSubCommentReq) (*emptypb.Empty, error) {
	err := s.cc.RemoveSubComment(ctx, req.Id, req.RootId, req.Uuid, req.UserUuid, req.Reply)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommentService) RemoveCommentDbAndCache(ctx context.Context, req *v1.RemoveCommentDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.cc.RemoveCommentDbAndCache(ctx, req.Id, req.CreationId, req.CreationType, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommentService) RemoveSubCommentDbAndCache(ctx context.Context, req *v1.RemoveSubCommentDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.cc.RemoveSubCommentDbAndCache(ctx, req.Id, req.RootId, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommentService) SetCommentAgree(ctx context.Context, req *v1.SetCommentAgreeReq) (*emptypb.Empty, error) {
	err := s.cc.SetCommentAgree(ctx, req.Id, req.CreationId, req.CreationType, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommentService) SetSubCommentAgree(ctx context.Context, req *v1.SetSubCommentAgreeReq) (*emptypb.Empty, error) {
	err := s.cc.SetSubCommentAgree(ctx, req.Id, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommentService) SetCommentAgreeDbAndCache(ctx context.Context, req *v1.SetCommentAgreeReq) (*emptypb.Empty, error) {
	err := s.cc.SetCommentAgreeDbAndCache(ctx, req.Id, req.CreationId, req.CreationType, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommentService) SetSubCommentAgreeDbAndCache(ctx context.Context, req *v1.SetSubCommentAgreeReq) (*emptypb.Empty, error) {
	err := s.cc.SetSubCommentAgreeDbAndCache(ctx, req.Id, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommentService) CancelCommentAgree(ctx context.Context, req *v1.CancelCommentAgreeReq) (*emptypb.Empty, error) {
	err := s.cc.CancelCommentAgree(ctx, req.Id, req.CreationId, req.CreationType, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommentService) CancelSubCommentAgree(ctx context.Context, req *v1.CancelSubCommentAgreeReq) (*emptypb.Empty, error) {
	err := s.cc.CancelSubCommentAgree(ctx, req.Id, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommentService) CancelCommentAgreeDbAndCache(ctx context.Context, req *v1.CancelCommentAgreeReq) (*emptypb.Empty, error) {
	err := s.cc.CancelCommentAgreeDbAndCache(ctx, req.Id, req.CreationId, req.CreationType, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommentService) CancelSubCommentAgreeDbAndCache(ctx context.Context, req *v1.CancelSubCommentAgreeReq) (*emptypb.Empty, error) {
	err := s.cc.CancelSubCommentAgreeDbAndCache(ctx, req.Id, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
