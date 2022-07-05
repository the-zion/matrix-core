package service

import (
	"context"
	"github.com/the-zion/matrix-core/api/bff/interface/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

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
