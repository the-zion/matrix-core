package service

import (
	"context"
	"github.com/the-zion/matrix-core/api/message/service/v1"
	"github.com/the-zion/matrix-core/app/message/service/internal/biz"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *MessageService) UploadProfileToCos(msg map[string]interface{}) error {
	return s.uc.UploadProfileToCos(msg)
}

func (s *MessageService) ProfileReview(ctx context.Context, req *v1.TextReviewReq) (*emptypb.Empty, error) {
	tr, err := s.TextReview(req)
	if err != nil {
		return nil, err
	}

	err = s.uc.ProfileReview(ctx, tr)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *MessageService) AvatarReview(ctx context.Context, req *v1.ImageReviewReq) (*emptypb.Empty, error) {
	err := s.uc.AvatarReview(ctx, &biz.ImageReview{
		Code:       req.JobsDetail.Code,
		Message:    req.JobsDetail.Message,
		JobId:      req.JobsDetail.JobId,
		State:      req.JobsDetail.State,
		Object:     req.JobsDetail.Object,
		Url:        req.JobsDetail.Url,
		Label:      req.JobsDetail.Label,
		Result:     req.JobsDetail.Result,
		Score:      req.JobsDetail.Score,
		Category:   req.JobsDetail.Category,
		SubLabel:   req.JobsDetail.SubLabel,
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

func (s *MessageService) CoverReview(ctx context.Context, req *v1.ImageReviewReq) (*emptypb.Empty, error) {
	err := s.uc.CoverReview(ctx, &biz.ImageReview{
		Code:       req.JobsDetail.Code,
		Message:    req.JobsDetail.Message,
		JobId:      req.JobsDetail.JobId,
		State:      req.JobsDetail.State,
		Object:     req.JobsDetail.Object,
		Url:        req.JobsDetail.Url,
		Label:      req.JobsDetail.Label,
		Result:     req.JobsDetail.Result,
		Score:      req.JobsDetail.Score,
		Category:   req.JobsDetail.Category,
		SubLabel:   req.JobsDetail.SubLabel,
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

func (s *MessageService) AddAvatarReviewDbAndCache(ctx context.Context, score, result int32, uuid, jobId, label, category, subLabel string) error {
	return s.uc.AddAvatarReviewDbAndCache(ctx, score, result, uuid, jobId, label, category, subLabel)
}

func (s *MessageService) AddCoverReviewDbAndCache(ctx context.Context, score, result int32, uuid, jobId, label, category, subLabel string) error {
	return s.uc.AddCoverReviewDbAndCache(ctx, score, result, uuid, jobId, label, category, subLabel)
}
