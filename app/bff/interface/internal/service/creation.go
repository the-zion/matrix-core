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
