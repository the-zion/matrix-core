package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type CommentRepo interface {
	GetLastCommentDraft(ctx context.Context, uuid string) (*CommentDraft, error)
	GetCommentList(ctx context.Context, page, creationId, creationType int32) ([]*Comment, error)
	CreateCommentDraft(ctx context.Context, uuid string) (int32, error)
	SendComment(ctx context.Context, id int32, uuid, ip string) error
}

type CommentUseCase struct {
	repo CommentRepo
	log  *log.Helper
}

func NewCommentUseCase(repo CommentRepo, logger log.Logger) *CommentUseCase {
	return &CommentUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "bff/biz/CommentUseCase")),
	}
}

func (r *CommentUseCase) CreateCommentDraft(ctx context.Context) (int32, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.CreateCommentDraft(ctx, uuid)
}

func (r *CommentUseCase) GetLastCommentDraft(ctx context.Context) (*CommentDraft, error) {
	uuid := ctx.Value("uuid").(string)
	return r.repo.GetLastCommentDraft(ctx, uuid)
}

func (r *CommentUseCase) GetCommentList(ctx context.Context, page, creationId, creationType int32) ([]*Comment, error) {
	return r.repo.GetCommentList(ctx, page, creationId, creationType)
}

func (r *CommentUseCase) SendComment(ctx context.Context, id int32) error {
	uuid := ctx.Value("uuid").(string)
	ip := ctx.Value("realIp").(string)
	return r.repo.SendComment(ctx, id, uuid, ip)
}
