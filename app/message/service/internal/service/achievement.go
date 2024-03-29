package service

import (
	"context"
)

func (s *MessageService) SetAchievementAgree(ctx context.Context, uuid, userUuid string) error {
	return s.ac.SetAchievementAgree(ctx, uuid, userUuid)
}

func (s *MessageService) CancelAchievementAgree(ctx context.Context, uuid, userUuid string) error {
	return s.ac.CancelAchievementAgree(ctx, uuid, userUuid)
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

func (s *MessageService) SetAchievementFollow(ctx context.Context, follow, followed string) error {
	return s.ac.SetAchievementFollow(ctx, follow, followed)
}

func (s *MessageService) CancelAchievementFollow(ctx context.Context, follow, followed string) error {
	return s.ac.CancelAchievementFollow(ctx, follow, followed)
}

func (s *MessageService) AddAchievementScore(ctx context.Context, uuid string, score int32) error {
	return s.ac.AddAchievementScore(ctx, uuid, score)
}

func (s *MessageService) SetUserMedalDbAndCache(ctx context.Context, medal, uuid string) error {
	return s.ac.SetUserMedalDbAndCache(ctx, medal, uuid)
}

func (s *MessageService) CancelUserMedalDbAndCache(ctx context.Context, medal, uuid string) error {
	return s.ac.CancelUserMedalDbAndCache(ctx, medal, uuid)
}

func (s *MessageService) AccessUserMedalDbAndCache(ctx context.Context, medal, uuid string) error {
	return s.ac.AccessUserMedalDbAndCache(ctx, medal, uuid)
}
