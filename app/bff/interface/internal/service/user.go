package service

import (
	"context"
	"github.com/the-zion/matrix-core/api/bff/interface/v1"
)

func (s *BffService) UserRegister(ctx context.Context, req *v1.UserRegisterReq) (*v1.UserRegisterReply, error) {
	err := s.uc.UserRegister(ctx, req.Email, req.Password, req.Code)
	if err != nil {
		return nil, err
	}
	return &v1.UserRegisterReply{}, nil
}

func (s *BffService) LoginByPassword(ctx context.Context, req *v1.LoginByPasswordReq) (*v1.LoginReply, error) {
	token, err := s.uc.LoginByPassword(ctx, req.Account, req.Password, req.Mode)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Token: token,
	}, nil
}

func (s *BffService) LoginByCode(ctx context.Context, req *v1.LoginByCodeReq) (*v1.LoginReply, error) {
	token, err := s.uc.LoginByCode(ctx, req.Phone, req.Code)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Token: token,
	}, nil
}

func (s *BffService) LoginPasswordReset(ctx context.Context, req *v1.LoginPasswordResetReq) (*v1.LoginPasswordResetReply, error) {
	err := s.uc.LoginPasswordReset(ctx, req.Account, req.Password, req.Code, req.Mode)
	if err != nil {
		return nil, err
	}
	return &v1.LoginPasswordResetReply{}, nil
}

func (s *BffService) SendPhoneCode(ctx context.Context, req *v1.SendPhoneCodeReq) (*v1.SendPhoneCodeReply, error) {
	err := s.uc.SendPhoneCode(ctx, req.Template, req.Phone)
	if err != nil {
		return nil, err
	}
	return &v1.SendPhoneCodeReply{}, nil
}

func (s *BffService) SendEmailCode(ctx context.Context, req *v1.SendEmailCodeReq) (*v1.SendEmailCodeReply, error) {
	err := s.uc.SendEmailCode(ctx, req.Template, req.Email)
	if err != nil {
		return nil, err
	}
	return &v1.SendEmailCodeReply{}, nil
}

func (s *BffService) GetUserProfile(ctx context.Context, req *v1.GetUserProfileReq) (*v1.GetUserProfileReply, error) {
	userProfile, err := s.uc.GetUserProfile(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserProfileReply{
		Uuid:      userProfile.Uuid,
		Username:  userProfile.Username,
		Avatar:    userProfile.Avatar,
		School:    userProfile.School,
		Homepage:  userProfile.Homepage,
		Introduce: userProfile.Introduce,
	}, nil
}
