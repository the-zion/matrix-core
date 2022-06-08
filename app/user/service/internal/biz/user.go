package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/user/service/v1"
)

type User struct {
	Uuid     string
	Password string
	Phone    string
	Email    string
	Wechat   string
	Github   string
}

type UserRepo interface {
	FindByAccount(ctx context.Context, account, mode string) (*User, error)
	GetUser(ctx context.Context, id int64) (*User, error)
	PasswordModify(ctx context.Context, id int64, password string) error
	VerifyPassword(ctx context.Context, id int64, password string) error
	VerifyCode(ctx context.Context, account, code, mode string) error
	SendCode(ctx context.Context, template int64, account, mode string) error
	SetUserPhone(ctx context.Context, id int64, phone string) error
	SetUserEmail(ctx context.Context, id int64, email string) error
}

type UserUseCase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUseCase(repo UserRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "user/biz/userUseCase")),
	}
}

func (r *UserUseCase) SendCode(ctx context.Context, template int64, account, mode string) error {
	err := r.repo.SendCode(ctx, template, account, mode)
	if err != nil {
		return v1.ErrorSendCodeFailed("send code failed: %s", err.Error())
	}
	return nil
}

func (r *UserUseCase) GetUser(ctx context.Context, id int64) (*User, error) {
	user, err := r.repo.GetUser(ctx, id)
	if err != nil {
		return nil, v1.ErrorGetUserFailed("get user failed: %s", err.Error())
	}
	return user, nil
}

func (r *UserUseCase) SetUserPhone(ctx context.Context, id int64, phone, password, code string) error {
	err := r.repo.VerifyCode(ctx, phone, code, "phone")
	if err != nil {
		return v1.ErrorVerifyCodeFailed("set phone failed: %s", err.Error())
	}

	err = r.repo.VerifyPassword(ctx, id, password)
	if err != nil {
		return v1.ErrorVerifyPasswordFailed("set phone failed: %s", err.Error())
	}

	err = r.repo.SetUserPhone(ctx, id, phone)
	if err != nil {
		return v1.ErrorSetPhoneFailed("set phone failed: %s", err.Error())
	}
	return nil
}

func (r *UserUseCase) SetUserEmail(ctx context.Context, id int64, email, password, code string) error {
	err := r.repo.VerifyCode(ctx, email, code, "email")
	if err != nil {
		return v1.ErrorVerifyCodeFailed("set email failed: %s", err.Error())
	}

	err = r.repo.VerifyPassword(ctx, id, password)
	if err != nil {
		return v1.ErrorVerifyPasswordFailed("set phone failed: %s", err.Error())
	}

	err = r.repo.SetUserEmail(ctx, id, email)
	if err != nil {
		return v1.ErrorSetEmailFailed("set email failed: %s", err.Error())
	}
	return nil
}
