package biz

import (
	"context"
	"errors"
	"github.com/Cube-v2/cube-core/app/user/service/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-playground/validator/v10"
)

var (
	ErrGetProfileFailed = errors.New("get user profile failed")
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
	//GetProfile(ctx context.Context, id int64) (*User, error)
}

type ProfileUseCase struct {
	key      string
	repo     ProfileRepo
	validate *validator.Validate
	log      *log.Helper
}

func NewProfileUseCase(conf *conf.Auth, repo ProfileRepo, validator *validator.Validate, logger log.Logger) *ProfileUseCase {
	return &ProfileUseCase{
		key:      conf.ApiKey,
		repo:     repo,
		validate: validator,
		log:      log.NewHelper(log.With(logger, "module", "user/biz/profileUseCase")),
	}
}

func (r *ProfileUseCase) GetUserProfile(ctx context.Context, id int64) (*Profile, error) {
	return &Profile{}, nil
}
