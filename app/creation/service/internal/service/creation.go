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

func (s *CreationService) GetCollectTalk(ctx context.Context, req *v1.GetCollectTalkReq) (*v1.GetTalkListReply, error) {
	reply := &v1.GetTalkListReply{Talk: make([]*v1.GetTalkListReply_Talk, 0)}
	talkList, err := s.cc.GetCollectTalk(ctx, req.Id, req.Page)
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

func (s *CreationService) GetCollectTalkCount(ctx context.Context, req *v1.GetCollectTalkCountReq) (*v1.GetCollectTalkCountReply, error) {
	count, err := s.cc.GetCollectTalkCount(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetCollectTalkCountReply{
		Count: count,
	}, nil
}

func (s *CreationService) GetCollectColumn(ctx context.Context, req *v1.GetCollectColumnReq) (*v1.GetColumnListReply, error) {
	reply := &v1.GetColumnListReply{Column: make([]*v1.GetColumnListReply_Column, 0)}
	columnList, err := s.cc.GetCollectColumn(ctx, req.Id, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range columnList {
		reply.Column = append(reply.Column, &v1.GetColumnListReply_Column{
			Id:   item.ColumnId,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *CreationService) GetCollectColumnCount(ctx context.Context, req *v1.GetCollectColumnCountReq) (*v1.GetCollectColumnCountReply, error) {
	count, err := s.cc.GetCollectColumnCount(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetCollectColumnCountReply{
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

func (s *CreationService) GetCollectionListInfo(ctx context.Context, req *v1.GetCollectionListInfoReq) (*v1.GetCollectionsReply, error) {
	reply := &v1.GetCollectionsReply{Collections: make([]*v1.GetCollectionsReply_Collections, 0)}
	collectionsListInfo, err := s.cc.GetCollectionListInfo(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	for _, item := range collectionsListInfo {
		reply.Collections = append(reply.Collections, &v1.GetCollectionsReply_Collections{
			Id:        item.Id,
			Name:      item.Name,
			Introduce: item.Introduce,
		})
	}
	return reply, nil
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

func (s *CreationService) AddCreationComment(ctx context.Context, req *v1.AddCreationCommentReq) (*emptypb.Empty, error) {
	var err error
	if req.CreationType == 1 {
		err = s.ac.AddArticleComment(ctx, req.CreationId, req.Uuid)
	}
	if req.CreationType == 3 {
		err = s.tc.AddTalkComment(ctx, req.CreationId, req.Uuid)
	}
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) ReduceCreationComment(ctx context.Context, req *v1.ReduceCreationCommentReq) (*emptypb.Empty, error) {
	var err error
	if req.CreationType == 1 {
		err = s.ac.ReduceArticleComment(ctx, req.CreationId, req.Uuid)
	}
	if req.CreationType == 3 {
		err = s.tc.ReduceTalkComment(ctx, req.CreationId, req.Uuid)
	}
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
