package biz

import (
	"context"
	v1 "github.com/Cube-v2/cube-core/api/user/service/v1"
	"github.com/go-kratos/kratos/v2/log"
)

type Profile struct {
	Username        string
	Sex             string
	Introduce       string
	Industry        string
	Address         string
	PersonalProfile string
	Tag             string
	Background      string
	Image           string
}

type ProfileRepo interface {
	GetProfile(ctx context.Context, id int64) (*Profile, error)
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

func (r *ProfileUseCase) GetUserProfile(ctx context.Context, id int64) (*Profile, error) {
	profile, err := r.repo.GetProfile(ctx, id)
	if err != nil {
		return nil, v1.ErrorGetProfileFailed("get user profile failed: %s", err.Error())
	}
	return &Profile{
		Username:        profile.Username,
		Sex:             profile.Sex,
		Introduce:       profile.Introduce,
		Industry:        profile.Industry,
		Address:         profile.Address,
		PersonalProfile: profile.PersonalProfile,
		Tag:             profile.Tag,
		Background:      profile.Background,
		Image:           profile.Image,
	}, nil
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
