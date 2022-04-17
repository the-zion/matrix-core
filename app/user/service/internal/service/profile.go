package service

import (
	"context"
	v1 "github.com/Cube-v2/cube-core/api/user/service/v1"
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
