package service

import (
	"context"
	v1 "github.com/Cube-v2/cube-core/api/user/service/v1"
)

func (s *UserService) GetUserAchievement(ctx context.Context, req *v1.GetUserAchievementReq) (*v1.GetUserAchievementReply, error) {
	ach, err := s.achc.GetUserAchievement(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserAchievementReply{
		Follow:   ach.Follow,
		Followed: ach.Followed,
		Agree:    ach.Agree,
		Collect:  ach.Collect,
		View:     ach.View,
	}, nil
}
