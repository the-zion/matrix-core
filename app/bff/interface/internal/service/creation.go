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

func (s *BffService) GetCollectArticleList(ctx context.Context, req *v1.GetCollectArticleListReq) (*v1.GetArticleListReply, error) {
	reply := &v1.GetArticleListReply{Article: make([]*v1.GetArticleListReply_Article, 0)}
	articleList, err := s.cc.GetCollectArticleList(ctx, req.Id, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range articleList {
		reply.Article = append(reply.Article, &v1.GetArticleListReply_Article{
			Id:      item.Id,
			Uuid:    item.Uuid,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
			Comment: item.Comment,
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

func (s *BffService) GetCollectTalkList(ctx context.Context, req *v1.GetCollectTalkListReq) (*v1.GetTalkListReply, error) {
	reply := &v1.GetTalkListReply{Talk: make([]*v1.GetTalkListReply_Talk, 0)}
	talkList, err := s.cc.GetCollectTalkList(ctx, req.Id, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range talkList {
		reply.Talk = append(reply.Talk, &v1.GetTalkListReply_Talk{
			Id:      item.Id,
			Uuid:    item.Uuid,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
			Comment: item.Comment,
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

func (s *BffService) GetCollectColumnList(ctx context.Context, req *v1.GetCollectColumnListReq) (*v1.GetColumnListReply, error) {
	reply := &v1.GetColumnListReply{Column: make([]*v1.GetColumnListReply_Column, 0)}
	columnList, err := s.cc.GetCollectColumnList(ctx, req.Id, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range columnList {
		reply.Column = append(reply.Column, &v1.GetColumnListReply_Column{
			Id:      item.Id,
			Uuid:    item.Uuid,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
		})
	}
	return reply, nil
}

func (s *BffService) GetCollectColumnCount(ctx context.Context, req *v1.GetCollectColumnCountReq) (*v1.GetCollectColumnCountReply, error) {
	count, err := s.cc.GetCollectColumnCount(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetCollectColumnCountReply{
		Count: count,
	}, nil
}

func (s *BffService) GetLastCollectionsDraft(ctx context.Context, _ *emptypb.Empty) (*v1.GetLastCollectionsDraftReply, error) {
	draft, err := s.cc.GetLastCollectionsDraft(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetLastCollectionsDraftReply{
		Id:     draft.Id,
		Status: draft.Status,
	}, nil
}

func (s *BffService) GetCollections(ctx context.Context, req *v1.GetCollectionsReq) (*v1.GetCollectionsReply, error) {
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

func (s *BffService) GetCollectionsList(ctx context.Context, req *v1.GetCollectionsListReq) (*v1.GetCollectionsListReply, error) {
	reply := &v1.GetCollectionsListReply{Collections: make([]*v1.GetCollectionsListReply_Collections, 0)}
	collections, err := s.cc.GetCollectionsList(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range collections {
		reply.Collections = append(reply.Collections, &v1.GetCollectionsListReply_Collections{
			Id: item.Id,
		})
	}
	return reply, nil
}

func (s *BffService) GetCollectionsListAll(ctx context.Context, _ *emptypb.Empty) (*v1.GetCollectionsListReply, error) {
	reply := &v1.GetCollectionsListReply{Collections: make([]*v1.GetCollectionsListReply_Collections, 0)}
	collections, err := s.cc.GetCollectionsListAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range collections {
		reply.Collections = append(reply.Collections, &v1.GetCollectionsListReply_Collections{
			Id: item.Id,
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

func (s *BffService) GetCollectionsListByVisitor(ctx context.Context, req *v1.GetCollectionsListByVisitorReq) (*v1.GetCollectionsListReply, error) {
	reply := &v1.GetCollectionsListReply{Collections: make([]*v1.GetCollectionsListReply_Collections, 0)}
	collections, err := s.cc.GetCollectionsListByVisitor(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range collections {
		reply.Collections = append(reply.Collections, &v1.GetCollectionsListReply_Collections{
			Id: item.Id,
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

func (s *BffService) CreateCollectionsDraft(ctx context.Context, _ *emptypb.Empty) (*v1.CreateCollectionsDraftReply, error) {
	id, err := s.cc.CreateCollectionsDraft(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.CreateCollectionsDraftReply{
		Id: id,
	}, nil
}

func (s *BffService) SendCollections(ctx context.Context, req *v1.SendCollectionsReq) (*emptypb.Empty, error) {
	err := s.cc.SendCollections(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SendCollectionsEdit(ctx context.Context, req *v1.SendCollectionsEditReq) (*emptypb.Empty, error) {
	err := s.cc.SendCollectionsEdit(ctx, req.Id)
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
			Id:      item.Id,
			Uuid:    item.Uuid,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
			Comment: item.Comment,
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
			Id:      item.Id,
			Uuid:    item.Uuid,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
			Comment: item.Comment,
		})
	}
	return reply, nil
}

func (s *BffService) GetColumnArticleList(ctx context.Context, req *v1.GetColumnArticleListReq) (*v1.GetArticleListReply, error) {
	reply := &v1.GetArticleListReply{Article: make([]*v1.GetArticleListReply_Article, 0)}
	articleList, err := s.ac.GetColumnArticleList(ctx, req.Id)
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
			Id:      item.Id,
			Uuid:    item.Uuid,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
			Comment: item.Comment,
		})
	}
	return reply, nil
}

func (s *BffService) GetUserArticleListSimple(ctx context.Context, req *v1.GetUserArticleListSimpleReq) (*v1.GetArticleListReply, error) {
	reply := &v1.GetArticleListReply{Article: make([]*v1.GetArticleListReply_Article, 0)}
	articleList, err := s.ac.GetUserArticleListSimple(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range articleList {
		reply.Article = append(reply.Article, &v1.GetArticleListReply_Article{
			Id:      item.Id,
			Uuid:    item.Uuid,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
			Comment: item.Comment,
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
			Id:      item.Id,
			Uuid:    item.Uuid,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
			Comment: item.Comment,
		})
	}
	return reply, nil
}

func (s *BffService) GetUserArticleListAll(ctx context.Context, _ *emptypb.Empty) (*v1.GetArticleListReply, error) {
	reply := &v1.GetArticleListReply{Article: make([]*v1.GetArticleListReply_Article, 0)}
	articleList, err := s.ac.GetUserArticleListAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range articleList {
		reply.Article = append(reply.Article, &v1.GetArticleListReply_Article{
			Id:      item.Id,
			Uuid:    item.Uuid,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
			Comment: item.Comment,
		})
	}
	return reply, nil
}

func (s *BffService) GetArticleStatistic(ctx context.Context, req *v1.GetArticleStatisticReq) (*v1.GetArticleStatisticReply, error) {
	articleStatistic, err := s.ac.GetArticleStatistic(ctx, req.Id, req.Uuid)
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

func (s *BffService) GetUserArticleAgree(ctx context.Context, _ *emptypb.Empty) (*v1.GetUserArticleAgreeReply, error) {
	agreeMap, err := s.ac.GetUserArticleAgree(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserArticleAgreeReply{
		Agree: agreeMap,
	}, nil
}

func (s *BffService) GetUserArticleCollect(ctx context.Context, _ *emptypb.Empty) (*v1.GetUserArticleCollectReply, error) {
	collectMap, err := s.ac.GetUserArticleCollect(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserArticleCollectReply{
		Collect: collectMap,
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

func (s *BffService) GetArticleSearch(ctx context.Context, req *v1.GetArticleSearchReq) (*v1.GetArticleSearchReply, error) {
	reply := &v1.GetArticleSearchReply{List: make([]*v1.GetArticleSearchReply_List, 0)}
	articleList, total, err := s.ac.GetArticleSearch(ctx, req.Page, req.Search, req.Time)
	if err != nil {
		return reply, err
	}
	for _, item := range articleList {
		reply.List = append(reply.List, &v1.GetArticleSearchReply_List{
			Id:      item.Id,
			Tags:    item.Tags,
			Title:   item.Title,
			Uuid:    item.Uuid,
			Text:    item.Text,
			Cover:   item.Cover,
			Update:  item.Update,
			Agree:   item.Agree,
			Collect: item.Collect,
			Comment: item.Comment,
			View:    item.View,
		})
	}
	reply.Total = total
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
			Id:      item.Id,
			Uuid:    item.Uuid,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
			Comment: item.Comment,
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
			Id:      item.Id,
			Uuid:    item.Uuid,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
			Comment: item.Comment,
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
			Id:      item.Id,
			Uuid:    item.Uuid,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
			Comment: item.Comment,
		})
	}
	return reply, nil
}

func (s *BffService) GetUserTalkListSimple(ctx context.Context, req *v1.GetUserTalkListSimpleReq) (*v1.GetTalkListReply, error) {
	reply := &v1.GetTalkListReply{Talk: make([]*v1.GetTalkListReply_Talk, 0)}
	talkList, err := s.tc.GetUserTalkListSimple(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range talkList {
		reply.Talk = append(reply.Talk, &v1.GetTalkListReply_Talk{
			Id:      item.Id,
			Uuid:    item.Uuid,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
			Comment: item.Comment,
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
			Id:      item.Id,
			Uuid:    item.Uuid,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
			Comment: item.Comment,
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

func (s *BffService) GetTalkSearch(ctx context.Context, req *v1.GetTalkSearchReq) (*v1.GetTalkSearchReply, error) {
	reply := &v1.GetTalkSearchReply{List: make([]*v1.GetTalkSearchReply_List, 0)}
	talkList, total, err := s.tc.GetTalkSearch(ctx, req.Page, req.Search, req.Time)
	if err != nil {
		return reply, err
	}
	for _, item := range talkList {
		reply.List = append(reply.List, &v1.GetTalkSearchReply_List{
			Id:      item.Id,
			Tags:    item.Tags,
			Title:   item.Title,
			Uuid:    item.Uuid,
			Text:    item.Text,
			Cover:   item.Cover,
			Update:  item.Update,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
			Comment: item.Comment,
		})
	}
	reply.Total = total
	return reply, nil
}

func (s *BffService) GetUserTalkAgree(ctx context.Context, _ *emptypb.Empty) (*v1.GetUserTalkAgreeReply, error) {
	agreeMap, err := s.tc.GetUserTalkAgree(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserTalkAgreeReply{
		Agree: agreeMap,
	}, nil
}

func (s *BffService) GetUserTalkCollect(ctx context.Context, _ *emptypb.Empty) (*v1.GetUserTalkCollectReply, error) {
	collectMap, err := s.tc.GetUserTalkCollect(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserTalkCollectReply{
		Collect: collectMap,
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

// ------------------------------------------column-------------------------------------------------

func (s *BffService) GetLastColumnDraft(ctx context.Context, _ *emptypb.Empty) (*v1.GetLastColumnDraftReply, error) {
	draft, err := s.coc.GetLastColumnDraft(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetLastColumnDraftReply{
		Id:     draft.Id,
		Status: draft.Status,
	}, nil
}

func (s *BffService) CreateColumnDraft(ctx context.Context, _ *emptypb.Empty) (*v1.CreateColumnDraftReply, error) {
	id, err := s.coc.CreateColumnDraft(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.CreateColumnDraftReply{
		Id: id,
	}, nil
}

func (s *BffService) SubscribeColumn(ctx context.Context, req *v1.SubscribeColumnReq) (*emptypb.Empty, error) {
	err := s.coc.SubscribeColumn(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) CancelSubscribeColumn(ctx context.Context, req *v1.CancelSubscribeColumnReq) (*emptypb.Empty, error) {
	err := s.coc.CancelSubscribeColumn(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SubscribeJudge(ctx context.Context, req *v1.SubscribeJudgeReq) (*v1.SubscribeJudgeReply, error) {
	subscribe, err := s.coc.SubscribeJudge(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.SubscribeJudgeReply{
		Subscribe: subscribe,
	}, nil
}

func (s *BffService) GetColumnList(ctx context.Context, req *v1.GetColumnListReq) (*v1.GetColumnListReply, error) {
	reply := &v1.GetColumnListReply{Column: make([]*v1.GetColumnListReply_Column, 0)}
	columnList, err := s.coc.GetColumnList(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range columnList {
		reply.Column = append(reply.Column, &v1.GetColumnListReply_Column{
			Id:      item.Id,
			Uuid:    item.Uuid,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
		})
	}
	return reply, nil
}

func (s *BffService) GetColumnListHot(ctx context.Context, req *v1.GetColumnListHotReq) (*v1.GetColumnListHotReply, error) {
	reply := &v1.GetColumnListHotReply{Column: make([]*v1.GetColumnListHotReply_Column, 0)}
	columnList, err := s.coc.GetColumnListHot(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range columnList {
		reply.Column = append(reply.Column, &v1.GetColumnListHotReply_Column{
			Id:      item.Id,
			Uuid:    item.Uuid,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
		})
	}
	return reply, nil
}

func (s *BffService) GetUserColumnList(ctx context.Context, req *v1.GetUserColumnListReq) (*v1.GetColumnListReply, error) {
	reply := &v1.GetColumnListReply{Column: make([]*v1.GetColumnListReply_Column, 0)}
	columnList, err := s.coc.GetUserColumnList(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range columnList {
		reply.Column = append(reply.Column, &v1.GetColumnListReply_Column{
			Id:      item.Id,
			Uuid:    item.Uuid,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
		})
	}
	return reply, nil
}

func (s *BffService) GetUserColumnListSimple(ctx context.Context, req *v1.GetUserColumnListSimpleReq) (*v1.GetColumnListReply, error) {
	reply := &v1.GetColumnListReply{Column: make([]*v1.GetColumnListReply_Column, 0)}
	columnList, err := s.coc.GetUserColumnListSimple(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range columnList {
		reply.Column = append(reply.Column, &v1.GetColumnListReply_Column{
			Id:      item.Id,
			Uuid:    item.Uuid,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
		})
	}
	return reply, nil
}

func (s *BffService) GetUserColumnListVisitor(ctx context.Context, req *v1.GetUserColumnListVisitorReq) (*v1.GetColumnListReply, error) {
	reply := &v1.GetColumnListReply{Column: make([]*v1.GetColumnListReply_Column, 0)}
	columnList, err := s.coc.GetUserColumnListVisitor(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range columnList {
		reply.Column = append(reply.Column, &v1.GetColumnListReply_Column{
			Id:      item.Id,
			Uuid:    item.Uuid,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
		})
	}
	return reply, nil
}

func (s *BffService) GetColumnCount(ctx context.Context, _ *emptypb.Empty) (*v1.GetColumnCountReply, error) {
	count, err := s.coc.GetColumnCount(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetColumnCountReply{
		Count: count,
	}, nil
}

func (s *BffService) GetColumnCountVisitor(ctx context.Context, req *v1.GetColumnCountVisitorReq) (*v1.GetColumnCountReply, error) {
	count, err := s.coc.GetColumnCountVisitor(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetColumnCountReply{
		Count: count,
	}, nil
}

func (s *BffService) GetColumnListStatistic(ctx context.Context, req *v1.GetColumnListStatisticReq) (*v1.GetColumnListStatisticReply, error) {
	reply := &v1.GetColumnListStatisticReply{Count: make([]*v1.GetColumnListStatisticReply_Count, 0)}
	statisticList, err := s.coc.GetColumnListStatistic(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	for _, item := range statisticList {
		reply.Count = append(reply.Count, &v1.GetColumnListStatisticReply_Count{
			Id:      item.Id,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
		})
	}
	return reply, nil
}

func (s *BffService) GetColumnStatistic(ctx context.Context, req *v1.GetColumnStatisticReq) (*v1.GetColumnStatisticReply, error) {
	columnStatistic, err := s.coc.GetColumnStatistic(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetColumnStatisticReply{
		Uuid:    columnStatistic.Uuid,
		Agree:   columnStatistic.Agree,
		Collect: columnStatistic.Collect,
		View:    columnStatistic.View,
	}, nil
}

func (s *BffService) GetSubscribeList(ctx context.Context, req *v1.GetSubscribeListReq) (*v1.GetSubscribeListReply, error) {
	reply := &v1.GetSubscribeListReply{Subscribe: make([]*v1.GetSubscribeListReply_Subscribe, 0)}
	statisticList, err := s.coc.GetSubscribeList(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range statisticList {
		reply.Subscribe = append(reply.Subscribe, &v1.GetSubscribeListReply_Subscribe{
			Id:      item.Id,
			Uuid:    item.Uuid,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
		})
	}
	return reply, nil
}

func (s *BffService) GetSubscribeListCount(ctx context.Context, req *v1.GetSubscribeListCountReq) (*v1.GetSubscribeListCountReply, error) {
	count, err := s.coc.GetSubscribeListCount(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetSubscribeListCountReply{
		Count: count,
	}, nil
}

func (s *BffService) GetColumnSubscribes(ctx context.Context, req *v1.GetColumnSubscribesReq) (*v1.GetColumnSubscribesReply, error) {
	reply := &v1.GetColumnSubscribesReply{Subscribes: make([]*v1.GetColumnSubscribesReply_Subscribes, 0)}
	subscribesList, err := s.coc.GetColumnSubscribes(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	for _, item := range subscribesList {
		reply.Subscribes = append(reply.Subscribes, &v1.GetColumnSubscribesReply_Subscribes{
			Id:             item.ColumnId,
			SubscribeJudge: item.Status,
		})
	}
	return reply, nil
}

func (s *BffService) GetColumnSearch(ctx context.Context, req *v1.GetColumnSearchReq) (*v1.GetColumnSearchReply, error) {
	reply := &v1.GetColumnSearchReply{List: make([]*v1.GetColumnSearchReply_List, 0)}
	columnList, total, err := s.coc.GetColumnSearch(ctx, req.Page, req.Search, req.Time)
	if err != nil {
		return reply, err
	}
	for _, item := range columnList {
		reply.List = append(reply.List, &v1.GetColumnSearchReply_List{
			Id:        item.Id,
			Tags:      item.Tags,
			Name:      item.Name,
			Uuid:      item.Uuid,
			Introduce: item.Introduce,
			Cover:     item.Cover,
			Update:    item.Update,
			Agree:     item.Agree,
			Collect:   item.Collect,
			View:      item.View,
		})
	}
	reply.Total = total
	return reply, nil
}

func (s *BffService) GetUserColumnAgree(ctx context.Context, _ *emptypb.Empty) (*v1.GetUserColumnAgreeReply, error) {
	agreeMap, err := s.coc.GetUserColumnAgree(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserColumnAgreeReply{
		Agree: agreeMap,
	}, nil
}

func (s *BffService) GetUserColumnCollect(ctx context.Context, _ *emptypb.Empty) (*v1.GetUserColumnCollectReply, error) {
	agreeMap, err := s.coc.GetUserColumnCollect(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserColumnCollectReply{
		Collect: agreeMap,
	}, nil
}

func (s *BffService) GetUserSubscribeColumn(ctx context.Context, _ *emptypb.Empty) (*v1.GetUserSubscribeColumnReply, error) {
	collectMap, err := s.coc.GetUserSubscribeColumn(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserSubscribeColumnReply{
		Subscribe: collectMap,
	}, nil
}

func (s *BffService) ColumnStatisticJudge(ctx context.Context, req *v1.ColumnStatisticJudgeReq) (*v1.ColumnStatisticJudgeReply, error) {
	judge, err := s.coc.ColumnStatisticJudge(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.ColumnStatisticJudgeReply{
		Agree:   judge.Agree,
		Collect: judge.Collect,
	}, nil
}

func (s *BffService) SendColumn(ctx context.Context, req *v1.SendColumnReq) (*emptypb.Empty, error) {
	err := s.coc.SendColumn(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SendColumnEdit(ctx context.Context, req *v1.SendColumnEditReq) (*emptypb.Empty, error) {
	err := s.coc.SendColumnEdit(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) DeleteColumn(ctx context.Context, req *v1.DeleteColumnReq) (*emptypb.Empty, error) {
	err := s.coc.DeleteColumn(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SetColumnAgree(ctx context.Context, req *v1.SetColumnAgreeReq) (*emptypb.Empty, error) {
	err := s.coc.SetColumnAgree(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) CancelColumnAgree(ctx context.Context, req *v1.CancelColumnAgreeReq) (*emptypb.Empty, error) {
	err := s.coc.CancelColumnAgree(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SetColumnCollect(ctx context.Context, req *v1.SetColumnCollectReq) (*emptypb.Empty, error) {
	err := s.coc.SetColumnCollect(ctx, req.Id, req.CollectionsId, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) CancelColumnCollect(ctx context.Context, req *v1.CancelColumnCollectReq) (*emptypb.Empty, error) {
	err := s.coc.CancelColumnCollect(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SetColumnView(ctx context.Context, req *v1.SetColumnViewReq) (*emptypb.Empty, error) {
	err := s.coc.SetColumnView(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) AddColumnIncludes(ctx context.Context, req *v1.AddColumnIncludesReq) (*emptypb.Empty, error) {
	err := s.coc.AddColumnIncludes(ctx, req.Id, req.ArticleId)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) DeleteColumnIncludes(ctx context.Context, req *v1.DeleteColumnIncludesReq) (*emptypb.Empty, error) {
	err := s.coc.DeleteColumnIncludes(ctx, req.Id, req.ArticleId)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) GetNews(ctx context.Context, req *v1.GetNewsReq) (*v1.GetNewsReply, error) {
	reply := &v1.GetNewsReply{News: make([]*v1.GetNewsReply_News, 0)}
	newsList, err := s.nc.GetNews(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range newsList {
		reply.News = append(reply.News, &v1.GetNewsReply_News{
			Id:     item.Id,
			Update: item.Update,
			Author: item.Author,
			Title:  item.Title,
			Text:   item.Text,
			Tags:   item.Tags,
			Cover:  item.Cover,
			Url:    item.Url,
		})
	}
	return reply, nil
}
