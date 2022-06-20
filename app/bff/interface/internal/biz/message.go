package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type MessageRepo interface {
	AvatarReview(ctx context.Context, ar *AvatarReview) error
}

type MessageUseCase struct {
	repo MessageRepo
	log  *log.Helper
}

func NewMessageUseCase(repo MessageRepo, logger log.Logger) *MessageUseCase {
	return &MessageUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "bff/biz/MessageUseCase")),
	}
}

func (r *MessageUseCase) AvatarReview(ctx context.Context, ar *AvatarReview) error {
	err := r.repo.AvatarReview(ctx, ar)
	if err != nil {
		return err
	}
	return nil
}
