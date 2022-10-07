package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/message/service/v1"
)

type MessageRepo interface {
	GetMailBoxLastTime(ctx context.Context, uuid string) (*MailBox, error)
	SetMailBoxLastTime(ctx context.Context, uuid string) error
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

func (r *MessageUseCase) SetMailBoxLastTime(ctx context.Context, uuid string) error {
	err := r.repo.SetMailBoxLastTime(ctx, uuid)
	if err != nil {
		return v1.ErrorSetMailboxLastTimeFailed("set mailbox last time failed: %s", err.Error())
	}
	return nil
}
