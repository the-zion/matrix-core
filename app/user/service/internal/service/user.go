package service

import (
	"context"
	v1 "github.com/Cube-v2/cube-core/api/user/service/v1"
)

func (s *UserService) SendCode(ctx context.Context, req *v1.SendCodeReq) (*v1.SendCodeReply, error) {
	code, err := s.uc.SendCode(ctx, req.Template, req.Account, req.Mode)
	if err != nil {
		return nil, err
	}
	return &v1.SendCodeReply{
		Code: code,
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *v1.GetUserReq) (*v1.GetUserReply, error) {
	user, err := s.uc.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserReply{
		Phone:  user.Phone,
		Email:  user.Email,
		Wechat: user.Wechat,
		Github: user.Github,
	}, nil
}
