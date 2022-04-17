package biz

import (
	"context"
	"errors"
	v1 "github.com/Cube-v2/cube-core/api/user/service/v1"
	"github.com/Cube-v2/cube-core/app/user/service/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-playground/validator/v10"
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
	validate *validator.Validate
	log      *log.Helper
}

func NewAuthUseCase(conf *conf.Auth, repo AuthRepo, userRepo UserRepo, validator *validator.Validate, logger log.Logger) *AuthUseCase {
	return &AuthUseCase{
		key:      conf.ApiKey,
		repo:     repo,
		userRepo: userRepo,
		validate: validator,
		log:      log.NewHelper(log.With(logger, "module", "user/biz/authUseCase")),
	}
}

//func (r *AuthUseCase) Login(ctx context.Context, account, password, code, mode string) (*Login, error) {
//	//if !loginVerify(r.validate, r.log, account, password, code, mode) {
//	//	return nil, v1.ErrorParamsIllegal("login failed: params illegal")
//	//}
//	if len(password) > 0 {
//		return loginByPassword(ctx, r, account, password, mode)
//	} else {
//		return loginByCode(ctx, r, account, code, mode)
//	}
//}

//func (r *AuthUseCase) LoginByPhone(ctx context.Context, account, password, code string) (*Login, error) {
//	if len(password) > 0 {
//		return loginByPassword(ctx, r, account, password, "phone")
//	} else {
//		return loginByCode(ctx, r, account, code, "phone")
//	}
//}
//
//func (r *AuthUseCase) LoginByEmail(ctx context.Context, account, password, code string) (*Login, error) {
//	if len(password) > 0 {
//		return loginByPassword(ctx, r, account, password, "email")
//	} else {
//		return loginByCode(ctx, r, account, code, "email")
//	}
//}

func (r *AuthUseCase) LoginByPassword(ctx context.Context, account, password, mode string) (*Login, error) {
	user, err := loginFindByAccount(ctx, r, account, mode)
	if err != nil {
		return nil, err
	}

	err = r.userRepo.VerifyPassword(ctx, user.Id, password)
	if errors.Is(err, ErrUnknownError) {
		return nil, v1.ErrorUnknownError("login failed: %s", err.Error())
	}
	if errors.Is(err, ErrPasswordError) {
		return nil, v1.ErrorVerifyPasswordFailed("login failed: %s", err.Error())
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
	if !loginPasswordVerify(r.validate, r.log, account, password, code, mode) {
		return nil, v1.ErrorParamsIllegal("login failed: params illegal")
	}
	return passWordForget(ctx, r, account, password, code, mode)
}

func loginFindByAccount(ctx context.Context, r *AuthUseCase, account, mode string) (*User, error) {
	user, err := r.userRepo.FindByAccount(ctx, account, mode)
	if err != nil {
		return nil, v1.ErrorGetUserFailed("login failed: %s", err.Error())
	}
	return user, nil
}

func loginVerifyCode(ctx context.Context, r *AuthUseCase, account, code, mode string) error {
	err := r.userRepo.VerifyCode(ctx, account, code, mode)
	if errors.Is(err, ErrUnknownError) {
		return v1.ErrorUnknownError("login failed: %s", err.Error())
	}
	if errors.Is(err, ErrCodeError) {
		return v1.ErrorVerifyCodeFailed("login failed: %s", err.Error())
	}
	return nil
}

func passWordForget(ctx context.Context, r *AuthUseCase, account, password, code, mode string) (*Login, error) {
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
