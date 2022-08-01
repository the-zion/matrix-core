package service

import (
	"context"
	"github.com/the-zion/matrix-core/api/bff/interface/v1"
)

func (s *BffService) GetAchievementList(ctx context.Context, req *v1.GetAchievementListReq) (*v1.GetAchievementListReply, error) {
	reply := &v1.GetAchievementListReply{Achievement: make([]*v1.GetAchievementListReply_Achievement, 0)}
	achievementList, err := s.achc.GetAchievementList(ctx, req.Uuids)
	if err != nil {
		return nil, err
	}
	for _, item := range achievementList {
		reply.Achievement = append(reply.Achievement, &v1.GetAchievementListReply_Achievement{
			Uuid:     item.Uuid,
			View:     item.View,
			Agree:    item.Agree,
			Follow:   item.Follow,
			Followed: item.Followed,
		})
	}
	return reply, nil
}

func (s *BffService) GetUserAchievement(ctx context.Context, req *v1.GetUserAchievementReq) (*v1.GetUserAchievementReply, error) {
	achievement, err := s.achc.GetUserAchievement(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserAchievementReply{
		Agree:    achievement.Agree,
		View:     achievement.View,
		Collect:  achievement.Collect,
		Follow:   achievement.Follow,
		Followed: achievement.Followed,
	}, nil
}
