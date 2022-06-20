package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	messageV1 "github.com/the-zion/matrix-core/api/message/service/v1"
	"github.com/the-zion/matrix-core/app/bff/interface/internal/biz"
	"golang.org/x/sync/singleflight"
)

var _ biz.MessageRepo = (*messageRepo)(nil)

type messageRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

func NewMessageRepo(data *Data, logger log.Logger) biz.MessageRepo {
	return &messageRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "bff/data/message")),
		sg:   &singleflight.Group{},
	}
}

type AvatarReview struct {
	Code       string
	Message    string
	JobId      string
	State      string
	Object     string
	Label      string
	Result     int32
	Category   string
	BucketId   string
	Region     string
	CosHeaders map[string]string
	EventName  string
}

func (r *messageRepo) AvatarReview(ctx context.Context, ar *biz.AvatarReview) error {
	arr := &messageV1.AvatarReviewReq{
		JobsDetail: &messageV1.AvatarReviewReq_JobsDetailStruct{},
	}
	arr.JobsDetail.Code = ar.Code
	arr.JobsDetail.Message = ar.Message
	arr.JobsDetail.JobId = ar.JobId
	arr.JobsDetail.State = ar.State
	arr.JobsDetail.Object = ar.Object
	arr.JobsDetail.Label = ar.Label
	arr.JobsDetail.Result = ar.Result
	arr.JobsDetail.Category = ar.Category
	arr.JobsDetail.BucketId = ar.BucketId
	arr.JobsDetail.Region = ar.Region
	arr.JobsDetail.CosHeaders = ar.CosHeaders
	arr.EventName = ar.EventName

	_, err := r.data.mc.AvatarReview(ctx, arr)
	if err != nil {
		return err
	}
	return nil
}
