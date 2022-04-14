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

func (r *AuthUseCase) Login(ctx context.Context, account, password, code, mode string) (string, error) {
	if !loginVerify(r.validate, r.log, account, password, code, mode) {
		return "", v1.ErrorParamsIllegal("login failed: params illegal")
	}
	//if len(account) > 0 && len(password) > 0 && len(mode) > 0 {
	//	return loginByPassword(ctx, r, req.Account, req.Password, req.Mode)
	//}
	//if len(req.Account) > 0 && len(req.Code) > 0 && len(req.Mode) > 0 {
	//	return loginByCode(ctx, r, req.Account, req.Code, req.Mode)
	//}
	return "", nil
}

func (r *AuthUseCase) Register(ctx context.Context, account, mode string) (*User, error) {
	var user *User
	var err error
	if mode == "phone" || mode == "email" || mode == "wechat" || mode == "github" {
		user, err = r.repo.UserRegister(ctx, account, strings.ToUpper(mode[:1])+mode[1:])
	} else {
		err = ErrUserRegisterFailed
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *AuthUseCase) LoginPassWordForget(ctx context.Context, req *v1.LoginPassWordForgetReq) (*v1.LoginReply, error) {
	if len(req.Account) > 0 && len(req.Code) > 0 && len(req.Password) > 0 && (req.Mode == "phone" || req.Mode == "email") {
		return passWordForget(ctx, r, req.Account, req.Password, req.Code, req.Mode)
	}
	return nil, v1.ErrorParamsIllegal("login failed: params illegal")
}

func loginByPassword(ctx context.Context, r *AuthUseCase, account, password, mode string) (*v1.LoginReply, error) {
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

func loginByCode(ctx context.Context, r *AuthUseCase, account, code, mode string) (*v1.LoginReply, error) {
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

func passWordForget(ctx context.Context, r *AuthUseCase, account, password, code, mode string) (*v1.LoginReply, error) {
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

func signToken(id int64, r *AuthUseCase) (*v1.LoginReply, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
	})
	signedString, err := claims.SignedString([]byte(r.key))
	if err != nil {
		r.log.Errorf("fail to sign token: id(%v) error(%v)", id, err.Error())
		return nil, v1.ErrorUnknownError("login failed: generate token failed")
	}
	return &v1.LoginReply{
		Token: signedString,
	}, nil
}
