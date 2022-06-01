package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/user/service/v1"
)

func (s *UserService) GetUserProfile(ctx context.Context, req *v1.GetUserProfileReq) (*v1.GetUserProfileReply, error) {
	profile, err := s.pc.GetUserProfile(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserProfileReply{
		Username:   profile.Username,
		Sex:        profile.Sex,
		Introduce:  profile.Introduce,
		Address:    profile.Address,
		Industry:   profile.Industry,
		Profile:    profile.PersonalProfile,
		Tag:        profile.Tag,
		Background: profile.Background,
		Image:      profile.Image,
	}, nil
}

func (s *UserService) SetUserProfile(ctx context.Context, req *v1.SetUserProfileReq) (*v1.SetUserProfileReply, error) {
	id := GetUserIdFromContext(ctx)
	err := s.pc.SetUserProfile(ctx, id, req.Sex, req.Introduce, req.Industry, req.Address, req.Profile, req.Tag)
	if err != nil {
		return nil, err
	}
	return &v1.SetUserProfileReply{
		Success: true,
	}, nil
}

func (s *UserService) SetUserName(ctx context.Context, req *v1.SetUserNameReq) (*v1.SetUserNameReply, error) {
	id := GetUserIdFromContext(ctx)
	err := s.pc.SetUserName(ctx, id, req.Name)
	if err != nil {
		return nil, err
	}
	return &v1.SetUserNameReply{
		Success: true,
	}, nil
}
