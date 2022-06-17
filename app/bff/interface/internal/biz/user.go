package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/bff/interface/v1"
)

type UserRepo interface {
	UserRegister(ctx context.Context, email, password, code string) error
	LoginByPassword(ctx context.Context, account, password, mode string) (string, error)
	LoginByCode(ctx context.Context, phone, code string) (string, error)
	LoginPasswordReset(ctx context.Context, account, password, code, mode string) error
	SendPhoneCode(ctx context.Context, template, phone string) error
	SendEmailCode(ctx context.Context, template, email string) error
	GetCosSessionKey(ctx context.Context) (*Credentials, error)
	GetUserProfile(ctx context.Context, uuid string) (*UserProfile, error)
	GetUserProfileUpdate(ctx context.Context, uuid string) (*UserProfileUpdate, error)
}

type UserUseCase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUseCase(repo UserRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "info/biz/userUseCase")),
	}
}

func (r *UserUseCase) UserRegister(ctx context.Context, email, password, code string) error {
	err := r.repo.UserRegister(ctx, email, password, code)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserUseCase) LoginByPassword(ctx context.Context, account, password, mode string) (string, error) {
	token, err := r.repo.LoginByPassword(ctx, account, password, mode)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *UserUseCase) LoginByCode(ctx context.Context, phone, code string) (string, error) {
	token, err := r.repo.LoginByCode(ctx, phone, code)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *UserUseCase) LoginPasswordReset(ctx context.Context, account, password, code, mode string) error {
	err := r.repo.LoginPasswordReset(ctx, account, password, code, mode)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserUseCase) SendPhoneCode(ctx context.Context, template, phone string) error {
	err := r.repo.SendPhoneCode(ctx, template, phone)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserUseCase) SendEmailCode(ctx context.Context, template, email string) error {
	err := r.repo.SendEmailCode(ctx, template, email)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserUseCase) GetCosSessionKey(ctx context.Context) (*Credentials, error) {
	credentials, err := r.repo.GetCosSessionKey(ctx)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}

func (r *UserUseCase) GetUserProfile(ctx context.Context) (*UserProfile, error) {
	uuid := ctx.Value("uuid").(string)
	userProfile, err := r.repo.GetUserProfile(ctx, uuid)
	if err != nil {
		return nil, v1.ErrorGetUserProfileFailed("get user profile failed: %s", err.Error())
	}
	return userProfile, nil
}

func (r *UserUseCase) GetUserProfileUpdate(ctx context.Context) (*UserProfileUpdate, error) {
	uuid := ctx.Value("uuid").(string)
	userProfile, err := r.repo.GetUserProfileUpdate(ctx, uuid)
	if err != nil {
		return nil, v1.ErrorGetUserProfileFailed("get user profile failed: %s", err.Error())
	}
	return userProfile, nil
}
