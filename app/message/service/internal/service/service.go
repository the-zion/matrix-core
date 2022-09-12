package service

import (
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "github.com/the-zion/matrix-core/api/message/service/v1"
	"github.com/the-zion/matrix-core/app/message/service/internal/biz"
)

var ProviderSet = wire.NewSet(NewMessageService)

type MessageService struct {
	v1.UnimplementedMessageServer
	uc    *biz.UserUseCase
	cc    *biz.CreationUseCase
	ac    *biz.AchievementCase
	commc *biz.CommentUseCase
	log   *log.Helper
}

func NewMessageService(uc *biz.UserUseCase, cc *biz.CreationUseCase, ac *biz.AchievementCase, commc *biz.CommentUseCase, logger log.Logger) *MessageService {
	return &MessageService{
		log:   log.NewHelper(log.With(logger, "module", "message/service")),
		uc:    uc,
		cc:    cc,
		ac:    ac,
		commc: commc,
	}
}

func (s *MessageService) TextReview(req *v1.TextReviewReq) (*biz.TextReview, error) {
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

	var section []map[string]interface{}

	for _, item := range req.JobsDetail.Section {
		se := make(map[string]interface{}, 0)
		se["PornInfoHitFlag"] = item.PornInfo.HitFlag
		se["PornInfoKeywords"] = item.PornInfo.Keywords
		se["AdsInfoHitFlag"] = item.AdsInfo.HitFlag
		se["AdsInfoKeywords"] = item.AdsInfo.Keywords
		se["IllegalInfoHitFlag"] = item.IllegalInfo.HitFlag
		se["IllegalInfoKeywords"] = item.IllegalInfo.Keywords
		se["AbuseInfoHitFlag"] = item.AbuseInfo.HitFlag
		se["AbuseInfoKeywords"] = item.AbuseInfo.Keywords
		section = append(section, se)
	}

	sectionMap, err := json.Marshal(section)
	if err != nil {
		return nil, err
	}

	tr.Section = string(sectionMap)
	return tr, nil
}
