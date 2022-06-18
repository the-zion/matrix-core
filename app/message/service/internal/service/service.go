package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "github.com/the-zion/matrix-core/api/message/service/v1"
	"github.com/the-zion/matrix-core/app/message/service/internal/biz"
)

var ProviderSet = wire.NewSet(NewMessageService)

type MessageService struct {
	v1.UnimplementedMessageServer
	uc  *biz.UserUseCase
	log *log.Helper
}

func NewMessageService(uc *biz.UserUseCase, logger log.Logger) *MessageService {
	return &MessageService{
		log: log.NewHelper(log.With(logger, "module", "message/service")),
		uc:  uc,
	}
}
