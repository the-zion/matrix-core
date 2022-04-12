package biz

import (
	"context"
	"errors"
	v1 "github.com/Cube-v2/cube-core/api/user/service/v1"
	"github.com/Cube-v2/cube-core/app/user/service/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrUserRegisterFailed = errors.New("user register failed")
)

type Auth struct {
	Id       int64
	Username string
	Password string
	Phone    string
	Email    string
	Wechat   string
	Github   string
}

type AuthRepo interface {
	UserRegister(ctx context.Context, u *Auth, mode string) (*Auth, error)
}

type AuthUseCase struct {
	key      string
	repo     AuthRepo
	userRepo UserRepo
	log      *log.Helper
}

func NewAuthUseCase(conf *conf.Auth, repo AuthRepo, userRepo UserRepo, logger log.Logger) *AuthUseCase {
	return &AuthUseCase{
		key:      conf.ApiKey,
		repo:     repo,
		userRepo: userRepo,
		log:      log.NewHelper(log.With(logger, "module", "usecase/auth")),
	}
}

func (receiver *AuthUseCase) Login(ctx context.Context, req *v1.LoginReq) (*v1.LoginReply, error) {
	if len(req.Account) > 0 && len(req.Password) > 0 && len(req.Mode) > 0 {
		return loginByAccountAndPassword(ctx, receiver, req.Account, req.Password, req.Mode)
	}
	if len(req.Account) > 0 && len(req.Code) > 0 && len(req.Mode) > 0 {
		return loginByAccountAndCode(ctx, receiver, req.Account, req.Code, req.Mode)
	}
	return &v1.LoginReply{}, v1.ErrorLoginFailed("login failed: params illegal")
}

func loginByAccountAndPassword(ctx context.Context, receiver *AuthUseCase, account, password, mode string) (*v1.LoginReply, error) {
	user, err := receiver.userRepo.FindByUserAccount(ctx, account, mode)
	if err != nil {
		return nil, v1.ErrorLoginFailed("login failed: %s", err.Error())
	}

	err = receiver.userRepo.VerifyPassword(ctx, user.Id, password)
	if err != nil {
		return nil, v1.ErrorLoginFailed("login failed: %s", err.Error())
	}

	return signToken(user.Id, receiver)
}

func loginByAccountAndCode(ctx context.Context, receiver *AuthUseCase, account, code, mode string) (*v1.LoginReply, error) {
	err := receiver.userRepo.VerifyCode(ctx, account, code, mode)
	if err != nil {
		return nil, v1.ErrorLoginFailed("login failed: %s", err.Error())
	}

	user, err := receiver.userRepo.FindByUserAccount(ctx, account, mode)
	if err != nil {
		user, err = receiver.Register(ctx, account, mode)
		if err != nil {
			return nil, v1.ErrorLoginFailed("login failed: %s", err.Error())
		}
	}

	return signToken(user.Id, receiver)
}

func (receiver *AuthUseCase) Register(ctx context.Context, account, mode string) (*Auth, error) {
	var user *Auth
	var err error
	switch mode {
	case "phone":
		user, err = receiver.repo.UserRegister(ctx, &Auth{Phone: account}, "Phone")
	case "email":
		user, err = receiver.repo.UserRegister(ctx, &Auth{Email: account}, "Email")
	case "wechat":
		user, err = receiver.repo.UserRegister(ctx, &Auth{Wechat: account}, "Wechat")
	case "github":
		user, err = receiver.repo.UserRegister(ctx, &Auth{Github: account}, "Github")
	default:
		err = ErrUserRegisterFailed
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func signToken(id int64, receiver *AuthUseCase) (*v1.LoginReply, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
	})
	signedString, err := claims.SignedString([]byte(receiver.key))
	if err != nil {
		return nil, v1.ErrorLoginFailed("login failed: generate token failed")
	}
	return &v1.LoginReply{
		Token: signedString,
	}, nil
}
