package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type MessageRepo interface {
	GetMailBoxLastTime(ctx context.Context, uuid string) (*MailBox, error)
	GetMessageNotification(ctx context.Context, uuid string, follows []string) (*Notification, error)
	SetMailBoxLastTime(ctx context.Context, uuid string, time int32) error
}

type MessageUseCase struct {
	repo     MessageRepo
	userRepo UserRepo
	re       Recovery
	log      *log.Helper
}

func NewMessageUseCase(repo MessageRepo, userRepo UserRepo, re Recovery, logger log.Logger) *MessageUseCase {
	return &MessageUseCase{
		repo:     repo,
		userRepo: userRepo,
		re:       re,
		log:      log.NewHelper(log.With(logger, "module", "bff/biz/MessageUseCase")),
	}
}

func (r *MessageUseCase) GetMessageNotification(ctx context.Context) (*Notification, error) {
	uuid := ctx.Value("uuid").(string)
	uuidMap, err := r.userRepo.GetUserFollows(ctx, uuid)
	if err != nil {
		return nil, err
	}
	follows := make([]string, 0)
	for key := range uuidMap {
		follows = append(follows, key)
	}
	return r.repo.GetMessageNotification(ctx, uuid, follows)
}

func (r *MessageUseCase) GetMailBoxLastTime(ctx context.Context) (*MailBox, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetMailBoxLastTime(ctx, uuid)
}

func (r *MessageUseCase) SetMailBoxLastTime(ctx context.Context, time int32) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.SetMailBoxLastTime(ctx, uuid, time)
}
