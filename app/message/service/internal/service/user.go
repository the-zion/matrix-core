package service

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/the-zion/matrix-core/api/message/service/v1"
	"github.com/the-zion/matrix-core/app/message/service/internal/biz"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *MessageService) SendCode(msgs ...*primitive.MessageExt) {
	s.uc.SendCode(msgs...)
}

func (s *MessageService) UploadProfileToCos(msg *primitive.MessageExt) error {
	return s.uc.UploadProfileToCos(msg)
}

func (s *MessageService) ProfileReview(ctx context.Context, req *v1.TextReviewReq) (*emptypb.Empty, error) {
	tr := s.TextReview(req)
	err := s.uc.ProfileReview(ctx, tr)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *MessageService) AvatarReview(ctx context.Context, req *v1.AvatarReviewReq) (*emptypb.Empty, error) {
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
	return &emptypb.Empty{}, nil
}

func (s *MessageService) CoverReview(ctx context.Context, req *v1.CoverReviewReq) (*emptypb.Empty, error) {
	err := s.uc.CoverReview(ctx, &biz.CoverReview{
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
	return &emptypb.Empty{}, nil
}

func (s *MessageService) SetFollowDbAndCache(ctx context.Context, uuid, userId string) error {
	return s.uc.SetFollowDbAndCache(ctx, uuid, userId)
}

func (s *MessageService) CancelFollowDbAndCache(ctx context.Context, uuid, userId string) error {
	return s.uc.CancelFollowDbAndCache(ctx, uuid, userId)
}
