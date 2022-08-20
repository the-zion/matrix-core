package service

import (
	"context"
	"github.com/the-zion/matrix-core/api/bff/interface/v1"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (s *BffService) GetUserMedal(ctx context.Context, req *v1.GetUserMedalReq) (*v1.GetUserMedalReply, error) {
	medal, err := s.achc.GetUserMedal(ctx, req.Uuid)
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

func (s *BffService) GetUserMedalProgress(ctx context.Context, _ *emptypb.Empty) (*v1.GetUserMedalProgressReply, error) {
	progress, err := s.achc.GetUserMedalProgress(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserMedalProgressReply{
		Creation:    progress.Article + progress.Talk,
		Agree:       progress.Agree,
		ActiveAgree: progress.ActiveAgree,
		View:        progress.View,
		Comment:     progress.Comment,
		Collect:     progress.Collect,
	}, nil
}

func (s *BffService) SetUserMedal(ctx context.Context, req *v1.SetUserMedalReq) (*emptypb.Empty, error) {
	err := s.achc.SetUserMedal(ctx, req.Medal)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) CancelUserMedalSet(ctx context.Context, req *v1.CancelUserMedalSetReq) (*emptypb.Empty, error) {
	err := s.achc.CancelUserMedalSet(ctx, req.Medal)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) AccessUserMedal(ctx context.Context, req *v1.AccessUserMedalReq) (*emptypb.Empty, error) {
	err := s.achc.AccessUserMedal(ctx, req.Medal)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
