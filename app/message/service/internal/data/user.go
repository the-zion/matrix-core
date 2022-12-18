package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/tencentyun/cos-go-sdk-v5"
	userV1 "github.com/the-zion/matrix-core/api/user/service/v1"
	"github.com/the-zion/matrix-core/app/message/service/internal/biz"
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

func (r *userRepo) UploadProfileToCos(msg map[string]interface{}) error {
	key := "profile/" + msg["uuid"].(string)

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: "text/html",
			XCosMetaXXX: &http.Header{},
		},
	}

	opt.XCosMetaXXX.Add("x-cos-meta-uuid", msg["uuid"].(string))
	opt.XCosMetaXXX.Add("x-cos-meta-update", msg["updated"].(string))

	m, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to marshal profile message: profile(%v)", msg))
	}

	f := strings.NewReader(string(m))
	_, err = r.data.cosUserCli.cos.Object.Put(
		context.Background(), key, f, opt,
	)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to upload profile to cos: profile(%v)", msg))
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

func (r *userRepo) AvatarIrregular(ctx context.Context, review *biz.ImageReview, uuid string) error {
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

func (r *userRepo) CoverIrregular(ctx context.Context, review *biz.ImageReview, uuid string) error {
	_, err := r.data.uc.CoverIrregular(ctx, &userV1.CoverIrregularReq{
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

func (r *userRepo) AddAvatarReviewDbAndCache(ctx context.Context, score, result int32, uuid, jobId, label, category, subLabel string) error {
	_, err := r.data.uc.AddAvatarReviewDbAndCache(ctx, &userV1.AddAvatarReviewDbAndCacheReq{
		Uuid:     uuid,
		Score:    score,
		JobId:    jobId,
		Label:    label,
		Result:   result,
		Category: category,
		SubLabel: subLabel,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) AddCoverReviewDbAndCache(ctx context.Context, score, result int32, uuid, jobId, label, category, subLabel string) error {
	_, err := r.data.uc.AddCoverReviewDbAndCache(ctx, &userV1.AddCoverReviewDbAndCacheReq{
		Uuid:     uuid,
		Score:    score,
		JobId:    jobId,
		Label:    label,
		Result:   result,
		Category: category,
		SubLabel: subLabel,
	})
	if err != nil {
		return err
	}
	return nil
}
