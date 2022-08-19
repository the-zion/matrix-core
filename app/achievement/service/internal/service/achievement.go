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

func (s *AchievementService) AddAchievementScore(ctx context.Context, req *v1.AddAchievementScoreReq) (*emptypb.Empty, error) {
	err := s.ac.AddAchievementScore(ctx, req.Uuid, req.Score)
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

func (s *AchievementService) GetUserMedal(ctx context.Context, req *v1.GetUserMedalReq) (*v1.GetUserMedalReply, error) {
	medal, err := s.ac.GetUserMedal(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserMedalReply{
		Creation1: medal.Creation1,
		Creation2: medal.Creation2,
		Creation3: medal.Creation3,
		Creation4: medal.Creation4,
		Creation5: medal.Creation5,
		Creation6: medal.Creation6,
		Creation7: medal.Creation7,
		Agree1:    medal.Agree1,
		Agree2:    medal.Agree2,
		Agree3:    medal.Agree3,
		Agree4:    medal.Agree4,
		Agree5:    medal.Agree5,
		Agree6:    medal.Agree6,
		View1:     medal.View1,
		View2:     medal.View2,
		View3:     medal.View3,
		Comment1:  medal.Comment1,
		Comment2:  medal.Comment2,
		Comment3:  medal.Comment3,
		Collect1:  medal.Collect1,
		Collect2:  medal.Collect2,
		Collect3:  medal.Collect3,
	}, nil
}

func (s *AchievementService) GetUserActive(ctx context.Context, req *v1.GetUserActiveReq) (*v1.GetUserActiveReply, error) {
	active, err := s.ac.GetUserActive(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserActiveReply{
		Agree: active.Agree,
	}, nil
}

func (s *AchievementService) SetUserMedal(ctx context.Context, req *v1.SetUserMedalReq) (*emptypb.Empty, error) {
	err := s.ac.SetUserMedal(ctx, req.Medal, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AchievementService) SetUserMedalDbAndCache(ctx context.Context, req *v1.SetUserMedalDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.ac.SetUserMedalDbAndCache(ctx, req.Medal, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AchievementService) CancelUserMedalSet(ctx context.Context, req *v1.CancelUserMedalSetReq) (*emptypb.Empty, error) {
	err := s.ac.CancelUserMedalSet(ctx, req.Medal, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AchievementService) CancelUserMedalDbAndCache(ctx context.Context, req *v1.CancelUserMedalDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.ac.CancelUserMedalDbAndCache(ctx, req.Medal, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AchievementService) AccessUserMedal(ctx context.Context, req *v1.AccessUserMedalReq) (*emptypb.Empty, error) {
	err := s.ac.AccessUserMedal(ctx, req.Medal, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AchievementService) AccessUserMedalDbAndCache(ctx context.Context, req *v1.AccessUserMedalReq) (*emptypb.Empty, error) {
	err := s.ac.AccessUserMedalDbAndCache(ctx, req.Medal, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
