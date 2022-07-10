package service

import (
	"context"
)

func (s *MessageService) SetAchievementAgree(ctx context.Context, uuid string) error {
	return s.ac.SetAchievementAgree(ctx, uuid)
}

func (s *MessageService) CancelAchievementAgree(ctx context.Context, uuid string) error {
	return s.ac.CancelAchievementAgree(ctx, uuid)
}

func (s *MessageService) SetAchievementView(ctx context.Context, uuid string) error {
	return s.ac.SetAchievementView(ctx, uuid)
}

func (s *MessageService) SetAchievementCollect(ctx context.Context, uuid string) error {
	return s.ac.SetAchievementCollect(ctx, uuid)
}

func (s *MessageService) CancelAchievementCollect(ctx context.Context, uuid string) error {
	return s.ac.CancelAchievementCollect(ctx, uuid)
}
