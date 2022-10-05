package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
	"github.com/the-zion/matrix-core/app/creation/service/internal/biz"
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

func (s *CreationService) GetLastCollectionsDraft(ctx context.Context, req *v1.GetLastCollectionsDraftReq) (*v1.GetLastCollectionsDraftReply, error) {
	draft, err := s.cc.GetLastCollectionsDraft(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetLastCollectionsDraftReply{
		Id:     draft.Id,
		Status: draft.Status,
	}, nil
}

func (s *CreationService) GetCollectionsContentReview(ctx context.Context, req *v1.GetCollectionsContentReviewReq) (*v1.GetCollectionsContentReviewReply, error) {
	reply := &v1.GetCollectionsContentReviewReply{Review: make([]*v1.GetCollectionsContentReviewReply_Review, 0)}
	reviewList, err := s.cc.GetCollectionsContentReview(ctx, req.Page, req.Uuid)
	if err != nil {
		return reply, err
	}
	for _, item := range reviewList {
		reply.Review = append(reply.Review, &v1.GetCollectionsContentReviewReply_Review{
			Id:         item.Id,
			CreationId: item.CreationId,
			Title:      item.Title,
			Kind:       item.Kind,
			Uuid:       item.Uuid,
			CreateAt:   item.CreateAt,
			JobId:      item.JobId,
			Label:      item.Label,
			Result:     item.Result,
			Section:    item.Section,
		})
	}
	return reply, nil
}

func (s *CreationService) GetCollectArticleList(ctx context.Context, req *v1.GetCollectArticleListReq) (*v1.GetArticleListReply, error) {
	reply := &v1.GetArticleListReply{Article: make([]*v1.GetArticleListReply_Article, 0)}
	articleList, err := s.cc.GetCollectArticleList(ctx, req.Id, req.Page)
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

func (s *CreationService) GetCollectTalkList(ctx context.Context, req *v1.GetCollectTalkListReq) (*v1.GetTalkListReply, error) {
	reply := &v1.GetTalkListReply{Talk: make([]*v1.GetTalkListReply_Talk, 0)}
	talkList, err := s.cc.GetCollectTalkList(ctx, req.Id, req.Page)
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

func (s *CreationService) GetCollectColumnList(ctx context.Context, req *v1.GetCollectColumnListReq) (*v1.GetColumnListReply, error) {
	reply := &v1.GetColumnListReply{Column: make([]*v1.GetColumnListReply_Column, 0)}
	columnList, err := s.cc.GetCollectColumnList(ctx, req.Id, req.Page)
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

func (s *CreationService) GetCollections(ctx context.Context, req *v1.GetCollectionsReq) (*v1.GetCollectionsReply, error) {
	collection, err := s.cc.GetCollections(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetCollectionsReply{
		Uuid:    collection.Uuid,
		Auth:    collection.Auth,
		Article: collection.Article,
		Column:  collection.Column,
		Talk:    collection.Talk,
	}, nil
}

func (s *CreationService) GetCollectionListInfo(ctx context.Context, req *v1.GetCollectionListInfoReq) (*v1.GetCollectionsListReply, error) {
	reply := &v1.GetCollectionsListReply{Collections: make([]*v1.GetCollectionsListReply_Collections, 0)}
	collectionsListInfo, err := s.cc.GetCollectionListInfo(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	for _, item := range collectionsListInfo {
		reply.Collections = append(reply.Collections, &v1.GetCollectionsListReply_Collections{
			Id: item.CollectionsId,
		})
	}
	return reply, nil
}

func (s *CreationService) GetCollectionsList(ctx context.Context, req *v1.GetCollectionsListReq) (*v1.GetCollectionsListReply, error) {
	reply := &v1.GetCollectionsListReply{Collections: make([]*v1.GetCollectionsListReply_Collections, 0)}
	collections, err := s.cc.GetCollectionsList(ctx, req.Uuid, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range collections {
		reply.Collections = append(reply.Collections, &v1.GetCollectionsListReply_Collections{
			Id: item.CollectionsId,
		})
	}
	return reply, nil
}

func (s *CreationService) GetCollectionsListAll(ctx context.Context, req *v1.GetCollectionsListAllReq) (*v1.GetCollectionsListReply, error) {
	reply := &v1.GetCollectionsListReply{Collections: make([]*v1.GetCollectionsListReply_Collections, 0)}
	collections, err := s.cc.GetCollectionsListAll(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range collections {
		reply.Collections = append(reply.Collections, &v1.GetCollectionsListReply_Collections{
			Id: item.CollectionsId,
		})
	}
	return reply, nil
}

func (s *CreationService) GetCollectionsListByVisitor(ctx context.Context, req *v1.GetCollectionsListReq) (*v1.GetCollectionsListReply, error) {
	reply := &v1.GetCollectionsListReply{Collections: make([]*v1.GetCollectionsListReply_Collections, 0)}
	collections, err := s.cc.GetCollectionsListByVisitor(ctx, req.Uuid, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range collections {
		reply.Collections = append(reply.Collections, &v1.GetCollectionsListReply_Collections{
			Id: item.CollectionsId,
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

func (s *CreationService) GetCreationUser(ctx context.Context, req *v1.GetCreationUserReq) (*v1.GetCreationUserReply, error) {
	creationUser, err := s.cc.GetCreationUser(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetCreationUserReply{
		Article:     creationUser.Article,
		Talk:        creationUser.Talk,
		Collections: creationUser.Collections,
		Column:      creationUser.Column,
		Collect:     creationUser.Collect,
		Subscribe:   creationUser.Subscribe,
	}, nil
}

func (s *CreationService) GetCreationUserVisitor(ctx context.Context, req *v1.GetCreationUserReq) (*v1.GetCreationUserReply, error) {
	creationUser, err := s.cc.GetCreationUserVisitor(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetCreationUserReply{
		Article:     creationUser.Article,
		Talk:        creationUser.Talk,
		Collections: creationUser.Collections,
		Column:      creationUser.Column,
	}, nil
}

func (s *CreationService) GetUserTimeLineList(ctx context.Context, req *v1.GetUserTimeLineListReq) (*v1.GetUserTimeLineListReply, error) {
	reply := &v1.GetUserTimeLineListReply{Timeline: make([]*v1.GetUserTimeLineListReply_TimeLine, 0)}
	timeline, err := s.cc.GetUserTimeLineList(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range timeline {
		reply.Timeline = append(reply.Timeline, &v1.GetUserTimeLineListReply_TimeLine{
			Id:         item.Id,
			Uuid:       item.Uuid,
			CreationId: item.CreationId,
			Mode:       item.Mode,
		})
	}
	return reply, nil
}

func (s *CreationService) CreateCollectionsDraft(ctx context.Context, req *v1.CreateCollectionsDraftReq) (*v1.CreateCollectionsDraftReply, error) {
	id, err := s.cc.CreateCollectionsDraft(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.CreateCollectionsDraftReply{
		Id: id,
	}, nil
}

func (s *CreationService) SendCollections(ctx context.Context, req *v1.SendCollectionsReq) (*emptypb.Empty, error) {
	err := s.cc.SendCollections(ctx, req.Id, req.Uuid, req.Ip)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CreateCollections(ctx context.Context, req *v1.CreateCollectionsReq) (*emptypb.Empty, error) {
	err := s.cc.CreateCollections(ctx, req.Id, req.Auth, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CreateCollectionsDbAndCache(ctx context.Context, req *v1.CreateCollectionsDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.cc.CreateCollectionsDbAndCache(ctx, req.Id, req.Auth, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CollectionsContentIrregular(ctx context.Context, req *v1.CreationContentIrregularReq) (*emptypb.Empty, error) {
	err := s.cc.CollectionsContentIrregular(ctx, &biz.TextReview{
		CreationId: req.Id,
		Uuid:       req.Uuid,
		JobId:      req.JobId,
		Title:      req.Title,
		Kind:       req.Kind,
		Label:      req.Label,
		Result:     req.Result,
		Section:    req.Section,
		Mode:       "add_collections_content_review_db_and_cache",
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) AddCollectionsContentReviewDbAndCache(ctx context.Context, req *v1.AddCreationContentReviewDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.cc.AddCollectionsContentReviewDbAndCache(ctx, &biz.TextReview{
		CreationId: req.CreationId,
		Uuid:       req.Uuid,
		JobId:      req.JobId,
		Title:      req.Title,
		Kind:       req.Kind,
		Label:      req.Label,
		Result:     req.Result,
		Section:    req.Section,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SendCollectionsEdit(ctx context.Context, req *v1.SendCollectionsEditReq) (*emptypb.Empty, error) {
	err := s.cc.SendCollectionsEdit(ctx, req.Id, req.Uuid, req.Ip)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) EditCollections(ctx context.Context, req *v1.EditCollectionsReq) (*emptypb.Empty, error) {
	err := s.cc.EditCollections(ctx, req.Id, req.Auth, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) EditCollectionsCos(ctx context.Context, req *v1.EditCollectionsCosReq) (*emptypb.Empty, error) {
	err := s.cc.EditCollectionsCos(ctx, req.Id, req.Auth, req.Uuid)
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

func (s *CreationService) DeleteCollectionsCache(ctx context.Context, req *v1.DeleteCollectionsCacheReq) (*emptypb.Empty, error) {
	err := s.cc.DeleteCollectionsCache(ctx, req.Id, req.Uuid)
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
