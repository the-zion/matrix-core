package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
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
		log:  log.NewHelper(log.With(logger, "module", "bff/biz/MessageUseCase")),
	}
}

func (r *MessageUseCase) GetMailBoxLastTime(ctx context.Context) (*MailBox, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetMailBoxLastTime(ctx, uuid)
}

func (r *MessageUseCase) SetMailBoxLastTime(ctx context.Context) error {
	uuid := ctx.Value("uuid").(string)
	return r.repo.SetMailBoxLastTime(ctx, uuid)
}
