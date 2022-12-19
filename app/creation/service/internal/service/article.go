package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
	"github.com/the-zion/matrix-core/app/creation/service/internal/biz"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *CreationService) GetArticleList(ctx context.Context, req *v1.GetArticleListReq) (*v1.GetArticleListReply, error) {
	articleList, err := s.ac.GetArticleList(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetArticleListReply{Article: make([]*v1.GetArticleListReply_Article, 0, len(articleList))}
	for _, item := range articleList {
		reply.Article = append(reply.Article, &v1.GetArticleListReply_Article{
			Id:   item.ArticleId,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *CreationService) GetArticleListHot(ctx context.Context, req *v1.GetArticleListHotReq) (*v1.GetArticleListHotReply, error) {
	articleList, err := s.ac.GetArticleListHot(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetArticleListHotReply{Article: make([]*v1.GetArticleListHotReply_Article, 0, len(articleList))}
	for _, item := range articleList {
		reply.Article = append(reply.Article, &v1.GetArticleListHotReply_Article{
			Id:   item.ArticleId,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *CreationService) GetColumnArticleList(ctx context.Context, req *v1.GetColumnArticleListReq) (*v1.GetArticleListReply, error) {
	articleList, err := s.ac.GetColumnArticleList(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetArticleListReply{Article: make([]*v1.GetArticleListReply_Article, 0, len(articleList))}
	for _, item := range articleList {
		reply.Article = append(reply.Article, &v1.GetArticleListReply_Article{
			Id:   item.ArticleId,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *CreationService) GetArticleCount(ctx context.Context, req *v1.GetArticleCountReq) (*v1.GetArticleCountReply, error) {
	count, err := s.ac.GetArticleCount(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetArticleCountReply{
		Count: count,
	}, nil
}

func (s *CreationService) GetArticleCountVisitor(ctx context.Context, req *v1.GetArticleCountVisitorReq) (*v1.GetArticleCountReply, error) {
	count, err := s.ac.GetArticleCountVisitor(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetArticleCountReply{
		Count: count,
	}, nil
}

func (s *CreationService) GetUserArticleList(ctx context.Context, req *v1.GetUserArticleListReq) (*v1.GetArticleListReply, error) {
	articleList, err := s.ac.GetUserArticleList(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetArticleListReply{Article: make([]*v1.GetArticleListReply_Article, 0, len(articleList))}
	for _, item := range articleList {
		reply.Article = append(reply.Article, &v1.GetArticleListReply_Article{
			Id:   item.ArticleId,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *CreationService) GetUserArticleListVisitor(ctx context.Context, req *v1.GetUserArticleListVisitorReq) (*v1.GetArticleListReply, error) {
	articleList, err := s.ac.GetUserArticleListVisitor(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetArticleListReply{Article: make([]*v1.GetArticleListReply_Article, 0, len(articleList))}
	for _, item := range articleList {
		reply.Article = append(reply.Article, &v1.GetArticleListReply_Article{
			Id:   item.ArticleId,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *CreationService) GetUserArticleListAll(ctx context.Context, req *v1.GetUserArticleListAllReq) (*v1.GetArticleListReply, error) {
	articleList, err := s.ac.GetUserArticleListAll(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetArticleListReply{Article: make([]*v1.GetArticleListReply_Article, 0, len(articleList))}
	for _, item := range articleList {
		reply.Article = append(reply.Article, &v1.GetArticleListReply_Article{
			Id:   item.ArticleId,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *CreationService) GetArticleStatistic(ctx context.Context, req *v1.GetArticleStatisticReq) (*v1.GetArticleStatisticReply, error) {
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

func (s *CreationService) GetUserArticleAgree(ctx context.Context, req *v1.GetUserArticleAgreeReq) (*v1.GetUserArticleAgreeReply, error) {
	agreeMap, err := s.ac.GetUserArticleAgree(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserArticleAgreeReply{
		Agree: agreeMap,
	}, nil
}

func (s *CreationService) GetUserArticleCollect(ctx context.Context, req *v1.GetUserArticleCollectReq) (*v1.GetUserArticleCollectReply, error) {
	collectMap, err := s.ac.GetUserArticleCollect(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserArticleCollectReply{
		Collect: collectMap,
	}, nil
}

func (s *CreationService) GetArticleListStatistic(ctx context.Context, req *v1.GetArticleListStatisticReq) (*v1.GetArticleListStatisticReply, error) {
	statisticList, err := s.ac.GetArticleListStatistic(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetArticleListStatisticReply{Count: make([]*v1.GetArticleListStatisticReply_Count, 0, len(statisticList))}
	for _, item := range statisticList {
		reply.Count = append(reply.Count, &v1.GetArticleListStatisticReply_Count{
			Id:      item.ArticleId,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
			Comment: item.Comment,
		})
	}
	return reply, nil
}

func (s *CreationService) GetLastArticleDraft(ctx context.Context, req *v1.GetLastArticleDraftReq) (*v1.GetLastArticleDraftReply, error) {
	draft, err := s.ac.GetLastArticleDraft(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetLastArticleDraftReply{
		Id:     draft.Id,
		Status: draft.Status,
	}, nil
}

func (s *CreationService) GetArticleSearch(ctx context.Context, req *v1.GetArticleSearchReq) (*v1.GetArticleSearchReply, error) {
	articleList, total, err := s.ac.GetArticleSearch(ctx, req.Page, req.Search, req.Time)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetArticleSearchReply{List: make([]*v1.GetArticleSearchReply_List, 0, len(articleList))}
	for _, item := range articleList {
		reply.List = append(reply.List, &v1.GetArticleSearchReply_List{
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

func (s *CreationService) GetArticleImageReview(ctx context.Context, req *v1.GetArticleImageReviewReq) (*v1.GetArticleImageReviewReply, error) {
	reviewList, err := s.ac.GetArticleImageReview(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetArticleImageReviewReply{Review: make([]*v1.GetArticleImageReviewReply_Review, 0, len(reviewList))}
	for _, item := range reviewList {
		reply.Review = append(reply.Review, &v1.GetArticleImageReviewReply_Review{
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

func (s *CreationService) GetArticleContentReview(ctx context.Context, req *v1.GetArticleContentReviewReq) (*v1.GetArticleContentReviewReply, error) {
	reviewList, err := s.ac.GetArticleContentReview(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetArticleContentReviewReply{Review: make([]*v1.GetArticleContentReviewReply_Review, 0, len(reviewList))}
	for _, item := range reviewList {
		reply.Review = append(reply.Review, &v1.GetArticleContentReviewReply_Review{
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

func (s *CreationService) ArticleImageIrregular(ctx context.Context, req *v1.CreationImageIrregularReq) (*emptypb.Empty, error) {
	err := s.ac.ArticleImageIrregular(ctx, &biz.ImageReview{
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
		Mode:       "add_article_image_review_db_and_cache",
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) ArticleContentIrregular(ctx context.Context, req *v1.CreationContentIrregularReq) (*emptypb.Empty, error) {
	err := s.ac.ArticleContentIrregular(ctx, &biz.TextReview{
		CreationId: req.Id,
		Uuid:       req.Uuid,
		JobId:      req.JobId,
		Title:      req.Title,
		Kind:       req.Kind,
		Label:      req.Label,
		Result:     req.Result,
		Section:    req.Section,
		Mode:       "add_article_content_review_db_and_cache",
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) AddArticleImageReviewDbAndCache(ctx context.Context, req *v1.AddCreationImageReviewDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.ac.AddArticleImageReviewDbAndCache(ctx, &biz.ImageReview{
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

func (s *CreationService) AddArticleContentReviewDbAndCache(ctx context.Context, req *v1.AddCreationContentReviewDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.ac.AddArticleContentReviewDbAndCache(ctx, &biz.TextReview{
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

func (s *CreationService) CreateArticle(ctx context.Context, req *v1.CreateArticleReq) (*emptypb.Empty, error) {
	err := s.ac.CreateArticle(ctx, req.Id, req.Auth, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) EditArticle(ctx context.Context, req *v1.EditArticleReq) (*emptypb.Empty, error) {
	err := s.ac.EditArticle(ctx, req.Id, req.Auth, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) DeleteArticle(ctx context.Context, req *v1.DeleteArticleReq) (*emptypb.Empty, error) {
	err := s.ac.DeleteArticle(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) DeleteArticleDraft(ctx context.Context, req *v1.DeleteArticleDraftReq) (*emptypb.Empty, error) {
	err := s.ac.DeleteArticleDraft(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CreateArticleDbCacheAndSearch(ctx context.Context, req *v1.CreateArticleDbCacheAndSearchReq) (*emptypb.Empty, error) {
	err := s.ac.CreateArticleDbCacheAndSearch(ctx, req.Id, req.Auth, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) EditArticleCosAndSearch(ctx context.Context, req *v1.EditArticleCosAndSearchReq) (*emptypb.Empty, error) {
	err := s.ac.EditArticleCosAndSearch(ctx, req.Id, req.Auth, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) DeleteArticleCacheAndSearch(ctx context.Context, req *v1.DeleteArticleCacheAndSearchReq) (*emptypb.Empty, error) {
	err := s.ac.DeleteArticleCacheAndSearch(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CreateArticleDraft(ctx context.Context, req *v1.CreateArticleDraftReq) (*v1.CreateArticleDraftReply, error) {
	id, err := s.ac.CreateArticleDraft(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.CreateArticleDraftReply{
		Id: id,
	}, nil
}

func (s *CreationService) ArticleDraftMark(ctx context.Context, req *v1.ArticleDraftMarkReq) (*emptypb.Empty, error) {
	err := s.ac.ArticleDraftMark(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) GetArticleDraftList(ctx context.Context, req *v1.GetArticleDraftListReq) (*v1.GetArticleDraftListReply, error) {
	draftList, err := s.ac.GetArticleDraftList(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetArticleDraftListReply{Draft: make([]*v1.GetArticleDraftListReply_Draft, 0, len(draftList))}
	for _, item := range draftList {
		reply.Draft = append(reply.Draft, &v1.GetArticleDraftListReply_Draft{
			Id: item.Id,
		})
	}
	return reply, nil
}

func (s *CreationService) SendArticle(ctx context.Context, req *v1.SendArticleReq) (*emptypb.Empty, error) {
	err := s.ac.SendArticle(ctx, req.Id, req.Uuid, req.Ip)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SendArticleEdit(ctx context.Context, req *v1.SendArticleEditReq) (*emptypb.Empty, error) {
	err := s.ac.SendArticleEdit(ctx, req.Id, req.Uuid, req.Ip)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SetArticleAgree(ctx context.Context, req *v1.SetArticleAgreeReq) (*emptypb.Empty, error) {
	err := s.ac.SetArticleAgree(ctx, req.Id, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SetArticleView(ctx context.Context, req *v1.SetArticleViewReq) (*emptypb.Empty, error) {
	err := s.ac.SetArticleView(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SetArticleViewDbAndCache(ctx context.Context, req *v1.SetArticleViewDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.ac.SetArticleViewDbAndCache(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SetArticleAgreeDbAndCache(ctx context.Context, req *v1.SetArticleAgreeDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.ac.SetArticleAgreeDbAndCache(ctx, req.Id, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SetArticleCollect(ctx context.Context, req *v1.SetArticleCollectReq) (*emptypb.Empty, error) {
	err := s.ac.SetArticleCollect(ctx, req.Id, req.CollectionsId, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SetArticleCollectDbAndCache(ctx context.Context, req *v1.SetArticleCollectDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.ac.SetArticleCollectDbAndCache(ctx, req.Id, req.CollectionsId, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CancelArticleAgree(ctx context.Context, req *v1.CancelArticleAgreeReq) (*emptypb.Empty, error) {
	err := s.ac.CancelArticleAgree(ctx, req.Id, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CancelArticleAgreeDbAndCache(ctx context.Context, req *v1.CancelArticleAgreeDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.ac.CancelArticleAgreeDbAndCache(ctx, req.Id, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CancelArticleCollect(ctx context.Context, req *v1.CancelArticleCollectReq) (*emptypb.Empty, error) {
	err := s.ac.CancelArticleCollect(ctx, req.Id, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CancelArticleCollectDbAndCache(ctx context.Context, req *v1.CancelArticleCollectDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.ac.CancelArticleCollectDbAndCache(ctx, req.Id, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) ArticleStatisticJudge(ctx context.Context, req *v1.ArticleStatisticJudgeReq) (*v1.ArticleStatisticJudgeReply, error) {
	judge, err := s.ac.ArticleStatisticJudge(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.ArticleStatisticJudgeReply{
		Agree:   judge.Agree,
		Collect: judge.Collect,
	}, nil
}
