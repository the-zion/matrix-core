package service

import (
	"context"
	v1 "github.com/Cube-v2/cube-core/api/user/service/v1"
)

func (s *UserService) Login(ctx context.Context, req *v1.LoginReq) (*v1.LoginReply, error) {
	token, err := s.ac.Login(ctx, req.Account, req.Password, req.Code, req.Mode)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Token: token,
	}, nil
}

func (s *UserService) LoginPassWordForget(ctx context.Context, req *v1.LoginPassWordForgetReq) (*v1.LoginReply, error) {
	return s.ac.LoginPassWordForget(ctx, req)
}
