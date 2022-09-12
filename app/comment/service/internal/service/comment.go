package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/comment/service/v1"
	"github.com/the-zion/matrix-core/app/comment/service/internal/biz"
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
		Comment:           commentUser.Comment,
		ArticleReply:      commentUser.ArticleReply,
		ArticleReplySub:   commentUser.ArticleReplySub,
		TalkReply:         commentUser.TalkReply,
		TalkReplySub:      commentUser.TalkReplySub,
		ArticleReplied:    commentUser.ArticleReplied,
		ArticleRepliedSub: commentUser.ArticleRepliedSub,
		TalkReplied:       commentUser.TalkReplied,
		TalkRepliedSub:    commentUser.TalkRepliedSub,
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

func (s *CommentService) GetUserSubCommentArticleReplyList(ctx context.Context, req *v1.GetUserSubCommentArticleReplyListReq) (*v1.GetUserSubCommentArticleReplyListReply, error) {
	reply := &v1.GetUserSubCommentArticleReplyListReply{List: make([]*v1.GetUserSubCommentArticleReplyListReply_List, 0)}
	commentList, err := s.cc.GetUserSubCommentArticleReplyList(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range commentList {
		reply.List = append(reply.List, &v1.GetUserSubCommentArticleReplyListReply_List{
			Id:             item.CommentId,
			CreationId:     item.CreationId,
			RootId:         item.RootId,
			ParentId:       item.ParentId,
			CreationAuthor: item.CreationAuthor,
			RootUser:       item.RootUser,
			Reply:          item.Reply,
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

func (s *CommentService) GetUserSubCommentTalkReplyList(ctx context.Context, req *v1.GetUserSubCommentTalkReplyListReq) (*v1.GetUserSubCommentTalkReplyListReply, error) {
	reply := &v1.GetUserSubCommentTalkReplyListReply{List: make([]*v1.GetUserSubCommentTalkReplyListReply_List, 0)}
	commentList, err := s.cc.GetUserSubCommentTalkReplyList(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range commentList {
		reply.List = append(reply.List, &v1.GetUserSubCommentTalkReplyListReply_List{
			Id:             item.CommentId,
			CreationId:     item.CreationId,
			RootId:         item.RootId,
			ParentId:       item.ParentId,
			CreationAuthor: item.CreationAuthor,
			RootUser:       item.RootUser,
			Reply:          item.Reply,
		})
	}
	return reply, nil
}

func (s *CommentService) GetUserCommentArticleRepliedList(ctx context.Context, req *v1.GetUserCommentArticleRepliedListReq) (*v1.GetUserCommentArticleRepliedListReply, error) {
	reply := &v1.GetUserCommentArticleRepliedListReply{List: make([]*v1.GetUserCommentArticleRepliedListReply_List, 0)}
	commentList, err := s.cc.GetUserCommentArticleRepliedList(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range commentList {
		reply.List = append(reply.List, &v1.GetUserCommentArticleRepliedListReply_List{
			Id:         item.CommentId,
			CreationId: item.CreationId,
			Uuid:       item.Uuid,
		})
	}
	return reply, nil
}

func (s *CommentService) GetUserSubCommentArticleRepliedList(ctx context.Context, req *v1.GetUserSubCommentArticleRepliedListReq) (*v1.GetUserSubCommentArticleRepliedListReply, error) {
	reply := &v1.GetUserSubCommentArticleRepliedListReply{List: make([]*v1.GetUserSubCommentArticleRepliedListReply_List, 0)}
	commentList, err := s.cc.GetUserSubCommentArticleRepliedList(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range commentList {
		reply.List = append(reply.List, &v1.GetUserSubCommentArticleRepliedListReply_List{
			Id:             item.CommentId,
			Uuid:           item.Uuid,
			CreationId:     item.CreationId,
			RootId:         item.RootId,
			ParentId:       item.ParentId,
			CreationAuthor: item.CreationAuthor,
			RootUser:       item.RootUser,
			Reply:          item.Reply,
			ReplyName:      item.ReplyName,
			RootName:       item.RootName,
			UserName:       item.UserName,
		})
	}
	return reply, nil
}

func (s *CommentService) GetUserCommentTalkRepliedList(ctx context.Context, req *v1.GetUserCommentTalkRepliedListReq) (*v1.GetUserCommentTalkRepliedListReply, error) {
	reply := &v1.GetUserCommentTalkRepliedListReply{List: make([]*v1.GetUserCommentTalkRepliedListReply_List, 0)}
	commentList, err := s.cc.GetUserCommentTalkRepliedList(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range commentList {
		reply.List = append(reply.List, &v1.GetUserCommentTalkRepliedListReply_List{
			Id:         item.CommentId,
			CreationId: item.CreationId,
			Uuid:       item.Uuid,
		})
	}
	return reply, nil
}

func (s *CommentService) GetUserSubCommentTalkRepliedList(ctx context.Context, req *v1.GetUserSubCommentTalkRepliedListReq) (*v1.GetUserSubCommentTalkRepliedListReply, error) {
	reply := &v1.GetUserSubCommentTalkRepliedListReply{List: make([]*v1.GetUserSubCommentTalkRepliedListReply_List, 0)}
	commentList, err := s.cc.GetUserSubCommentTalkRepliedList(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range commentList {
		reply.List = append(reply.List, &v1.GetUserSubCommentTalkRepliedListReply_List{
			Id:             item.CommentId,
			Uuid:           item.Uuid,
			CreationId:     item.CreationId,
			RootId:         item.RootId,
			ParentId:       item.ParentId,
			CreationAuthor: item.CreationAuthor,
			RootUser:       item.RootUser,
			Reply:          item.Reply,
			ReplyName:      item.ReplyName,
			RootName:       item.RootName,
			UserName:       item.UserName,
		})
	}
	return reply, nil
}

func (s *CommentService) GetCommentContentReview(ctx context.Context, req *v1.GetCommentContentReviewReq) (*v1.GetCommentContentReviewReply, error) {
	reply := &v1.GetCommentContentReviewReply{Review: make([]*v1.GetCommentContentReviewReply_Review, 0)}
	reviewList, err := s.cc.GetCommentContentReview(ctx, req.Page, req.Uuid)
	if err != nil {
		return reply, err
	}
	for _, item := range reviewList {
		reply.Review = append(reply.Review, &v1.GetCommentContentReviewReply_Review{
			Id:        item.Id,
			CommentId: item.CommentId,
			Comment:   item.Comment,
			Kind:      item.Kind,
			Uuid:      item.Uuid,
			CreateAt:  item.CreateAt,
			JobId:     item.JobId,
			Label:     item.Label,
			Result:    item.Result,
			Section:   item.Section,
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

func (s *CommentService) CommentContentIrregular(ctx context.Context, req *v1.CommentContentIrregularReq) (*emptypb.Empty, error) {
	err := s.cc.CommentContentIrregular(ctx, &biz.TextReview{
		CommentId: req.Id,
		Uuid:      req.Uuid,
		JobId:     req.JobId,
		Comment:   req.Comment,
		Kind:      req.Kind,
		Label:     req.Label,
		Result:    req.Result,
		Section:   req.Section,
		Mode:      "add_comment_content_review_db_and_cache",
	})
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
	err := s.cc.RemoveComment(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommentService) RemoveSubComment(ctx context.Context, req *v1.RemoveSubCommentReq) (*emptypb.Empty, error) {
	err := s.cc.RemoveSubComment(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommentService) RemoveCommentDbAndCache(ctx context.Context, req *v1.RemoveCommentDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.cc.RemoveCommentDbAndCache(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommentService) RemoveSubCommentDbAndCache(ctx context.Context, req *v1.RemoveSubCommentDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.cc.RemoveSubCommentDbAndCache(ctx, req.Id, req.Uuid)
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

func (s *CommentService) AddCommentContentReviewDbAndCache(ctx context.Context, req *v1.AddCommentContentReviewDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.cc.AddCommentContentReviewDbAndCache(ctx, &biz.TextReview{
		CommentId: req.CommentId,
		Uuid:      req.Uuid,
		JobId:     req.JobId,
		Comment:   req.Comment,
		Kind:      req.Kind,
		Label:     req.Label,
		Result:    req.Result,
		Section:   req.Section,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
