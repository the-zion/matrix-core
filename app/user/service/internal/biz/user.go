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
)

type User struct {
}

type UserRepo interface {
	FindByUserAccount(ctx context.Context, account, mode string) (*Auth, error)
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
		log:  log.NewHelper(log.With(logger, "module", "usecase/user")),
	}
}

func (receiver *UserUseCase) SendCode(ctx context.Context, req *v1.SendCodeReq) (*v1.SendCodeReply, error) {
	if len(req.Account) > 0 && len(req.Mode) > 0 && (req.Mode == "phone" || req.Mode == "email") && (req.Template > 0 && req.Template < 4) {
		code, err := receiver.repo.SendCode(ctx, req.Template, req.Account, req.Mode)
		if err != nil {
			return nil, v1.ErrorSendCodeFailed("send code failed: %s", err.Error())
		}
		return &v1.SendCodeReply{
			Code: code,
		}, nil
	}
	return nil, v1.ErrorSendCodeFailed("send code failed: params illegal")
}
