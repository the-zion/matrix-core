package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type UserRepo interface {
	UserRegister(ctx context.Context, email, password, code string) error
	LoginByPassword(ctx context.Context, account, password, mode string) (string, error)
	LoginByCode(ctx context.Context, phone, code string) (string, error)
	LoginPasswordReset(ctx context.Context, account, password, code, mode string) error
	SendPhoneCode(ctx context.Context, template, phone string) error
	SendEmailCode(ctx context.Context, template, email string) error
	GetCosSessionKey(ctx context.Context) (*Credentials, error)
	GetAccount(ctx context.Context, uuid string) (*UserAccount, error)
	GetProfile(ctx context.Context, uuid string) (*UserProfile, error)
	GetProfileUpdate(ctx context.Context, uuid string) (*UserProfileUpdate, error)
	SetProfileUpdate(ctx context.Context, profile *UserProfileUpdate) error
	SetUserPhone(ctx context.Context, uuid, phone, code string) error
	SetUserEmail(ctx context.Context, uuid, email, code string) error
}

type UserUseCase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUseCase(repo UserRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "bff/biz/UserUseCase")),
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

func (r *UserUseCase) GetAccount(ctx context.Context) (*UserAccount, error) {
	uuid := ctx.Value("uuid").(string)
	account, err := r.repo.GetAccount(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (r *UserUseCase) GetProfile(ctx context.Context) (*UserProfile, error) {
	uuid := ctx.Value("uuid").(string)
	userProfile, err := r.repo.GetProfile(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return userProfile, nil
}

func (r *UserUseCase) GetProfileUpdate(ctx context.Context) (*UserProfileUpdate, error) {
	uuid := ctx.Value("uuid").(string)
	userProfile, err := r.repo.GetProfileUpdate(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return userProfile, nil
}

func (r *UserUseCase) SetUserProfile(ctx context.Context, profile *UserProfileUpdate) error {
	uuid := ctx.Value("uuid").(string)
	profile.Uuid = uuid
	err := r.repo.SetProfileUpdate(ctx, profile)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserUseCase) SetUserPhone(ctx context.Context, phone, code string) error {
	uuid := ctx.Value("uuid").(string)
	err := r.repo.SetUserPhone(ctx, uuid, phone, code)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserUseCase) SetUserEmail(ctx context.Context, email, code string) error {
	uuid := ctx.Value("uuid").(string)
	err := r.repo.SetUserEmail(ctx, uuid, email, code)
	if err != nil {
		return err
	}
	return nil
}
