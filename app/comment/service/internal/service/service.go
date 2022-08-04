package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "github.com/the-zion/matrix-core/api/comment/service/v1"
	"github.com/the-zion/matrix-core/app/comment/service/internal/biz"
)

var ProviderSet = wire.NewSet(NewCommentService)

type CommentService struct {
	v1.UnimplementedCommentServer
	cc  *biz.CommentUseCase
	log *log.Helper
}

func NewCommentService(cc *biz.CommentUseCase, logger log.Logger) *CommentService {
	return &CommentService{
		log: log.NewHelper(log.With(logger, "module", "comment/service")),
		cc:  cc,
	}
}
