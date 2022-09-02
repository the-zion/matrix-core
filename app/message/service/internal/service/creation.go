package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/message/service/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *MessageService) ToReviewCreateArticle(id int32, uuid string) error {
	return s.cc.ToReviewCreateArticle(id, uuid)
}

func (s *MessageService) ToReviewEditArticle(id int32, uuid string) error {
	return s.cc.ToReviewEditArticle(id, uuid)
}

func (s *MessageService) ArticleCreateReview(ctx context.Context, req *v1.TextReviewReq) (*emptypb.Empty, error) {
	tr := s.TextReview(req)
	err := s.cc.ArticleCreateReview(ctx, tr)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *MessageService) ArticleEditReview(ctx context.Context, req *v1.TextReviewReq) (*emptypb.Empty, error) {
	tr := s.TextReview(req)
	err := s.cc.ArticleEditReview(ctx, tr)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *MessageService) CreateArticleDbCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return s.cc.CreateArticleDbCacheAndSearch(ctx, id, auth, uuid)
}

func (s *MessageService) EditArticleCosAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return s.cc.EditArticleCosAndSearch(ctx, id, auth, uuid)
}

func (s *MessageService) DeleteArticleCacheAndSearch(ctx context.Context, id int32, uuid string) error {
	return s.cc.DeleteArticleCacheAndSearch(ctx, id, uuid)
}

func (s *MessageService) SetArticleViewDbAndCache(ctx context.Context, id int32, uuid string) error {
	return s.cc.SetArticleViewDbAndCache(ctx, id, uuid)
}

func (s *MessageService) SetArticleAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return s.cc.SetArticleAgreeDbAndCache(ctx, id, uuid, userUuid)
}

func (s *MessageService) CancelArticleAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return s.cc.CancelArticleAgreeDbAndCache(ctx, id, uuid, userUuid)
}

func (s *MessageService) SetArticleCollectDbAndCache(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	return s.cc.SetArticleCollectDbAndCache(ctx, id, collectionsId, uuid, userUuid)
}

func (s *MessageService) CancelArticleCollectDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return s.cc.CancelArticleCollectDbAndCache(ctx, id, uuid, userUuid)
}

func (s *MessageService) ToReviewCreateTalk(id int32, uuid string) error {
	return s.cc.ToReviewCreateTalk(id, uuid)
}

func (s *MessageService) ToReviewEditTalk(id int32, uuid string) error {
	return s.cc.ToReviewEditTalk(id, uuid)
}

func (s *MessageService) TalkCreateReview(ctx context.Context, req *v1.TextReviewReq) (*emptypb.Empty, error) {
	tr := s.TextReview(req)
	err := s.cc.TalkCreateReview(ctx, tr)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *MessageService) TalkEditReview(ctx context.Context, req *v1.TextReviewReq) (*emptypb.Empty, error) {
	tr := s.TextReview(req)
	err := s.cc.TalkEditReview(ctx, tr)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *MessageService) CreateTalkDbCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return s.cc.CreateTalkDbCacheAndSearch(ctx, id, auth, uuid)
}

func (s *MessageService) EditTalkCosAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return s.cc.EditTalkCosAndSearch(ctx, id, auth, uuid)
}

func (s *MessageService) DeleteTalkCacheAndSearch(ctx context.Context, id int32, uuid string) error {
	return s.cc.DeleteTalkCacheAndSearch(ctx, id, uuid)
}

func (s *MessageService) SetTalkViewDbAndCache(ctx context.Context, id int32, uuid string) error {
	return s.cc.SetTalkViewDbAndCache(ctx, id, uuid)
}

func (s *MessageService) SetTalkAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return s.cc.SetTalkAgreeDbAndCache(ctx, id, uuid, userUuid)
}

func (s *MessageService) CancelTalkAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return s.cc.CancelTalkAgreeDbAndCache(ctx, id, uuid, userUuid)
}

func (s *MessageService) SetTalkCollectDbAndCache(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	return s.cc.SetTalkCollectDbAndCache(ctx, id, collectionsId, uuid, userUuid)
}

func (s *MessageService) CancelTalkCollectDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return s.cc.CancelTalkCollectDbAndCache(ctx, id, uuid, userUuid)
}

func (s *MessageService) ToReviewCreateColumn(id int32, uuid string) error {
	return s.cc.ToReviewCreateColumn(id, uuid)
}

