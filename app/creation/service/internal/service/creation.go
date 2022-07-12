package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *CreationService) GetLeaderBoard(ctx context.Context, _ *emptypb.Empty) (*v1.GetLeaderBoardReply, error) {
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

func (s *CreationService) GetCollectArticle(ctx context.Context, req *v1.GetCollectArticleReq) (*v1.GetArticleListReply, error) {
	reply := &v1.GetArticleListReply{Article: make([]*v1.GetArticleListReply_Article, 0)}
	articleList, err := s.cc.GetCollectArticle(ctx, req.Id, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range articleList {
		reply.Article = append(reply.Article, &v1.GetArticleListReply_Article{
			Id:   item.ArticleId,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *CreationService) GetCollectArticleCount(ctx context.Context, req *v1.GetCollectArticleCountReq) (*v1.GetCollectArticleCountReply, error) {
	count, err := s.cc.GetCollectArticleCount(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetCollectArticleCountReply{
		Count: count,
	}, nil
}

func (s *CreationService) GetCollection(ctx context.Context, req *v1.GetCollectionReq) (*v1.GetCollectionReply, error) {
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

func (s *CreationService) GetCollections(ctx context.Context, req *v1.GetCollectionsReq) (*v1.GetCollectionsReply, error) {
	reply := &v1.GetCollectionsReply{Collections: make([]*v1.GetCollectionsReply_Collections, 0)}
	collections, err := s.cc.GetCollections(ctx, req.Uuid, req.Page)
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

func (s *CreationService) GetCollectionsByVisitor(ctx context.Context, req *v1.GetCollectionsReq) (*v1.GetCollectionsReply, error) {
	reply := &v1.GetCollectionsReply{Collections: make([]*v1.GetCollectionsReply_Collections, 0)}
	collections, err := s.cc.GetCollectionsByVisitor(ctx, req.Uuid, req.Page)
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

func (s *CreationService) GetCollectionsCount(ctx context.Context, req *v1.GetCollectionsCountReq) (*v1.GetCollectionsCountReply, error) {
	count, err := s.cc.GetCollectionsCount(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetCollectionsCountReply{
		Count: count,
	}, nil
}

func (s *CreationService) GetCollectionsVisitorCount(ctx context.Context, req *v1.GetCollectionsCountReq) (*v1.GetCollectionsCountReply, error) {
	count, err := s.cc.GetCollectionsVisitorCount(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetCollectionsCountReply{
		Count: count,
	}, nil
}

func (s *CreationService) CreateCollections(ctx context.Context, req *v1.CreateCollectionsReq) (*emptypb.Empty, error) {
	err := s.cc.CreateCollections(ctx, req.Uuid, req.Name, req.Introduce, req.Auth)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) EditCollections(ctx context.Context, req *v1.EditCollectionsReq) (*emptypb.Empty, error) {
	err := s.cc.EditCollections(ctx, req.Id, req.Uuid, req.Name, req.Introduce, req.Auth)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) DeleteCollections(ctx context.Context, req *v1.DeleteCollectionsReq) (*emptypb.Empty, error) {
	err := s.cc.DeleteCollections(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
