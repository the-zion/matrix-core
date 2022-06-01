package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/user/service/v1"
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

func (s *UserService) LoginByCode(ctx context.Context, req *v1.LoginByCodeReq) (*v1.LoginReply, error) {
	login, err := s.ac.LoginByCode(ctx, req.Phone, req.Code)
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

func (s *UserService) SendPhoneCode(ctx context.Context, req *v1.SendPhoneCodeReq) (*v1.SendPhoneCodeReply, error) {
	err := s.ac.SendPhoneCode(ctx, req.Template, req.Phone)
	if err != nil {
		return nil, err
	}
	return &v1.SendPhoneCodeReply{
		Success: true,
	}, nil
}

func (s *UserService) SendEmailCode(ctx context.Context, req *v1.SendEmailCodeReq) (*v1.SendEmailCodeReply, error) {
	err := s.ac.SendEmailCode(ctx, req.Template, req.Email)
	if err != nil {
		return nil, err
	}
	return &v1.SendEmailCodeReply{
		Success: true,
	}, nil
}
