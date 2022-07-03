package service

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	v1 "github.com/the-zion/matrix-core/api/message/service/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *MessageService) ToReviewArticleDraft(msg *primitive.MessageExt) error {
	return s.cc.ToReviewArticleDraft(msg)
}

func (s *MessageService) ArticleDraftReview(ctx context.Context, req *v1.TextReviewReq) (*emptypb.Empty, error) {
	tr := s.TextReview(req)
	err := s.cc.ArticleDraftReview(ctx, tr)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *MessageService) CreateArticleCacheAndSearch(ctx context.Context, uuid string, id int32) error {
	return s.cc.CreateArticleCacheAndSearch(ctx, uuid, id)
}
