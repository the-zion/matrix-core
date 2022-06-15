package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/user/service/v1"
)

type ProfileRepo interface {
	SetProfile(ctx context.Context, id int64, sex, introduce, industry, address, profile, tag string) error
	SetName(ctx context.Context, id int64, name string) error
}

type ProfileUseCase struct {
	repo ProfileRepo
	log  *log.Helper
}

func NewProfileUseCase(repo ProfileRepo, logger log.Logger) *ProfileUseCase {
	return &ProfileUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "user/biz/profileUseCase")),
	}
}

func (r *ProfileUseCase) SetUserProfile(ctx context.Context, id int64, sex, introduce, industry, address, profile, tag string) error {
	err := r.repo.SetProfile(ctx, id, sex, introduce, industry, address, profile, tag)
	if err != nil {
		return v1.ErrorSetProfileFailed("set user profile failed: %s", err.Error())
	}
	return nil
}

func (r *ProfileUseCase) SetUserName(ctx context.Context, id int64, name string) error {
	err := r.repo.SetName(ctx, id, name)
	if err != nil {
		return v1.ErrorSetUsernameFailed("set user name failed: %s", err.Error())
	}
	return nil
}
