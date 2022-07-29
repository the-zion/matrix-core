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

func (s *MessageService) CreateArticleCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return s.cc.CreateArticleCacheAndSearch(ctx, id, auth, uuid)
}

func (s *MessageService) EditArticleCosAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return s.cc.EditArticleCosAndSearch(ctx, id, auth, uuid)
}

func (s *MessageService) DeleteArticleCacheAndSearch(ctx context.Context, id int32, uuid string) error {
	return s.cc.DeleteArticleCacheAndSearch(ctx, id, uuid)
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

func (s *MessageService) CreateTalkCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return s.cc.CreateTalkCacheAndSearch(ctx, id, auth, uuid)
}

func (s *MessageService) EditTalkCosAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return s.cc.EditTalkCosAndSearch(ctx, id, auth, uuid)
}

func (s *MessageService) DeleteTalkCacheAndSearch(ctx context.Context, id int32, uuid string) error {
	return s.cc.DeleteTalkCacheAndSearch(ctx, id, uuid)
}

func (s *MessageService) ToReviewCreateColumn(id int32, uuid string) error {
	return s.cc.ToReviewCreateColumn(id, uuid)
}

func (s *MessageService) ToReviewEditColumn(id int32, uuid string) error {
	return s.cc.ToReviewEditColumn(id, uuid)
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

func (s *MessageService) CreateColumnCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return s.cc.CreateColumnCacheAndSearch(ctx, id, auth, uuid)
}

func (s *MessageService) EditColumnCosAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	return s.cc.EditColumnCosAndSearch(ctx, id, auth, uuid)
}

func (s *MessageService) DeleteColumnCacheAndSearch(ctx context.Context, id int32, uuid string) error {
	return s.cc.DeleteColumnCacheAndSearch(ctx, id, uuid)
}
