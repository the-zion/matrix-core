package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/user/service/v1"
	"github.com/the-zion/matrix-core/app/user/service/internal/biz"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserService) GetAccount(ctx context.Context, req *v1.GetAccountReq) (*v1.GetAccountReply, error) {
	user, err := s.uc.GetAccount(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetAccountReply{
		Phone:    user.Phone,
		Email:    user.Email,
		Qq:       user.Qq,
		Wechat:   user.Wechat,
		Weibo:    user.Weibo,
		Github:   user.Github,
		Password: user.Password,
	}, nil
}

func (s *UserService) GetProfile(ctx context.Context, req *v1.GetProfileReq) (*v1.GetProfileReply, error) {
	profile, err := s.uc.GetProfile(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetProfileReply{
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

func (s *UserService) GetProfileUpdate(ctx context.Context, req *v1.GetProfileUpdateReq) (*v1.GetProfileUpdateReply, error) {
	profile, err := s.uc.GetProfileUpdate(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetProfileUpdateReply{
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

func (s *UserService) SetProfileUpdate(ctx context.Context, req *v1.SetProfileUpdateReq) (*emptypb.Empty, error) {
	profile := &biz.ProfileUpdate{}
	profile.Uuid = req.Uuid
	profile.Username = req.Username
	profile.School = req.School
	profile.Company = req.Company
	profile.Job = req.Job
	profile.Homepage = req.Homepage
	profile.Introduce = req.Introduce
	err := s.uc.SetProfileUpdate(ctx, profile)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) ProfileReviewPass(ctx context.Context, req *v1.ProfileReviewPassReq) (*emptypb.Empty, error) {
	err := s.uc.ProfileReviewPass(ctx, req.Uuid, req.Update)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) ProfileReviewNotPass(ctx context.Context, req *v1.ProfileReviewNotPassReq) (*emptypb.Empty, error) {
	err := s.uc.ProfileReviewNotPass(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
