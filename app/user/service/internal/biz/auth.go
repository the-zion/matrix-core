package biz

import (
	"context"
	"fmt"
	v1 "github.com/Cube-v2/cube-core/api/user/service/v1"
	"github.com/Cube-v2/cube-core/app/user/service/internal/conf"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"strings"
)

type Login struct {
	Id    int64
	Token string
}

type AuthRepo interface {
	UserRegister(ctx context.Context, account, mode string) (*User, error)
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

func (r *AuthUseCase) LoginByCode(ctx context.Context, account, code, mode string) (*Login, error) {
	err := r.userRepo.VerifyCode(ctx, account, code, mode)
	if err != nil {
		return nil, v1.ErrorVerifyCodeFailed("login failed: %s", err.Error())
	}

	user, err := r.userRepo.FindByAccount(ctx, account, mode)
	if kerrors.IsNotFound(err) {
		user, err = r.Register(ctx, account, mode)
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

func (r *AuthUseCase) Register(ctx context.Context, account, mode string) (*User, error) {
	user, err := r.repo.UserRegister(ctx, account, strings.ToUpper(mode[:1])+mode[1:])
	if err != nil {
		return nil, err
	}
	return user, nil
}

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
