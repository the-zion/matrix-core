package biz

import (
	"context"
	"errors"
	v1 "github.com/Cube-v2/cube-core/api/user/service/v1"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrProfileNotFound = errors.New("profile not found")
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
