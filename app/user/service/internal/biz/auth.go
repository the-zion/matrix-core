package biz

import (
	"context"
	"errors"
	v1 "github.com/Cube-v2/cube-core/api/user/service/v1"
	"github.com/Cube-v2/cube-core/app/user/service/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"strings"
)

var (
	ErrUserRegisterFailed = errors.New("user register failed")
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
	user, err := loginFindByAccount(ctx, r, account, mode)
	if err != nil {
		return nil, err
	}

	errFormat := "login failed: %s"
	err = VerifyPassword(ctx, r.userRepo, user.Id, password, errFormat)
	if err != nil {
		return nil, err
	}

	return signToken(user.Id, r)
}

func (r *AuthUseCase) LoginByCode(ctx context.Context, account, code, mode string) (*Login, error) {
	err := loginVerifyCode(ctx, r, account, code, mode)
	if err != nil {
		return nil, err
	}

	user, err := loginFindByAccount(ctx, r, account, mode)
	if err != nil {
		user, err = r.Register(ctx, account, mode)
		if err != nil {
			return nil, v1.ErrorUnknownError("login failed: %s", err.Error())
		}
	}

	return signToken(user.Id, r)
}

func (r *AuthUseCase) Register(ctx context.Context, account, mode string) (*User, error) {
	user, err := r.repo.UserRegister(ctx, account, strings.ToUpper(mode[:1])+mode[1:])
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *AuthUseCase) LoginPasswordForget(ctx context.Context, account, password, code, mode string) (*Login, error) {
	err := loginVerifyCode(ctx, r, account, code, mode)
	if err != nil {
		return nil, err
	}
	user, err := loginFindByAccount(ctx, r, account, mode)
	if err != nil {
		return nil, err
	}
	err = r.userRepo.PasswordModify(ctx, user.Id, password)
	if err != nil {
		return nil, v1.ErrorUnknownError("login failed: %s", err.Error())
	}
	return signToken(user.Id, r)
}

func loginFindByAccount(ctx context.Context, r *AuthUseCase, account, mode string) (*User, error) {
	user, err := r.userRepo.FindByAccount(ctx, account, mode)
	if err != nil {
		return nil, v1.ErrorGetUserFailed("login failed: %s", err.Error())
	}
	return user, nil
}

func loginVerifyCode(ctx context.Context, r *AuthUseCase, account, code, mode string) error {
	return VerifyCode(ctx, r.userRepo, account, code, mode, "login failed: %s")
}

func signToken(id int64, r *AuthUseCase) (*Login, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
	})
	signedString, err := claims.SignedString([]byte(r.key))
	if err != nil {
		r.log.Errorf("fail to sign token: id(%v) error(%v)", id, err.Error())
		return nil, v1.ErrorUnknownError("login failed: generate token failed")
	}
	return &Login{
		Id:    id,
		Token: signedString,
	}, nil
}
