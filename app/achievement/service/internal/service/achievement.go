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
