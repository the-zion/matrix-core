package service

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/the-zion/matrix-core/api/message/service/v1"
	"github.com/the-zion/matrix-core/app/message/service/internal/biz"
)

func (s *MessageService) SendCode(msgs ...*primitive.MessageExt) {
	s.uc.SendCode(msgs...)
}

func (s *MessageService) ProfileReview(ctx context.Context, req *v1.ProfileReviewReq) (*v1.ProfileReviewReply, error) {
	fmt.Println(req)
	return &v1.ProfileReviewReply{}, nil
}

func (s *MessageService) AvatarReview(ctx context.Context, req *v1.AvatarReviewReq) (*v1.AvatarReviewReply, error) {
	err := s.uc.AvatarReview(ctx, &biz.AvatarReview{
		Code:       req.JobsDetail.Code,
		Message:    req.JobsDetail.Message,
		JobId:      req.JobsDetail.JobId,
		State:      req.JobsDetail.State,
		Object:     req.JobsDetail.Object,
		Label:      req.JobsDetail.Label,
		Result:     req.JobsDetail.Result,
		Category:   req.JobsDetail.Category,
		BucketId:   req.JobsDetail.BucketId,
		Region:     req.JobsDetail.Region,
		CosHeaders: req.JobsDetail.CosHeaders,
		EventName:  req.EventName,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AvatarReviewReply{}, nil
}
