package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/user/service/v1"
	"github.com/the-zion/matrix-core/app/user/service/internal/biz"
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

func (s *UserService) GetUserProfile(ctx context.Context, req *v1.GetUserProfileReq) (*v1.GetUserProfileReply, error) {
	profile, err := s.uc.GetUserProfile(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserProfileReply{
		Uuid:      profile.Uuid,
		Username:  profile.Username,
		Avatar:    profile.Avatar,
		School:    profile.School,
		Company:   profile.Company,
		Job:       profile.Job,
		Homepage:  profile.Homepage,
		Introduce: profile.Introduce,
	}, nil
}

func (s *UserService) GetUserProfileUpdate(ctx context.Context, req *v1.GetUserProfileUpdateReq) (*v1.GetUserProfileUpdateReply, error) {
	profile, err := s.uc.GetUserProfileUpdate(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserProfileUpdateReply{
		Username:  profile.Username,
		Avatar:    profile.Avatar,
		School:    profile.School,
		Company:   profile.Company,
		Job:       profile.Job,
		Homepage:  profile.Homepage,
		Introduce: profile.Introduce,
		Status:    profile.Status,
	}, nil
}

func (s *UserService) SetUserProfile(ctx context.Context, req *v1.SetUserProfileReq) (*v1.SetUserProfileReply, error) {
	profile := &biz.ProfileUpdate{}
	profile.Uuid = req.Uuid
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
	return &v1.SetUserProfileReply{}, nil
}

func (s *UserService) SetUserPhone(ctx context.Context, req *v1.SetUserPhoneReq) (*v1.SetUserPhoneReply, error) {
	return nil, nil
}

func (s *UserService) SetUserEmail(ctx context.Context, req *v1.SetUserEmailReq) (*v1.SetUserEmailReply, error) {
	return nil, nil
}
