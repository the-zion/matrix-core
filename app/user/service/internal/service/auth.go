package service

import (
	"context"
	v1 "github.com/Cube-v2/cube-core/api/user/service/v1"
)

func (s *UserService) LoginByPassword(ctx context.Context, req *v1.LoginByPasswordReq) (*v1.LoginReply, error) {
	login, err := s.ac.LoginByPassword(ctx, req.Account, req.Password, req.Mode)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Id:    login.Id,
		Token: login.Token,
	}, nil
}

func (s *UserService) LoginByEmail(ctx context.Context, req *v1.LoginByCodeReq) (*v1.LoginReply, error) {
	login, err := s.ac.LoginByCode(ctx, req.Account, req.Code, req.Mode)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Id:    login.Id,
		Token: login.Token,
	}, nil
}

func (s *UserService) LoginPassWordForget(ctx context.Context, req *v1.LoginPassWordForgetReq) (*v1.LoginReply, error) {
	login, err := s.ac.LoginPasswordForget(ctx, req.Account, req.Password, req.Code, req.Mode)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Id:    login.Id,
		Token: login.Token,
	}, nil
}
