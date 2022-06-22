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

func (s *MessageService) UploadProfileToCos(msgs ...*primitive.MessageExt) {
	s.uc.UploadProfileToCos(msgs...)
}

func (s *MessageService) ProfileReview(ctx context.Context, req *v1.ProfileReviewReq) (*emptypb.Empty, error) {
	tr := &biz.TextReview{
		Code:         req.JobsDetail.Code,
		Message:      req.JobsDetail.Message,
		JobId:        req.JobsDetail.JobId,
		DataId:       req.JobsDetail.DataId,
		State:        req.JobsDetail.State,
		CreationTime: req.JobsDetail.CreationTime,
		Object:       req.JobsDetail.Object,
		Label:        req.JobsDetail.Label,
		Result:       req.JobsDetail.Result,
		BucketId:     req.JobsDetail.BucketId,
		Region:       req.JobsDetail.Region,
		CosHeaders:   req.JobsDetail.CosHeaders,
	}

	var section []*biz.Section

	for _, item := range req.JobsDetail.Section {
		se := &biz.Section{
			Label:  item.Label,
			Result: item.Result,
			PornInfo: &biz.SectionPornInfo{
				HitFlag:  item.PornInfo.HitFlag,
				Score:    item.PornInfo.Score,
				Keywords: item.PornInfo.Keywords,
			},
			AdsInfo: &biz.SectionAdsInfo{
				HitFlag:  item.AdsInfo.HitFlag,
				Score:    item.AdsInfo.Score,
				Keywords: item.AdsInfo.Keywords,
			},
			IllegalInfo: &biz.SectionIllegalInfo{
				HitFlag:  item.IllegalInfo.HitFlag,
				Score:    item.IllegalInfo.Score,
				Keywords: item.IllegalInfo.Keywords,
			},
			AbuseInfo: &biz.SectionAbuseInfo{
				HitFlag:  item.AbuseInfo.HitFlag,
				Score:    item.AbuseInfo.Score,
				Keywords: item.AbuseInfo.Keywords,
			},
		}
		section = append(section, se)
	}

	tr.Section = section

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
