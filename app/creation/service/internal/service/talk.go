package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
	"github.com/the-zion/matrix-core/app/creation/service/internal/biz"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *CreationService) GetLastTalkDraft(ctx context.Context, req *v1.GetLastTalkDraftReq) (*v1.GetLastTalkDraftReply, error) {
	draft, err := s.tc.GetLastTalkDraft(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetLastTalkDraftReply{
		Id:     draft.Id,
		Status: draft.Status,
	}, nil
}

func (s *CreationService) GetTalkList(ctx context.Context, req *v1.GetTalkListReq) (*v1.GetTalkListReply, error) {
	reply := &v1.GetTalkListReply{Talk: make([]*v1.GetTalkListReply_Talk, 0)}
	talkList, err := s.tc.GetTalkList(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range talkList {
		reply.Talk = append(reply.Talk, &v1.GetTalkListReply_Talk{
			Id:   item.TalkId,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *CreationService) GetTalkListHot(ctx context.Context, req *v1.GetTalkListHotReq) (*v1.GetTalkListHotReply, error) {
	reply := &v1.GetTalkListHotReply{Talk: make([]*v1.GetTalkListHotReply_Talk, 0)}
	talkList, err := s.tc.GetTalkListHot(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range talkList {
		reply.Talk = append(reply.Talk, &v1.GetTalkListHotReply_Talk{
			Id:   item.TalkId,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *CreationService) GetUserTalkList(ctx context.Context, req *v1.GetUserTalkListReq) (*v1.GetTalkListReply, error) {
	reply := &v1.GetTalkListReply{Talk: make([]*v1.GetTalkListReply_Talk, 0)}
	talkList, err := s.tc.GetUserTalkList(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range talkList {
		reply.Talk = append(reply.Talk, &v1.GetTalkListReply_Talk{
			Id:   item.TalkId,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *CreationService) GetUserTalkListVisitor(ctx context.Context, req *v1.GetUserTalkListVisitorReq) (*v1.GetTalkListReply, error) {
	reply := &v1.GetTalkListReply{Talk: make([]*v1.GetTalkListReply_Talk, 0)}
	talkList, err := s.tc.GetUserTalkListVisitor(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range talkList {
		reply.Talk = append(reply.Talk, &v1.GetTalkListReply_Talk{
			Id:   item.TalkId,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *CreationService) GetTalkCount(ctx context.Context, req *v1.GetTalkCountReq) (*v1.GetTalkCountReply, error) {
	count, err := s.tc.GetTalkCount(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetTalkCountReply{
		Count: count,
	}, nil
}

func (s *CreationService) GetTalkCountVisitor(ctx context.Context, req *v1.GetTalkCountVisitorReq) (*v1.GetTalkCountReply, error) {
	count, err := s.tc.GetTalkCountVisitor(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetTalkCountReply{
		Count: count,
	}, nil
}

func (s *CreationService) GetTalkListStatistic(ctx context.Context, req *v1.GetTalkListStatisticReq) (*v1.GetTalkListStatisticReply, error) {
	reply := &v1.GetTalkListStatisticReply{Count: make([]*v1.GetTalkListStatisticReply_Count, 0)}
	statisticList, err := s.tc.GetTalkListStatistic(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	for _, item := range statisticList {
		reply.Count = append(reply.Count, &v1.GetTalkListStatisticReply_Count{
			Id:      item.TalkId,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
			Comment: item.Comment,
		})
	}
	return reply, nil
}

func (s *CreationService) GetTalkStatistic(ctx context.Context, req *v1.GetTalkStatisticReq) (*v1.GetTalkStatisticReply, error) {
	talkStatistic, err := s.tc.GetTalkStatistic(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetTalkStatisticReply{
		Uuid:    talkStatistic.Uuid,
		Agree:   talkStatistic.Agree,
		Collect: talkStatistic.Collect,
		View:    talkStatistic.View,
		Comment: talkStatistic.Comment,
	}, nil
}

func (s *CreationService) GetTalkSearch(ctx context.Context, req *v1.GetTalkSearchReq) (*v1.GetTalkSearchReply, error) {
	reply := &v1.GetTalkSearchReply{List: make([]*v1.GetTalkSearchReply_List, 0)}
	talkList, total, err := s.tc.GetTalkSearch(ctx, req.Page, req.Search, req.Time)
	if err != nil {
		return reply, err
	}
	for _, item := range talkList {
		reply.List = append(reply.List, &v1.GetTalkSearchReply_List{
			Id:     item.Id,
			Tags:   item.Tags,
			Title:  item.Title,
			Uuid:   item.Uuid,
			Text:   item.Text,
			Cover:  item.Cover,
			Update: item.Update,
		})
	}
	reply.Total = total
	return reply, nil
}

func (s *CreationService) GetUserTalkAgree(ctx context.Context, req *v1.GetUserTalkAgreeReq) (*v1.GetUserTalkAgreeReply, error) {
	agreeMap, err := s.tc.GetUserTalkAgree(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserTalkAgreeReply{
		Agree: agreeMap,
	}, nil
}

func (s *CreationService) GetUserTalkCollect(ctx context.Context, req *v1.GetUserTalkCollectReq) (*v1.GetUserTalkCollectReply, error) {
	collectMap, err := s.tc.GetUserTalkCollect(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserTalkCollectReply{
		Collect: collectMap,
	}, nil
}

func (s *CreationService) GetTalkImageReview(ctx context.Context, req *v1.GetTalkImageReviewReq) (*v1.GetTalkImageReviewReply, error) {
	reply := &v1.GetTalkImageReviewReply{Review: make([]*v1.GetTalkImageReviewReply_Review, 0)}
	reviewList, err := s.tc.GetTalkImageReview(ctx, req.Page, req.Uuid)
	if err != nil {
		return reply, err
	}
	for _, item := range reviewList {
		reply.Review = append(reply.Review, &v1.GetTalkImageReviewReply_Review{
			Id:         item.Id,
			CreationId: item.CreationId,
			Kind:       item.Kind,
			Uid:        item.Uid,
			Uuid:       item.Uuid,
			CreateAt:   item.CreateAt,
			JobId:      item.JobId,
			Url:        item.Url,
			Label:      item.Label,
			Result:     item.Result,
			Score:      item.Score,
			Category:   item.Category,
			SubLabel:   item.SubLabel,
		})
	}
	return reply, nil
}

func (s *CreationService) TalkImageIrregular(ctx context.Context, req *v1.CreationImageIrregularReq) (*emptypb.Empty, error) {
	err := s.tc.TalkImageIrregular(ctx, &biz.ImageReview{
		CreationId: req.Id,
		Kind:       req.Kind,
		Uid:        req.Uid,
		Uuid:       req.Uuid,
		JobId:      req.JobId,
		Url:        req.Url,
		Label:      req.Label,
		Result:     req.Result,
		Score:      req.Score,
		Category:   req.Category,
		SubLabel:   req.SubLabel,
		Mode:       "add_talk_image_review_db_and_cache",
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) AddTalkImageReviewDbAndCache(ctx context.Context, req *v1.AddCreationImageReviewDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.tc.AddTalkImageReviewDbAndCache(ctx, &biz.ImageReview{
		CreationId: req.CreationId,
		Kind:       req.Kind,
		Uid:        req.Uid,
		Uuid:       req.Uuid,
		JobId:      req.JobId,
		Url:        req.Url,
		Label:      req.Label,
		Result:     req.Result,
		Score:      req.Score,
		Category:   req.Category,
		SubLabel:   req.SubLabel,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CreateTalkDraft(ctx context.Context, req *v1.CreateTalkDraftReq) (*v1.CreateTalkDraftReply, error) {
	id, err := s.tc.CreateTalkDraft(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.CreateTalkDraftReply{
		Id: id,
	}, nil
}

func (s *CreationService) SendTalk(ctx context.Context, req *v1.SendTalkReq) (*emptypb.Empty, error) {
	err := s.tc.SendTalk(ctx, req.Id, req.Uuid, req.Ip)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SendTalkEdit(ctx context.Context, req *v1.SendTalkEditReq) (*emptypb.Empty, error) {
	err := s.tc.SendTalkEdit(ctx, req.Id, req.Uuid, req.Ip)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) DeleteTalk(ctx context.Context, req *v1.DeleteTalkReq) (*emptypb.Empty, error) {
	err := s.tc.DeleteTalk(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CreateTalk(ctx context.Context, req *v1.CreateTalkReq) (*emptypb.Empty, error) {
	err := s.tc.CreateTalk(ctx, req.Id, req.Auth, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) EditTalk(ctx context.Context, req *v1.EditTalkReq) (*emptypb.Empty, error) {
	err := s.tc.EditTalk(ctx, req.Id, req.Auth, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CreateTalkDbCacheAndSearch(ctx context.Context, req *v1.CreateTalkDbCacheAndSearchReq) (*emptypb.Empty, error) {
	err := s.tc.CreateTalkDbCacheAndSearch(ctx, req.Id, req.Auth, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) EditTalkCosAndSearch(ctx context.Context, req *v1.EditTalkCosAndSearchReq) (*emptypb.Empty, error) {
	err := s.tc.EditTalkCosAndSearch(ctx, req.Id, req.Auth, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) DeleteTalkCacheAndSearch(ctx context.Context, req *v1.DeleteTalkCacheAndSearchReq) (*emptypb.Empty, error) {
	err := s.tc.DeleteTalkCacheAndSearch(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SetTalkAgree(ctx context.Context, req *v1.SetTalkAgreeReq) (*emptypb.Empty, error) {
	err := s.tc.SetTalkAgree(ctx, req.Id, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SetTalkAgreeDbAndCache(ctx context.Context, req *v1.SetTalkAgreeDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.tc.SetTalkAgreeDbAndCache(ctx, req.Id, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CancelTalkAgree(ctx context.Context, req *v1.CancelTalkAgreeReq) (*emptypb.Empty, error) {
	err := s.tc.CancelTalkAgree(ctx, req.Id, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CancelTalkAgreeDbAndCache(ctx context.Context, req *v1.CancelTalkAgreeDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.tc.CancelTalkAgreeDbAndCache(ctx, req.Id, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CancelTalkCollect(ctx context.Context, req *v1.CancelTalkCollectReq) (*emptypb.Empty, error) {
	err := s.tc.CancelTalkCollect(ctx, req.Id, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CancelTalkCollectDbAndCache(ctx context.Context, req *v1.CancelTalkCollectDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.tc.CancelTalkCollectDbAndCache(ctx, req.Id, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SetTalkView(ctx context.Context, req *v1.SetTalkViewReq) (*emptypb.Empty, error) {
	err := s.tc.SetTalkView(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SetTalkViewDbAndCache(ctx context.Context, req *v1.SetTalkViewDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.tc.SetTalkViewDbAndCache(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SetTalkCollect(ctx context.Context, req *v1.SetTalkCollectReq) (*emptypb.Empty, error) {
	err := s.tc.SetTalkCollect(ctx, req.Id, req.CollectionsId, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SetTalkCollectDbAndCache(ctx context.Context, req *v1.SetTalkCollectDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.tc.SetTalkCollectDbAndCache(ctx, req.Id, req.CollectionsId, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) TalkStatisticJudge(ctx context.Context, req *v1.TalkStatisticJudgeReq) (*v1.TalkStatisticJudgeReply, error) {
	judge, err := s.tc.TalkStatisticJudge(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.TalkStatisticJudgeReply{
		Agree:   judge.Agree,
		Collect: judge.Collect,
	}, nil
}
