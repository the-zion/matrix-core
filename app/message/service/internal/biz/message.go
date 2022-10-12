package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/message/service/v1"
)

type MessageRepo interface {
	GetMailBoxLastTime(ctx context.Context, uuid string) (*MailBox, error)
	GetMessageNotification(ctx context.Context, uuid string, follows []string) (*Notification, error)
	GetMessageSystemNotification(ctx context.Context, page int32, uuid string) ([]*SystemNotification, error)
	SetMailBoxLastTime(ctx context.Context, uuid string, time int32) error
	RemoveMailBoxCommentCount(ctx context.Context, uuid string) error
	RemoveMailBoxSubCommentCount(ctx context.Context, uuid string) error
	RemoveMailBoxSystemNotificationCount(ctx context.Context, uuid string) error
	AddMailBoxSystemNotification(ctx context.Context, contentId int32, notificationType string, title string, uuid string, label string, result int32, section string, text string, uid string, comment string) (*SystemNotification, error)
	AddMailBoxSystemNotificationToCache(ctx context.Context, notification *SystemNotification) error
}

type MessageUseCase struct {
	repo MessageRepo
	re   Recovery
	log  *log.Helper
}

func NewMessageUseCase(repo MessageRepo, re Recovery, logger log.Logger) *MessageUseCase {
	return &MessageUseCase{
		repo: repo,
		re:   re,
		log:  log.NewHelper(log.With(logger, "module", "message/biz/messageUseCase")),
	}
}

func (r *MessageUseCase) GetMailBoxLastTime(ctx context.Context, uuid string) (*MailBox, error) {
	mailbox, err := r.repo.GetMailBoxLastTime(ctx, uuid)
	if err != nil {
		return nil, v1.ErrorGetMailboxLastTimeFailed("get mailbox last time failed: %s", err.Error())
	}
	return mailbox, nil
}

func (r *MessageUseCase) GetMessageNotification(ctx context.Context, uuid string, follows []string) (*Notification, error) {
	mailbox, err := r.repo.GetMessageNotification(ctx, uuid, follows)
	if err != nil {
		return nil, v1.ErrorGetMessageNotificationFailed("get message notification failed: %s", err.Error())
	}
	return mailbox, nil
}

func (r *MessageUseCase) GetMessageSystemNotification(ctx context.Context, page int32, uuid string) ([]*SystemNotification, error) {
	notificationList, err := r.repo.GetMessageSystemNotification(ctx, page, uuid)
	if err != nil {
		return nil, v1.ErrorGetMessageNotificationFailed("get message system notification failed: %s", err.Error())
	}
	return notificationList, nil
}

func (r *MessageUseCase) SetMailBoxLastTime(ctx context.Context, uuid string, time int32) error {
	err := r.repo.SetMailBoxLastTime(ctx, uuid, time)
	if err != nil {
		return v1.ErrorSetMailboxLastTimeFailed("set mailbox last time failed: %s", err.Error())
	}
	return nil
}

func (r *MessageUseCase) RemoveMailBoxCommentCount(ctx context.Context, uuid string) error {
	err := r.repo.RemoveMailBoxCommentCount(ctx, uuid)
	if err != nil {
		return v1.ErrorRemoveMailboxCommentFailed("remove mail box comment count: %s", err.Error())
	}
	return nil
}

func (r *MessageUseCase) RemoveMailBoxSubCommentCount(ctx context.Context, uuid string) error {
	err := r.repo.RemoveMailBoxSubCommentCount(ctx, uuid)
	if err != nil {
		return v1.ErrorRemoveMailboxCommentFailed("remove mail box sub comment count: %s", err.Error())
	}
	return nil
}

func (r *MessageUseCase) RemoveMailBoxSystemNotificationCount(ctx context.Context, uuid string) error {
	err := r.repo.RemoveMailBoxSystemNotificationCount(ctx, uuid)
	if err != nil {
		return v1.ErrorRemoveMailboxSystemNotificationFailed("remove mail box system notification failed: %s", err.Error())
	}
	return nil
}
