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
	Username string
	Password string
	Phone    string
	Email    string
	Wechat   string
	Github   string
}

type UserRepo interface {
	FindByAccount(ctx context.Context, account, mode string) (*User, error)
	GetUser(ctx context.Context, username string) (*User, error)
	PasswordModify(ctx context.Context, id int64, password string) error
	VerifyPassword(ctx context.Context, id int64, password string) error
	VerifyCode(ctx context.Context, account, code, mode string) error
	SendCode(ctx context.Context, template int64, account, mode string) (string, error)
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

func (receiver *UserUseCase) SendCode(ctx context.Context, req *v1.SendCodeReq) (*v1.SendCodeReply, error) {
	if len(req.Account) > 0 && len(req.Mode) > 0 && (req.Mode == "phone" || req.Mode == "email") && (req.Template > 0 && req.Template < 4) {
		code, err := receiver.repo.SendCode(ctx, req.Template, req.Account, req.Mode)
		if errors.Is(err, ErrUnknownError) {
			return nil, v1.ErrorUnknownError("send code failed: %s", err.Error())
		}
		if errors.Is(err, ErrSendCodeError) {
			return nil, v1.ErrorSendCodeFailed("send code failed: %s", err.Error())
		}

		return &v1.SendCodeReply{
			Code: code,
		}, nil
	}
	return nil, v1.ErrorParamsIllegal("send code failed: params illegal")
}

func (receiver *UserUseCase) GetUser(ctx context.Context, req *v1.GetUserReq) (*v1.GetUserReply, error) {
	if len(req.Username) > 0 {
		user, err := receiver.repo.GetUser(ctx, req.Username)
		if err != nil {
			return nil, v1.ErrorGetUserFailed("get user failed: %s", err.Error())
		}
		return &v1.GetUserReply{
			Phone:  user.Phone,
			Email:  user.Email,
			Wechat: user.Wechat,
			Github: user.Github,
		}, nil
	}
	return nil, v1.ErrorParamsIllegal("get user failed: params illegal")
}
