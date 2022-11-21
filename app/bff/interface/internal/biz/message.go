package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type MessageRepo interface {
	GetMailBoxLastTime(ctx context.Context, uuid string) (*MailBox, error)
	GetMessageNotification(ctx context.Context, uuid string, follows []string) (*Notification, error)
	GetMessageSystemNotification(ctx context.Context, page int32, uuid string) ([]*SystemNotification, error)
	SetMailBoxLastTime(ctx context.Context, uuid string, time int32) error
	RemoveMailBoxCommentCount(ctx context.Context, uuid string) error
	RemoveMailBoxSubCommentCount(ctx context.Context, uuid string) error
	RemoveMailBoxSystemNotificationCount(ctx context.Context, uuid string) error
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
	follows := make([]string, 0, len(uuidMap))
	for key := range uuidMap {
		follows = append(follows, key)
	}
	return r.repo.GetMessageNotification(ctx, uuid, follows)
}

func (r *MessageUseCase) GetMailBoxLastTime(ctx context.Context) (*MailBox, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetMailBoxLastTime(ctx, uuid)
}

func (r *MessageUseCase) GetMessageSystemNotification(ctx context.Context, page int32) ([]*SystemNotification, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetMessageSystemNotification(ctx, page, uuid)
}

func (r *MessageUseCase) SetMailBoxLastTime(ctx context.Context, time int32) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.SetMailBoxLastTime(ctx, uuid, time)
}

func (r *MessageUseCase) RemoveMailBoxCommentCount(ctx context.Context) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.RemoveMailBoxCommentCount(ctx, uuid)
}

func (r *MessageUseCase) RemoveMailBoxSubCommentCount(ctx context.Context) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.RemoveMailBoxSubCommentCount(ctx, uuid)
}

func (r *MessageUseCase) RemoveMailBoxSystemNotificationCount(ctx context.Context) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.RemoveMailBoxSystemNotificationCount(ctx, uuid)
}
