package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/user/service/v1"
)

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

func (s *UserService) SetUserPhone(ctx context.Context, req *v1.SetUserPhoneReq) (*v1.SetUserPhoneReply, error) {
	id := GetUserIdFromContext(ctx)
	err := s.uc.SetUserPhone(ctx, id, req.Phone, req.Password, req.Code)
	if err != nil {
		return nil, err
	}
	return &v1.SetUserPhoneReply{
		Success: true,
	}, nil
}

func (s *UserService) SetUserEmail(ctx context.Context, req *v1.SetUserEmailReq) (*v1.SetUserEmailReply, error) {
	id := GetUserIdFromContext(ctx)
	err := s.uc.SetUserEmail(ctx, id, req.Email, req.Password, req.Code)
	if err != nil {
		return nil, err
	}
	return &v1.SetUserEmailReply{
		Success: true,
	}, nil
}
