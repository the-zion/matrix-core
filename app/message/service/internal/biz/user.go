package biz

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
)

type UserRepo interface {
	SendCode(msgs ...*primitive.MessageExt)
	AvatarIrregular(ctx context.Context, review *AvatarReview, uuid string) error
	UploadProfileToCos(msg *primitive.MessageExt) error
	ProfileReviewPass(ctx context.Context, uuid, update string) error
	ProfileReviewNotPass(ctx context.Context, uuid string) error
	SetFollowDbAndCache(ctx context.Context, uuid, userId string) error
	CancelFollowDbAndCache(ctx context.Context, uuid, userId string) error
}

type UserUseCase struct {
	repo UserRepo
	jwt  Jwt
	log  *log.Helper
}

func NewUserUseCase(repo UserRepo, jwt Jwt, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		repo: repo,
		jwt:  jwt,
		log:  log.NewHelper(log.With(logger, "module", "message/biz/userUseCase")),
	}
}

func (r *UserUseCase) SendCode(msgs ...*primitive.MessageExt) {
	r.repo.SendCode(msgs...)
}

func (r *UserUseCase) UploadProfileToCos(msg *primitive.MessageExt) error {
	return r.repo.UploadProfileToCos(msg)
}

func (r *UserUseCase) AvatarReview(ctx context.Context, ar *AvatarReview) error {
	var err error
	var token string
	var ok bool

	if token, ok = ar.CosHeaders["x-cos-meta-token"]; !ok || token == "" {
		r.log.Info("token not exist，%v", ar)
		return nil
	}

	uuid, err := r.jwt.JwtCheck(token)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get uuid from token: %s", token))
	}

	if ar.State != "Success" {
		r.log.Info("avatar upload review failed，%v", ar)
		return nil
	}

	if ar.Result == 0 {
		return nil
	} else {
		err = r.repo.AvatarIrregular(ctx, ar, uuid)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *UserUseCase) CoverReview(ctx context.Context, cr *CoverReview) error {
	fmt.Println(cr)
	return nil
}

func (r *UserUseCase) ProfileReview(ctx context.Context, tr *TextReview) error {
	var uuid, updated string
	var ok bool

	if uuid, ok = tr.CosHeaders["x-cos-meta-uuid"]; !ok || uuid == "" {
		r.log.Info("uuid not exist，%v", tr)
		return nil
	}

	if updated, ok = tr.CosHeaders["x-cos-meta-update"]; !ok || updated == "" {
		r.log.Info("updated not exist，%v", tr)
		return nil
	}

	if tr.State != "Success" {
		r.log.Info("profile review failed，%v", tr)
		return nil
	}

	var err error
	if tr.Result == 0 {
		err = r.repo.ProfileReviewPass(ctx, uuid, updated)
	} else {
		r.log.Info("profile review not pass，%v", tr)
		err = r.repo.ProfileReviewNotPass(ctx, uuid)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *UserUseCase) SetFollowDbAndCache(ctx context.Context, uuid, userId string) error {
	return r.repo.SetFollowDbAndCache(ctx, uuid, userId)
}

func (r *UserUseCase) CancelFollowDbAndCache(ctx context.Context, uuid, userId string) error {
	return r.repo.CancelFollowDbAndCache(ctx, uuid, userId)
}
