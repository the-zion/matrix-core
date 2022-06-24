package service

import (
	"context"
	"github.com/the-zion/matrix-core/api/bff/interface/v1"
	"github.com/the-zion/matrix-core/app/bff/interface/internal/biz"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *BffService) UserRegister(ctx context.Context, req *v1.UserRegisterReq) (*emptypb.Empty, error) {
	err := s.uc.UserRegister(ctx, req.Email, req.Password, req.Code)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
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

func (s *BffService) LoginPasswordReset(ctx context.Context, req *v1.LoginPasswordResetReq) (*emptypb.Empty, error) {
	err := s.uc.LoginPasswordReset(ctx, req.Account, req.Password, req.Code, req.Mode)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SendPhoneCode(ctx context.Context, req *v1.SendPhoneCodeReq) (*emptypb.Empty, error) {
	err := s.uc.SendPhoneCode(ctx, req.Template, req.Phone)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SendEmailCode(ctx context.Context, req *v1.SendEmailCodeReq) (*emptypb.Empty, error) {
	err := s.uc.SendEmailCode(ctx, req.Template, req.Email)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) GetCosSessionKey(ctx context.Context, _ *emptypb.Empty) (*v1.GetCosSessionKeyReply, error) {
	credentials, err := s.uc.GetCosSessionKey(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetCosSessionKeyReply{
		TmpSecretId:  credentials.TmpSecretID,
		TmpSecretKey: credentials.TmpSecretKey,
		SessionToken: credentials.SessionToken,
		StartTime:    credentials.StartTime,
		ExpiredTime:  credentials.ExpiredTime,
	}, nil
}

func (s *BffService) GetProfile(ctx context.Context, _ *emptypb.Empty) (*v1.GetProfileReply, error) {
	userProfile, err := s.uc.GetProfile(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetProfileReply{
		Uuid:      userProfile.Uuid,
		Username:  userProfile.Username,
		Avatar:    userProfile.Avatar,
		School:    userProfile.School,
		Company:   userProfile.Company,
		Job:       userProfile.Job,
		Homepage:  userProfile.Homepage,
		Introduce: userProfile.Introduce,
	}, nil
}

func (s *BffService) GetProfileUpdate(ctx context.Context, _ *emptypb.Empty) (*v1.GetProfileUpdateReply, error) {
	userProfile, err := s.uc.GetProfileUpdate(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetProfileUpdateReply{
		Username:  userProfile.Username,
		Avatar:    userProfile.Avatar,
		School:    userProfile.School,
		Company:   userProfile.Company,
		Job:       userProfile.Job,
		Homepage:  userProfile.Homepage,
		Introduce: userProfile.Introduce,
		Status:    userProfile.Status,
	}, nil
}

func (s *BffService) SetProfileUpdate(ctx context.Context, req *v1.SetProfileUpdateReq) (*emptypb.Empty, error) {
	profile := &biz.UserProfileUpdate{}
	profile.Username = req.Username
	profile.School = req.School
	profile.Company = req.Company
	profile.Job = req.Job
	profile.Homepage = req.Homepage
	profile.Introduce = req.Introduce
	err := s.uc.SetUserProfile(ctx, profile)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
