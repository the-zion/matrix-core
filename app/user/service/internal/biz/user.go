package biz

import (
	"context"
	"errors"
	v1 "github.com/Cube-v2/cube-core/api/user/service/v1"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrPasswordError = errors.New("password error")
	ErrSendCodeError = errors.New("send code error")
	ErrCodeError     = errors.New("code error")
	ErrUnknownError  = errors.New("unknown error")
)

type User struct {
	Id       int64
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
	SendCode(ctx context.Context, template int64, account, mode string) (string, error)
	SetUserPhone(ctx context.Context, id int64) (string, error)
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

func (r *UserUseCase) SendCode(ctx context.Context, template int64, account, mode string) (string, error) {
	code, err := r.repo.SendCode(ctx, template, account, mode)
	if errors.Is(err, ErrUnknownError) {
		return "", v1.ErrorUnknownError("send code failed: %s", err.Error())
	}
	if errors.Is(err, ErrSendCodeError) {
		return "", v1.ErrorSendCodeFailed("send code failed: %s", err.Error())
	}

	return code, nil
}

func (r *UserUseCase) GetUser(ctx context.Context, id int64) (*User, error) {
	user, err := r.repo.GetUser(ctx, id)
	if err != nil {
		return nil, v1.ErrorGetUserFailed("get user failed: %s", err.Error())
	}
	return user, nil
}

func (r *UserUseCase) SetUserPhone(ctx context.Context, id int64, phone, code string) (string, error) {
	phone, err := r.repo.SetUserPhone(ctx, id)
	if err != nil {
		return "", v1.ErrorGetUserFailed("get user failed: %s", err.Error())
	}
	return phone, nil
}
