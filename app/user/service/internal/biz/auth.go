package biz

import (
	"context"
	v1 "github.com/Cube-v2/cube-core/api/user/service/v1"
	"github.com/Cube-v2/cube-core/app/user/service/internal/conf"
	"github.com/golang-jwt/jwt/v4"
)

type AuthUseCase struct {
	key      string
	userRepo UserRepo
}

func NewLoginUseCase(conf *conf.Auth, userRepo UserRepo) *AuthUseCase {
	return &AuthUseCase{
		key:      conf.ApiKey,
		userRepo: userRepo,
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
	return &v1.LoginReply{}, v1.ErrorLoginFailed("login: params illegal")
}

func loginByPhoneAndPassword(ctx context.Context, receiver *AuthUseCase, phone, password string) (*v1.LoginReply, error) {
	user, err := receiver.userRepo.FindByUserPhone(ctx, phone)
	if err != nil {
		return nil, v1.ErrorLoginFailed("account not found: %s", err.Error())
	}

	err = receiver.userRepo.VerifyPassword(ctx, user, password)
	if err != nil {
		return nil, v1.ErrorLoginFailed("password not match")
	}

	return signToken(user.Id, receiver)
}

func loginByPhoneAndCode(ctx context.Context, receiver *AuthUseCase, phone, code string) (*v1.LoginReply, error) {
	user, err := receiver.userRepo.FindByUserPhone(ctx, phone)
	if err != nil {
		return nil, v1.ErrorLoginFailed("login failed: %s", err.Error())
	}

	err = receiver.userRepo.VerifyCode(ctx, user, code)
	if err != nil {
		return nil, v1.ErrorLoginFailed("code not match")
	}

	return signToken(user.Id, receiver)
}

func loginByEmailAndPassword(ctx context.Context, receiver *AuthUseCase, email, password string) (*v1.LoginReply, error) {
	user, err := receiver.userRepo.FindByUserEmail(ctx, email)
	if err != nil {
		return nil, v1.ErrorLoginFailed("account not found: %s", err.Error())
	}

	err = receiver.userRepo.VerifyPassword(ctx, user, password)
	if err != nil {
		return nil, v1.ErrorLoginFailed("code not match")
	}

	return signToken(user.Id, receiver)
}

func loginByEmailAndCode(ctx context.Context, receiver *AuthUseCase, email, code string) (*v1.LoginReply, error) {
	user, err := receiver.userRepo.FindByUserEmail(ctx, email)
	if err != nil {
		return nil, v1.ErrorLoginFailed("login failed: %s", err.Error())
	}

	err = receiver.userRepo.VerifyCode(ctx, user, code)
	if err != nil {
		return nil, v1.ErrorLoginFailed("code not match")
	}

	return signToken(user.Id, receiver)
}

func (receiver *AuthUseCase) Register(ctx context.Context, account, mode string) error {
	var err error
	switch mode {
	case "phone":
		err = receiver.userRepo.UserRegister(ctx, &User{Phone: account})
	case "email":
		err = receiver.userRepo.UserRegister(ctx, &User{Email: account})
	case "wechat":
		err = receiver.userRepo.UserRegister(ctx, &User{Wechat: account})
	case "github":
		err = receiver.userRepo.UserRegister(ctx, &User{Github: account})
	}
	if err != nil {
		return v1.ErrorRegisterFailed("account created failed: %s", err.Error())
	}
	return nil
}

func signToken(id int64, receiver *AuthUseCase) (*v1.LoginReply, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
	})
	signedString, err := claims.SignedString([]byte(receiver.key))
	if err != nil {
		return nil, v1.ErrorLoginFailed("generate token failed: %s", err.Error())
	}
	return &v1.LoginReply{
		Token: signedString,
	}, nil
}
