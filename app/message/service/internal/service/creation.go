package service

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	v1 "github.com/the-zion/matrix-core/api/message/service/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *MessageService) ToReviewArticleDraft(msgs ...*primitive.MessageExt) {
	s.cc.ToReviewArticleDraft(msgs...)
}

func (s *MessageService) ArticleDraftReview(ctx context.Context, req *v1.TextReviewReq) (*emptypb.Empty, error) {
	tr := s.TextReview(req)
	err := s.cc.ArticleDraftReview(ctx, tr)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
