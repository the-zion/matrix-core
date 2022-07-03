package biz

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-kratos/kratos/v2/log"
)

type UserRepo interface {
	SendCode(msgs ...*primitive.MessageExt)
	UploadProfileToCos(msg *primitive.MessageExt) error
	ProfileReviewPass(ctx context.Context, uuid, update string) error
	ProfileReviewNotPass(ctx context.Context, uuid string) error
}

type UserUseCase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUseCase(repo UserRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		repo: repo,
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
	fmt.Println(ar)
	return nil
}

func (r *UserUseCase) ProfileReview(ctx context.Context, tr *TextReview) error {
	uuid := tr.CosHeaders["x-cos-meta-uuid"]
	if uuid == "" {
		r.log.Info("uuid not exist，%v", tr)
		return nil
	}

	updated := tr.CosHeaders["x-cos-meta-update"]
	if updated == "" {
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
