package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/achievement/service/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AchievementService) SetAchievementAgree(ctx context.Context, req *v1.SetAchievementAgreeReq) (*emptypb.Empty, error) {
	err := s.ac.SetAchievementAgree(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AchievementService) CancelAchievementAgree(ctx context.Context, req *v1.CancelAchievementAgreeReq) (*emptypb.Empty, error) {
	err := s.ac.CancelAchievementAgree(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AchievementService) SetAchievementView(ctx context.Context, req *v1.SetAchievementViewReq) (*emptypb.Empty, error) {
	err := s.ac.SetAchievementView(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AchievementService) SetAchievementCollect(ctx context.Context, req *v1.SetAchievementCollectReq) (*emptypb.Empty, error) {
	err := s.ac.SetAchievementCollect(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AchievementService) CancelAchievementCollect(ctx context.Context, req *v1.CancelAchievementCollectReq) (*emptypb.Empty, error) {
	err := s.ac.CancelAchievementCollect(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AchievementService) SetAchievementFollow(ctx context.Context, req *v1.SetAchievementFollowReq) (*emptypb.Empty, error) {
	err := s.ac.SetAchievementFollow(ctx, req.Follow, req.Followed)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AchievementService) CancelAchievementFollow(ctx context.Context, req *v1.CancelAchievementFollowReq) (*emptypb.Empty, error) {
	err := s.ac.CancelAchievementFollow(ctx, req.Follow, req.Followed)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AchievementService) GetAchievementList(ctx context.Context, req *v1.GetAchievementListReq) (*v1.GetAchievementListReply, error) {
	reply := &v1.GetAchievementListReply{Achievement: make([]*v1.GetAchievementListReply_Achievement, 0)}
	achievementList, err := s.ac.GetAchievementList(ctx, req.Uuids)
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

func (s *AchievementService) GetUserAchievement(ctx context.Context, req *v1.GetUserAchievementReq) (*v1.GetUserAchievementReply, error) {
	achievement, err := s.ac.GetUserAchievement(ctx, req.Uuid)
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
