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
	var reply *v1.LoginReply
	var err error
	if req.Phone != "" && req.Password != "" {
		reply, err = receiver.LoginByPhoneAndPassword(ctx, req)
	}
	if req.Phone != "" && req.Code != "" {
		reply, err = receiver.LoginByPhoneAndCode(ctx, req)
	}
	if req.Email != "" && req.Password != "" {
		reply, err = receiver.LoginByEmailAndPassword(ctx, req)
	}
	if req.Email != "" && req.Code != "" {
		reply, err = receiver.LoginByEmailAndCode(ctx, req)
	}
	return reply, err
}

func (receiver *AuthUseCase) LoginByPhoneAndPassword(ctx context.Context, req *v1.LoginReq) (*v1.LoginReply, error) {
	user, err := receiver.userRepo.FindByUserPhone(ctx, req.Phone)
	if err != nil {
		return nil, v1.ErrorLoginFailed("user not found: %s", err.Error())
	}

	err = receiver.userRepo.VerifyPassword(ctx, user, req.Password)
	if err != nil {
		return nil, v1.ErrorLoginFailed("password not match")
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
	})
	signedString, err := claims.SignedString([]byte(receiver.key))
	if err != nil {
		return nil, v1.ErrorLoginFailed("generate token failed: %s", err.Error())
	}
	return &v1.LoginReply{
		Token: signedString,
	}, nil
}

func (receiver *AuthUseCase) LoginByPhoneAndCode(ctx context.Context, req *v1.LoginReq) (*v1.LoginReply, error) {
	return nil, nil
}

func (receiver *AuthUseCase) LoginByEmailAndPassword(ctx context.Context, req *v1.LoginReq) (*v1.LoginReply, error) {
	return nil, nil
}

func (receiver *AuthUseCase) LoginByEmailAndCode(ctx context.Context, req *v1.LoginReq) (*v1.LoginReply, error) {
	return nil, nil
}
