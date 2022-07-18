package service

import (
	"context"
	"github.com/the-zion/matrix-core/api/bff/interface/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ------------------------------------------creation-------------------------------------------------

func (s *BffService) GetLeaderBoard(ctx context.Context, _ *emptypb.Empty) (*v1.GetLeaderBoardReply, error) {
	reply := &v1.GetLeaderBoardReply{Board: make([]*v1.GetLeaderBoardReply_Board, 0)}
	boardList, err := s.cc.GetLeaderBoard(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range boardList {
		reply.Board = append(reply.Board, &v1.GetLeaderBoardReply_Board{
			Id:   item.Id,
			Uuid: item.Uuid,
			Mode: item.Mode,
		})
	}
	return reply, nil
}

func (s *BffService) GetCollectArticle(ctx context.Context, req *v1.GetCollectArticleReq) (*v1.GetArticleListReply, error) {
	reply := &v1.GetArticleListReply{Article: make([]*v1.GetArticleListReply_Article, 0)}
	articleList, err := s.cc.GetCollectArticle(ctx, req.Id, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range articleList {
		reply.Article = append(reply.Article, &v1.GetArticleListReply_Article{
			Id:   item.Id,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *BffService) GetCollectArticleCount(ctx context.Context, req *v1.GetCollectArticleCountReq) (*v1.GetCollectArticleCountReply, error) {
	count, err := s.cc.GetCollectArticleCount(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetCollectArticleCountReply{
		Count: count,
	}, nil
}

func (s *BffService) GetCollectTalk(ctx context.Context, req *v1.GetCollectTalkReq) (*v1.GetTalkListReply, error) {
	reply := &v1.GetTalkListReply{Talk: make([]*v1.GetTalkListReply_Talk, 0)}
	talkList, err := s.cc.GetCollectTalk(ctx, req.Id, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range talkList {
		reply.Talk = append(reply.Talk, &v1.GetTalkListReply_Talk{
			Id:   item.Id,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *BffService) GetCollectTalkCount(ctx context.Context, req *v1.GetCollectTalkCountReq) (*v1.GetCollectTalkCountReply, error) {
	count, err := s.cc.GetCollectTalkCount(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetCollectTalkCountReply{
		Count: count,
	}, nil
}

func (s *BffService) GetCollection(ctx context.Context, req *v1.GetCollectionReq) (*v1.GetCollectionReply, error) {
	collection, err := s.cc.GetCollection(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetCollectionReply{
		Uuid:      collection.Uuid,
		Name:      collection.Name,
		Introduce: collection.Introduce,
		Auth:      collection.Auth,
	}, nil
}

func (s *BffService) GetCollections(ctx context.Context, req *v1.GetCollectionsReq) (*v1.GetCollectionsReply, error) {
	reply := &v1.GetCollectionsReply{Collections: make([]*v1.GetCollectionsReply_Collections, 0)}
	collections, err := s.cc.GetCollections(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range collections {
		reply.Collections = append(reply.Collections, &v1.GetCollectionsReply_Collections{
			Id:        item.Id,
			Name:      item.Name,
			Introduce: item.Introduce,
		})
	}
	return reply, nil
}

func (s *BffService) GetCollectionsCount(ctx context.Context, _ *emptypb.Empty) (*v1.GetCollectionsCountReply, error) {
	count, err := s.cc.GetCollectionsCount(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetCollectionsCountReply{
		Count: count,
	}, nil
}

func (s *BffService) GetCollectionsByVisitor(ctx context.Context, req *v1.GetCollectionsByVisitorReq) (*v1.GetCollectionsReply, error) {
	reply := &v1.GetCollectionsReply{Collections: make([]*v1.GetCollectionsReply_Collections, 0)}
	collections, err := s.cc.GetCollectionsByVisitor(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range collections {
		reply.Collections = append(reply.Collections, &v1.GetCollectionsReply_Collections{
			Id:        item.Id,
			Name:      item.Name,
			Introduce: item.Introduce,
		})
	}
	return reply, nil
}

func (s *BffService) GetCollectionsVisitorCount(ctx context.Context, req *v1.GetCollectionsVisitorCountReq) (*v1.GetCollectionsCountReply, error) {
	count, err := s.cc.GetCollectionsVisitorCount(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetCollectionsCountReply{
		Count: count,
	}, nil
}

func (s *BffService) CreateCollections(ctx context.Context, req *v1.CreateCollectionsReq) (*emptypb.Empty, error) {
	err := s.cc.CreateCollections(ctx, req.Name, req.Introduce, req.Auth)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) EditCollections(ctx context.Context, req *v1.EditCollectionsReq) (*emptypb.Empty, error) {
	err := s.cc.EditCollections(ctx, req.Id, req.Name, req.Introduce, req.Auth)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) DeleteCollections(ctx context.Context, req *v1.DeleteCollectionsReq) (*emptypb.Empty, error) {
	err := s.cc.DeleteCollections(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ------------------------------------------article-------------------------------------------------

func (s *BffService) GetArticleList(ctx context.Context, req *v1.GetArticleListReq) (*v1.GetArticleListReply, error) {
	reply := &v1.GetArticleListReply{Article: make([]*v1.GetArticleListReply_Article, 0)}
	articleList, err := s.ac.GetArticleList(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range articleList {
		reply.Article = append(reply.Article, &v1.GetArticleListReply_Article{
			Id:   item.Id,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *BffService) GetArticleListHot(ctx context.Context, req *v1.GetArticleListHotReq) (*v1.GetArticleListHotReply, error) {
	reply := &v1.GetArticleListHotReply{Article: make([]*v1.GetArticleListHotReply_Article, 0)}
	articleList, err := s.ac.GetArticleListHot(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range articleList {
		reply.Article = append(reply.Article, &v1.GetArticleListHotReply_Article{
			Id:   item.Id,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *BffService) GetArticleCountVisitor(ctx context.Context, req *v1.GetArticleCountVisitorReq) (*v1.GetArticleCountReply, error) {
	count, err := s.ac.GetArticleCountVisitor(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetArticleCountReply{
		Count: count,
	}, nil
}

func (s *BffService) GetArticleCount(ctx context.Context, _ *emptypb.Empty) (*v1.GetArticleCountReply, error) {
	count, err := s.ac.GetArticleCount(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetArticleCountReply{
		Count: count,
	}, nil
}

func (s *BffService) GetUserArticleList(ctx context.Context, req *v1.GetUserArticleListReq) (*v1.GetArticleListReply, error) {
	reply := &v1.GetArticleListReply{Article: make([]*v1.GetArticleListReply_Article, 0)}
	articleList, err := s.ac.GetUserArticleList(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range articleList {
		reply.Article = append(reply.Article, &v1.GetArticleListReply_Article{
			Id:   item.Id,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *BffService) GetUserArticleListVisitor(ctx context.Context, req *v1.GetUserArticleListVisitorReq) (*v1.GetArticleListReply, error) {
	reply := &v1.GetArticleListReply{Article: make([]*v1.GetArticleListReply_Article, 0)}
	articleList, err := s.ac.GetUserArticleListVisitor(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range articleList {
		reply.Article = append(reply.Article, &v1.GetArticleListReply_Article{
			Id:   item.Id,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *BffService) GetArticleStatistic(ctx context.Context, req *v1.GetArticleStatisticReq) (*v1.GetArticleStatisticReply, error) {
	articleStatistic, err := s.ac.GetArticleStatistic(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetArticleStatisticReply{
		Uuid:    articleStatistic.Uuid,
		Agree:   articleStatistic.Agree,
		Collect: articleStatistic.Collect,
		View:    articleStatistic.View,
		Comment: articleStatistic.Comment,
	}, nil
}

func (s *BffService) GetArticleListStatistic(ctx context.Context, req *v1.GetArticleListStatisticReq) (*v1.GetArticleListStatisticReply, error) {
	reply := &v1.GetArticleListStatisticReply{Count: make([]*v1.GetArticleListStatisticReply_Count, 0)}
	statisticList, err := s.ac.GetArticleListStatistic(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	for _, item := range statisticList {
		reply.Count = append(reply.Count, &v1.GetArticleListStatisticReply_Count{
			Id:      item.Id,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
			Comment: item.Comment,
		})
	}
	return reply, nil
}

func (s *BffService) GetLastArticleDraft(ctx context.Context, _ *emptypb.Empty) (*v1.GetLastArticleDraftReply, error) {
	draft, err := s.ac.GetLastArticleDraft(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetLastArticleDraftReply{
		Id:     draft.Id,
		Status: draft.Status,
	}, nil
}

func (s *BffService) CreateArticleDraft(ctx context.Context, _ *emptypb.Empty) (*v1.CreateArticleDraftReply, error) {
	id, err := s.ac.CreateArticleDraft(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.CreateArticleDraftReply{
		Id: id,
	}, nil
}

func (s *BffService) ArticleDraftMark(ctx context.Context, req *v1.ArticleDraftMarkReq) (*emptypb.Empty, error) {
	err := s.ac.ArticleDraftMark(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) GetArticleDraftList(ctx context.Context, _ *emptypb.Empty) (*v1.GetArticleDraftListReply, error) {
	reply := &v1.GetArticleDraftListReply{Draft: make([]*v1.GetArticleDraftListReply_Draft, 0)}
	draftList, err := s.ac.GetArticleDraftList(ctx)
	if err != nil {
		return reply, err
	}
	for _, item := range draftList {
		reply.Draft = append(reply.Draft, &v1.GetArticleDraftListReply_Draft{
			Id: item.Id,
		})
	}
	return reply, nil
}

func (s *BffService) SendArticle(ctx context.Context, req *v1.SendArticleReq) (*emptypb.Empty, error) {
	err := s.ac.SendArticle(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SendArticleEdit(ctx context.Context, req *v1.SendArticleEditReq) (*emptypb.Empty, error) {
	err := s.ac.SendArticleEdit(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) DeleteArticle(ctx context.Context, req *v1.DeleteArticleReq) (*emptypb.Empty, error) {
	err := s.ac.DeleteArticle(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SetArticleAgree(ctx context.Context, req *v1.SetArticleAgreeReq) (*emptypb.Empty, error) {
	err := s.ac.SetArticleAgree(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SetArticleView(ctx context.Context, req *v1.SetArticleViewReq) (*emptypb.Empty, error) {
	err := s.ac.SetArticleView(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SetArticleCollect(ctx context.Context, req *v1.SetArticleCollectReq) (*emptypb.Empty, error) {
	err := s.ac.SetArticleCollect(ctx, req.Id, req.CollectionsId, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) CancelArticleAgree(ctx context.Context, req *v1.CancelArticleAgreeReq) (*emptypb.Empty, error) {
	err := s.ac.CancelArticleAgree(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) CancelArticleCollect(ctx context.Context, req *v1.CancelArticleCollectReq) (*emptypb.Empty, error) {
	err := s.ac.CancelArticleCollect(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) ArticleStatisticJudge(ctx context.Context, req *v1.ArticleStatisticJudgeReq) (*v1.ArticleStatisticJudgeReply, error) {
	judge, err := s.ac.ArticleStatisticJudge(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.ArticleStatisticJudgeReply{
		Agree:   judge.Agree,
		Collect: judge.Collect,
	}, nil
}

// ------------------------------------------talk-------------------------------------------------

func (s *BffService) GetTalkList(ctx context.Context, req *v1.GetTalkListReq) (*v1.GetTalkListReply, error) {
	reply := &v1.GetTalkListReply{Talk: make([]*v1.GetTalkListReply_Talk, 0)}
	talkList, err := s.tc.GetTalkList(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range talkList {
		reply.Talk = append(reply.Talk, &v1.GetTalkListReply_Talk{
			Id:   item.Id,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *BffService) GetTalkListHot(ctx context.Context, req *v1.GetTalkListHotReq) (*v1.GetTalkListHotReply, error) {
	reply := &v1.GetTalkListHotReply{Talk: make([]*v1.GetTalkListHotReply_Talk, 0)}
	talkList, err := s.tc.GetTalkListHot(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range talkList {
		reply.Talk = append(reply.Talk, &v1.GetTalkListHotReply_Talk{
			Id:   item.Id,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *BffService) GetUserTalkList(ctx context.Context, req *v1.GetUserTalkListReq) (*v1.GetTalkListReply, error) {
	reply := &v1.GetTalkListReply{Talk: make([]*v1.GetTalkListReply_Talk, 0)}
	talkList, err := s.tc.GetUserTalkList(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range talkList {
		reply.Talk = append(reply.Talk, &v1.GetTalkListReply_Talk{
			Id:   item.Id,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *BffService) GetUserTalkListVisitor(ctx context.Context, req *v1.GetUserTalkListVisitorReq) (*v1.GetTalkListReply, error) {
	reply := &v1.GetTalkListReply{Talk: make([]*v1.GetTalkListReply_Talk, 0)}
	talkList, err := s.tc.GetUserTalkListVisitor(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range talkList {
		reply.Talk = append(reply.Talk, &v1.GetTalkListReply_Talk{
			Id:   item.Id,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *BffService) GetTalkCount(ctx context.Context, _ *emptypb.Empty) (*v1.GetTalkCountReply, error) {
	count, err := s.tc.GetTalkCount(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetTalkCountReply{
		Count: count,
	}, nil
}

func (s *BffService) GetTalkCountVisitor(ctx context.Context, req *v1.GetTalkCountVisitorReq) (*v1.GetTalkCountReply, error) {
	count, err := s.tc.GetTalkCountVisitor(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetTalkCountReply{
		Count: count,
	}, nil
}

func (s *BffService) GetTalkListStatistic(ctx context.Context, req *v1.GetTalkListStatisticReq) (*v1.GetTalkListStatisticReply, error) {
	reply := &v1.GetTalkListStatisticReply{Count: make([]*v1.GetTalkListStatisticReply_Count, 0)}
	statisticList, err := s.tc.GetTalkListStatistic(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	for _, item := range statisticList {
		reply.Count = append(reply.Count, &v1.GetTalkListStatisticReply_Count{
			Id:      item.Id,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
			Comment: item.Comment,
		})
	}
	return reply, nil
}

func (s *BffService) GetTalkStatistic(ctx context.Context, req *v1.GetTalkStatisticReq) (*v1.GetTalkStatisticReply, error) {
	talkStatistic, err := s.tc.GetTalkStatistic(ctx, req.Id)
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

func (s *BffService) GetLastTalkDraft(ctx context.Context, _ *emptypb.Empty) (*v1.GetLastTalkDraftReply, error) {
	draft, err := s.tc.GetLastTalkDraft(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetLastTalkDraftReply{
		Id:     draft.Id,
		Status: draft.Status,
	}, nil
}

func (s *BffService) CreateTalkDraft(ctx context.Context, _ *emptypb.Empty) (*v1.CreateTalkDraftReply, error) {
	id, err := s.tc.CreateTalkDraft(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.CreateTalkDraftReply{
		Id: id,
	}, nil
}

func (s *BffService) SendTalk(ctx context.Context, req *v1.SendTalkReq) (*emptypb.Empty, error) {
	err := s.tc.SendTalk(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SendTalkEdit(ctx context.Context, req *v1.SendTalkEditReq) (*emptypb.Empty, error) {
	err := s.tc.SendTalkEdit(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) DeleteTalk(ctx context.Context, req *v1.DeleteTalkReq) (*emptypb.Empty, error) {
	err := s.tc.DeleteTalk(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SetTalkAgree(ctx context.Context, req *v1.SetTalkAgreeReq) (*emptypb.Empty, error) {
	err := s.tc.SetTalkAgree(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) CancelTalkAgree(ctx context.Context, req *v1.CancelTalkAgreeReq) (*emptypb.Empty, error) {
	err := s.tc.CancelTalkAgree(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) CancelTalkCollect(ctx context.Context, req *v1.CancelTalkCollectReq) (*emptypb.Empty, error) {
	err := s.tc.CancelTalkCollect(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SetTalkView(ctx context.Context, req *v1.SetTalkViewReq) (*emptypb.Empty, error) {
	err := s.tc.SetTalkView(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SetTalkCollect(ctx context.Context, req *v1.SetTalkCollectReq) (*emptypb.Empty, error) {
	err := s.tc.SetTalkCollect(ctx, req.Id, req.CollectionsId, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) TalkStatisticJudge(ctx context.Context, req *v1.TalkStatisticJudgeReq) (*v1.TalkStatisticJudgeReply, error) {
	judge, err := s.tc.TalkStatisticJudge(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.TalkStatisticJudgeReply{
		Agree:   judge.Agree,
		Collect: judge.Collect,
	}, nil
}
