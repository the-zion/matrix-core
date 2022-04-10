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
	ErrUserNotFound       = errors.New("user not found")
	ErrPasswordError      = errors.New("password error")
	ErrCodeError          = errors.New("code error")
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
	FindByUserPhone(ctx context.Context, phone string) (*Auth, error)
	FindByUserEmail(ctx context.Context, email string) (*Auth, error)
	VerifyPassword(ctx context.Context, id int64, password string) error
	VerifyCode(ctx context.Context, key, code string) error
	UserRegister(ctx context.Context, u *Auth, mode string) (*Auth, error)
}

type AuthUseCase struct {
	key  string
	repo AuthRepo
	log  *log.Helper
}

func NewAuthUseCase(conf *conf.Auth, repo AuthRepo, logger log.Logger) *AuthUseCase {
	return &AuthUseCase{
		key:  conf.ApiKey,
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "usecase/auth")),
	}
}

func (receiver *AuthUseCase) Login(ctx context.Context, req *v1.LoginReq) (*v1.LoginReply, error) {
	if len(req.Phone) > 0 && len(req.Password) > 0 {
		return loginByPhoneAndPassword(ctx, receiver, req.Phone, req.Password)
	}
	if len(req.Phone) > 0 && len(req.Code) > 0 {
		return loginByPhoneAndCode(ctx, receiver, req.Phone, req.Code)
	}
	if len(req.Email) > 0 && len(req.Password) > 0 {
		return loginByEmailAndPassword(ctx, receiver, req.Email, req.Password)
	}
	if len(req.Email) > 0 && len(req.Code) > 0 {
		return loginByEmailAndCode(ctx, receiver, req.Email, req.Code)
	}
	return &v1.LoginReply{}, v1.ErrorLoginFailed("login failed: params illegal")
}

func loginByPhoneAndPassword(ctx context.Context, receiver *AuthUseCase, phone, password string) (*v1.LoginReply, error) {
	user, err := receiver.repo.FindByUserPhone(ctx, phone)
	if err != nil {
		return nil, v1.ErrorLoginFailed("login failed: %s", err.Error())
	}

	err = receiver.repo.VerifyPassword(ctx, user.Id, password)
	if err != nil {
		return nil, v1.ErrorLoginFailed("login failed: %s", err.Error())
	}

	return signToken(user.Id, receiver)
}

func loginByPhoneAndCode(ctx context.Context, receiver *AuthUseCase, phone, code string) (*v1.LoginReply, error) {
	err := receiver.repo.VerifyCode(ctx, phone, code)
	if err != nil {
		return nil, v1.ErrorLoginFailed("login failed: %s", err.Error())
	}

	user, err := receiver.repo.FindByUserPhone(ctx, phone)
	if err != nil {
		user, err = receiver.Register(ctx, phone, "phone")
		if err != nil {
			return nil, v1.ErrorLoginFailed("login failed: %s", err.Error())
		}
	}

	return signToken(user.Id, receiver)
}

func loginByEmailAndPassword(ctx context.Context, receiver *AuthUseCase, email, password string) (*v1.LoginReply, error) {
	user, err := receiver.repo.FindByUserEmail(ctx, email)
	if err != nil {
		return nil, v1.ErrorLoginFailed("account not found: %s", err.Error())
	}

	err = receiver.repo.VerifyPassword(ctx, user.Id, password)
	if err != nil {
		return nil, v1.ErrorLoginFailed("login failed: %s", err.Error())
	}

	return signToken(user.Id, receiver)
}

func loginByEmailAndCode(ctx context.Context, receiver *AuthUseCase, email, code string) (*v1.LoginReply, error) {
	err := receiver.repo.VerifyCode(ctx, email, code)
	if err != nil {
		return nil, v1.ErrorLoginFailed("login failed: %s", err.Error())
	}

	user, err := receiver.repo.FindByUserEmail(ctx, email)
	if err != nil {
		user, err = receiver.Register(ctx, email, "email")
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
