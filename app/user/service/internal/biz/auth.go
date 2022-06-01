package biz

import (
	"context"
	"fmt"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	v1 "github.com/the-zion/matrix-core/api/user/service/v1"
	"github.com/the-zion/matrix-core/app/user/service/internal/conf"
)

type Login struct {
	Id    int64
	Token string
}

type AuthRepo interface {
	FindByAccount(ctx context.Context, account, mode string) (*User, error)
	RegisterWithPhone(ctx context.Context, phone string) (*User, error)
	SendPhoneCode(ctx context.Context, template, phone string) error
	SendEmailCode(ctx context.Context, template, phone string) error
	//SendCode(ctx context.Context, template int64, account, mode string) error
	//UserRegister(ctx context.Context, account, mode string) (*User, error)
	VerifyCode(ctx context.Context, account, code, mode string) error
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
		log:      log.NewHelper(log.With(logger, "module", "user/biz/authUseCase")),
	}
}

func (r *AuthUseCase) LoginByPassword(ctx context.Context, account, password, mode string) (*Login, error) {
	user, err := r.userRepo.FindByAccount(ctx, account, mode)
	if err != nil {
		return nil, v1.ErrorGetUserFailed("login failed: %s", err.Error())
	}

	err = r.userRepo.VerifyPassword(ctx, user.Id, password)
	if err != nil {
		return nil, v1.ErrorVerifyPasswordFailed("login failed: %s", err.Error())
	}

	token, err := signToken(user.Id, r.key)
	if err != nil {
		return nil, v1.ErrorUnknownError("login failed: %s", err.Error())
	}

	return &Login{
		Id:    user.Id,
		Token: token,
	}, nil
}

func (r *AuthUseCase) LoginByCode(ctx context.Context, phone, code string) (*Login, error) {
	err := r.repo.VerifyCode(ctx, phone, code, "email")
	if err != nil {
		return nil, v1.ErrorVerifyCodeFailed("login failed: %s", err.Error())
	}

	user, err := r.repo.FindByAccount(ctx, phone, code)
	if kerrors.IsNotFound(err) {
		user, err = r.repo.RegisterWithPhone(ctx, phone)
		if err != nil {
			return nil, v1.ErrorRegisterFailed("login failed: %s", err.Error())
		}
	}
	if err != nil {
		return nil, v1.ErrorUnknownError("login failed: %s", err.Error())
	}

	token, err := signToken(user.Id, r.key)
	if err != nil {
		return nil, v1.ErrorUnknownError("login failed: %s", err.Error())
	}

	return &Login{
		Id:    user.Id,
		Token: token,
	}, nil
}

//func (r *AuthUseCase) Register(ctx context.Context, account, mode string) (*User, error) {
//	user, err := r.repo.UserRegister(ctx, account, strings.ToUpper(mode[:1])+mode[1:])
//	if err != nil {
//		return nil, err
//	}
//	return user, nil
//}

func (r *AuthUseCase) LoginPasswordForget(ctx context.Context, account, password, code, mode string) (*Login, error) {
	err := r.userRepo.VerifyCode(ctx, account, code, mode)
	if err != nil {
		return nil, v1.ErrorVerifyCodeFailed("login failed: %s", err.Error())
	}

	user, err := r.userRepo.FindByAccount(ctx, account, mode)
	if err != nil {
		return nil, v1.ErrorGetUserFailed("login failed: %s", err.Error())
	}

	err = r.userRepo.PasswordModify(ctx, user.Id, password)
	if err != nil {
		return nil, v1.ErrorUnknownError("login failed: %s", err.Error())
	}

	token, err := signToken(user.Id, r.key)
	if err != nil {
		return nil, v1.ErrorUnknownError("login failed: %s", err.Error())
	}

	return &Login{
		Id:    user.Id,
		Token: token,
	}, nil
}

func (r *AuthUseCase) SendPhoneCode(ctx context.Context, template, phone string) error {
	err := r.repo.SendPhoneCode(ctx, template, phone)
	if err != nil {
		return v1.ErrorSendCodeFailed("send code failed: %s", err.Error())
	}
	return nil
}

func (r *AuthUseCase) SendEmailCode(ctx context.Context, template, email string) error {
	err := r.repo.SendEmailCode(ctx, template, email)
	if err != nil {
		return v1.ErrorSendCodeFailed("send code failed: %s", err.Error())
	}
	return nil
}

func signToken(id int64, key string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
	})
	signedString, err := claims.SignedString([]byte(key))
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to sign token: id(%v)", id))
	}
	return signedString, nil
}
