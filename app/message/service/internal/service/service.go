package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "github.com/the-zion/matrix-core/api/message/service/v1"
	"github.com/the-zion/matrix-core/app/message/service/internal/biz"
)

var ProviderSet = wire.NewSet(NewMessageService)

type MessageService struct {
	v1.UnimplementedMessageServer
	uc  *biz.UserUseCase
	cc  *biz.CreationUseCase
	ac  *biz.AchievementCase
	log *log.Helper
}

func NewMessageService(uc *biz.UserUseCase, cc *biz.CreationUseCase, ac *biz.AchievementCase, logger log.Logger) *MessageService {
	return &MessageService{
		log: log.NewHelper(log.With(logger, "module", "message/service")),
		uc:  uc,
		cc:  cc,
		ac:  ac,
	}
}

func (s *MessageService) TextReview(req *v1.TextReviewReq) *biz.TextReview {
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
	return tr
}