func (s *MessageService) ToReviewEditColumn(id int32, uuid string) error {
	return s.cc.ToReviewEditColumn(id, uuid)
}

func (s *MessageService) ToReviewCreateCollections(id int32, uuid string) error {
	return s.cc.ToReviewCreateCollections(id, uuid)
}

func (s *MessageService) ToReviewEditCollections(id int32, uuid string) error {
	return s.cc.ToReviewEditCollections(id, uuid)
}

func (s *MessageService) ColumnCreateReview(ctx context.Context, req *v1.TextReviewReq) (*emptypb.Empty, error) {
	tr := s.TextReview(req)
	err := s.cc.ColumnCreateReview(ctx, tr)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *MessageService) ColumnEditReview(ctx context.Context, req *v1.TextReviewReq) (*emptypb.Empty, error) {
	tr := s.TextReview(req)
	err := s.cc.ColumnEditReview(ctx, tr)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *MessageService) CreateColumnDbCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return s.cc.CreateColumnDbCacheAndSearch(ctx, id, auth, uuid)
}

func (s *MessageService) EditColumnCosAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return s.cc.EditColumnCosAndSearch(ctx, id, auth, uuid)
}

func (s *MessageService) DeleteColumnCacheAndSearch(ctx context.Context, id int32, uuid string) error {
	return s.cc.DeleteColumnCacheAndSearch(ctx, id, uuid)
}

func (s *MessageService) SetColumnViewDbAndCache(ctx context.Context, id int32, uuid string) error {
	return s.cc.SetColumnViewDbAndCache(ctx, id, uuid)
}

func (s *MessageService) SetColumnAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return s.cc.SetColumnAgreeDbAndCache(ctx, id, uuid, userUuid)
}

func (s *MessageService) CancelColumnAgreeDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return s.cc.CancelColumnAgreeDbAndCache(ctx, id, uuid, userUuid)
}

func (s *MessageService) SetColumnCollectDbAndCache(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	return s.cc.SetColumnCollectDbAndCache(ctx, id, collectionsId, uuid, userUuid)
}

func (s *MessageService) CancelColumnCollectDbAndCache(ctx context.Context, id int32, uuid, userUuid string) error {
	return s.cc.CancelColumnCollectDbAndCache(ctx, id, uuid, userUuid)
}

func (s *MessageService) AddColumnIncludesDbAndCache(ctx context.Context, id, articleId int32, uuid string) error {
	return s.cc.AddColumnIncludesDbAndCache(ctx, id, articleId, uuid)
}

func (s *MessageService) DeleteColumnIncludesDbAndCache(ctx context.Context, id, articleId int32, uuid string) error {
	return s.cc.DeleteColumnIncludesDbAndCache(ctx, id, articleId, uuid)
}

func (s *MessageService) SetColumnSubscribeDbAndCache(ctx context.Context, id int32, uuid string) error {
	return s.cc.SetColumnSubscribeDbAndCache(ctx, id, uuid)
}

func (s *MessageService) CancelColumnSubscribeDbAndCache(ctx context.Context, id int32, uuid string) error {
	return s.cc.CancelColumnSubscribeDbAndCache(ctx, id, uuid)
}

func (s *MessageService) CollectionsCreateReview(ctx context.Context, req *v1.TextReviewReq) (*emptypb.Empty, error) {
	tr := s.TextReview(req)
	err := s.cc.CollectionsCreateReview(ctx, tr)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *MessageService) CollectionsEditReview(ctx context.Context, req *v1.TextReviewReq) (*emptypb.Empty, error) {
	tr := s.TextReview(req)
	err := s.cc.CollectionsEditReview(ctx, tr)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *MessageService) CreateCollectionsDbAndCache(ctx context.Context, id, auth int32, uuid string) error {
	return s.cc.CreateCollectionsDbAndCache(ctx, id, auth, uuid)
}

func (s *MessageService) EditCollectionsCos(ctx context.Context, id, auth int32, uuid string) error {
	return s.cc.EditCollectionsCos(ctx, id, auth, uuid)
}

func (s *MessageService) DeleteCollectionsCache(ctx context.Context, id int32, uuid string) error {
	return s.cc.DeleteCollectionsCache(ctx, id, uuid)
}

func (s *MessageService) AddCreationComment(ctx context.Context, createId, createType int32, uuid string) {
	s.cc.AddCreationComment(ctx, createId, createType, uuid)
}
