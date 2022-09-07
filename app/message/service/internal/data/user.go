package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentyun/cos-go-sdk-v5"
	userV1 "github.com/the-zion/matrix-core/api/user/service/v1"
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

func (r *userRepo) UploadProfileToCos(msg *primitive.MessageExt) error {
	m := map[string]interface{}{"Uuid": ""}
	err := json.Unmarshal(msg.Body, &m)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to unmarshal profile: profile(%v)", msg.Body))
	}
	key := "profile/" + m["Uuid"].(string)

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: "text/html",
			XCosMetaXXX: &http.Header{},
		},
	}

	opt.XCosMetaXXX.Add("x-cos-meta-uuid", m["Uuid"].(string))
	opt.XCosMetaXXX.Add("x-cos-meta-update", m["Updated"].(string))

	f := strings.NewReader(string(msg.Body))
	_, err = r.data.cosUserCli.cos.Object.Put(
		context.Background(), key, f, opt,
	)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to upload profile to cos: profile(%v)", msg.Body))
	}
	return nil
}

func (r *userRepo) ProfileReviewPass(ctx context.Context, uuid, update string) error {
	_, err := r.data.uc.ProfileReviewPass(ctx, &userV1.ProfileReviewPassReq{
		Uuid:   uuid,
		Update: update,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) ProfileReviewNotPass(ctx context.Context, uuid string) error {
	_, err := r.data.uc.ProfileReviewNotPass(ctx, &userV1.ProfileReviewNotPassReq{
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
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

func (r *userRepo) AvatarIrregular(ctx context.Context, review *biz.AvatarReview, uuid string) error {
	_, err := r.data.uc.AvatarIrregular(ctx, &userV1.AvatarIrregularReq{
		Uuid:     uuid,
		JobId:    review.JobId,
		Url:      review.Url,
		Label:    review.Label,
		Result:   review.Result,
		Score:    review.Score,
		Category: review.Category,
		SubLabel: review.SubLabel,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) SetFollowDbAndCache(ctx context.Context, uuid, userId string) error {
	_, err := r.data.uc.SetFollowDbAndCache(ctx, &userV1.SetFollowDbAndCacheReq{
		Uuid:     uuid,
		UserUuid: userId,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) CancelFollowDbAndCache(ctx context.Context, uuid, userId string) error {
	_, err := r.data.uc.CancelFollowDbAndCache(ctx, &userV1.CancelFollowDbAndCacheReq{
		Uuid:     uuid,
		UserUuid: userId,
	})
	if err != nil {
		return err
	}
	return nil
}
