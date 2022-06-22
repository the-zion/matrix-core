package data

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/the-zion/matrix-core/app/message/service/internal/biz"
	"github.com/the-zion/matrix-core/app/message/service/internal/pkg/util"
	"net/http"
	"strings"
)

var _ biz.UserRepo = (*userRepo)(nil)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "message/data/user")),
	}
}

func (r *userRepo) UploadProfileToCos(msgs ...*primitive.MessageExt) {
	for _, i := range msgs {
		m := map[string]string{"Uuid": ""}
		err := json.Unmarshal(i.Body, &m)
		if err != nil {
			log.Errorf("fail to unmarshal profile: err(%v)", err)
		}
		key := "profile/" + m["Uuid"]

		opt := &cos.ObjectPutOptions{
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				ContentType: "text/html",
				XCosMetaXXX: &http.Header{},
			},
		}

		opt.XCosMetaXXX.Add("x-cos-meta-uuid", m["Uuid"])

		f := strings.NewReader(string(i.Body))
		_, err = r.data.cosCli.Object.Put(
			context.Background(), key, f, opt,
		)
		if err != nil {
			log.Errorf("fail to upload profile to cos: err(%v)", err)
		}
	}
}

func (r *userRepo) SendCode(msgs ...*primitive.MessageExt) {
	for _, i := range msgs {
		body := strings.Split(string(i.Body), ";")
		if body[3] == "phone" {
			request := r.data.phoneCodeCli.request
			client := r.data.phoneCodeCli.client
			request.TemplateId = common.StringPtr(util.GetPhoneTemplate(body[2]))
			request.TemplateParamSet = common.StringPtrs([]string{body[1]})
			request.PhoneNumberSet = common.StringPtrs([]string{body[0]})
			_, err := client.SendSms(request)
			if err != nil {
				r.log.Errorf("fail to send phone code: code(%s) error: %v", body[1], err.Error())
			}
		}

		if body[3] == "email" {
			m := r.data.goMailCli.message
			d := r.data.goMailCli.dialer
			m.SetHeader("To", body[0])
			m.SetHeader("Subject", "matrix 魔方技术")
			m.SetBody("text/html", util.GetEmailTemplate(body[2], body[1]))
			err := d.DialAndSend(m)
			if err != nil {
				r.log.Errorf("fail to send email code: code(%s) error: %v", body[1], err.Error())
			}
		}
	}
}

//func (r *userRepo) ProfileReview(msgs ...*primitive.MessageExt) error {
//	tr := &biz.TextReview{
//		Code:         req.JobsDetail.Code,
//		Message:      req.JobsDetail.Message,
//		JobId:        req.JobsDetail.JobId,
//		DataId:       req.JobsDetail.DataId,
//		State:        req.JobsDetail.State,
//		CreationTime: req.JobsDetail.CreationTime,
//		Object:       req.JobsDetail.Object,
//		Label:        req.JobsDetail.Label,
//		Result:       req.JobsDetail.Result,
//		BucketId:     req.JobsDetail.BucketId,
//		Region:       req.JobsDetail.Region,
//		CosHeaders:   req.JobsDetail.CosHeaders,
//	}
//
//	var section []*biz.Section
//
//	for _, item := range req.JobsDetail.Section {
//		se := &biz.Section{
//			Label:  item.Label,
//			Result: item.Result,
//		}
//
//		if item.PornInfo != nil {
//			se.PornInfo = &biz.SectionPornInfo{
//				HitFlag:  item.PornInfo.HitFlag,
//				Score:    item.PornInfo.Score,
//				Keywords: item.PornInfo.Keywords,
//			}
//		}
//
//		if item.AdsInfo != nil {
//			se.AdsInfo = &biz.SectionAdsInfo{
//				HitFlag:  item.AdsInfo.HitFlag,
//				Score:    item.AdsInfo.Score,
//				Keywords: item.AdsInfo.Keywords,
//			}
//		}
//
//		if item.IllegalInfo != nil {
//			se.IllegalInfo = &biz.SectionIllegalInfo{
//				HitFlag:  item.IllegalInfo.HitFlag,
//				Score:    item.IllegalInfo.Score,
//				Keywords: item.IllegalInfo.Keywords,
//			}
//		}
//
//		if item.AbuseInfo != nil {
//			se.AbuseInfo = &biz.SectionAbuseInfo{
//				HitFlag:  item.AbuseInfo.HitFlag,
//				Score:    item.AbuseInfo.Score,
//				Keywords: item.AbuseInfo.Keywords,
//			}
//		}
//		section = append(section, se)
//
//	}
//
//	tr.Section = section
//}
