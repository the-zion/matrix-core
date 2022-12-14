package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/user/service/v1"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserService) UserRegister(ctx context.Context, req *v1.UserRegisterReq) (*emptypb.Empty, error) {
	err := s.ac.UserRegister(ctx, req.Email, req.Password, req.Code)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) LoginByPassword(ctx context.Context, req *v1.LoginByPasswordReq) (*v1.LoginReply, error) {
	token, err := s.ac.LoginByPassword(ctx, req.Account, req.Password, req.Mode)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Token: token,
	}, nil
}

func (s *UserService) LoginByCode(ctx context.Context, req *v1.LoginByCodeReq) (*v1.LoginReply, error) {
	token, err := s.ac.LoginByCode(ctx, req.Phone, req.Code)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Token: token,
	}, nil
}

func (s *UserService) LoginByWeChat(ctx context.Context, req *v1.LoginByWeChatReq) (*v1.LoginReply, error) {
	token, err := s.ac.LoginByWechat(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Token: token,
	}, nil
}

func (s *UserService) LoginByQQ(ctx context.Context, req *v1.LoginByQQReq) (*v1.LoginReply, error) {
	token, err := s.ac.LoginByQQ(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Token: token,
	}, nil
}

func (s *UserService) LoginByGithub(ctx context.Context, req *v1.LoginByGithubReq) (*v1.LoginReply, error) {
	github, err := s.ac.LoginByGithub(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Token: github.Token,
	}, nil
}

func (s *UserService) LoginPasswordReset(ctx context.Context, req *v1.LoginPasswordResetReq) (*emptypb.Empty, error) {
	err := s.ac.LoginPasswordReset(ctx, req.Account, req.Password, req.Code, req.Mode)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) SetUserPhone(ctx context.Context, req *v1.SetUserPhoneReq) (*emptypb.Empty, error) {
	err := s.ac.SetUserPhone(ctx, req.Uuid, req.Phone, req.Code)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) SetUserEmail(ctx context.Context, req *v1.SetUserEmailReq) (*emptypb.Empty, error) {
	err := s.ac.SetUserEmail(ctx, req.Uuid, req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) SetUserPassword(ctx context.Context, req *v1.SetUserPasswordReq) (*emptypb.Empty, error) {
	err := s.ac.SetUserPassword(ctx, req.Uuid, req.Password)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) SendPhoneCode(ctx context.Context, req *v1.SendPhoneCodeReq) (*emptypb.Empty, error) {
	err := s.ac.SendPhoneCode(ctx, req.Template, req.Phone)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) SendEmailCode(ctx context.Context, req *v1.SendEmailCodeReq) (*emptypb.Empty, error) {
	err := s.ac.SendEmailCode(ctx, req.Template, req.Email)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) GetCosSessionKey(ctx context.Context, req *v1.GetCosSessionKeyReq) (*v1.GetCosSessionKeyReply, error) {
	credentials, err := s.ac.GetCosSessionKey(ctx, req.Uuid)
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

func (s *UserService) ChangeUserPassword(ctx context.Context, req *v1.ChangeUserPasswordReq) (*emptypb.Empty, error) {
	err := s.ac.ChangeUserPassword(ctx, req.Uuid, req.Oldpassword, req.Password)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) UnbindUserPhone(ctx context.Context, req *v1.UnbindUserPhoneReq) (*emptypb.Empty, error) {
	err := s.ac.UnbindUserPhone(ctx, req.Uuid, req.Phone, req.Code)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) UnbindUserEmail(ctx context.Context, req *v1.UnbindUserEmailReq) (*emptypb.Empty, error) {
	err := s.ac.UnbindUserEmail(ctx, req.Uuid, req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
