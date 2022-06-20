package server

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/the-zion/matrix-core/app/message/service/internal/conf"
	"github.com/the-zion/matrix-core/app/message/service/internal/service"
)

type CodeMqConsumerServer struct {
	c  rocketmq.PushConsumer
	ms *service.MessageService
}

func NewCodeMqConsumerServer(conf *conf.Server, messageService *service.MessageService, logger log.Logger) *CodeMqConsumerServer {
	l := log.NewHelper(log.With(logger, "server", "message/server/rocketmq-code-consumer"))
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(conf.Rocketmq.Code.GroupName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.Rocketmq.Code.ServerAddress})),
		consumer.WithCredentials(primitive.Credentials{
			SecretKey: conf.Rocketmq.Code.SecretKey,
			AccessKey: conf.Rocketmq.Code.AccessKey,
		}),
		consumer.WithNamespace(conf.Rocketmq.Code.NameSpace),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		consumer.WithConsumerModel(consumer.Clustering),
	)
	if err != nil {
		l.Fatalf("init consumer error: %v", err)
	}

	err = c.Subscribe("code", consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: "phone || email",
	}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		messageService.SendCode(msgs...)
		return consumer.ConsumeRetryLater, nil
	})
	if err != nil {
		l.Fatalf("consumer subscribe error: %v", err)
	}

	return &CodeMqConsumerServer{
		c: c,
	}
}

func (s *CodeMqConsumerServer) Start(ctx context.Context) error {
	log.Info("mq code consumer starting")
	return s.c.Start()
}

func (s *CodeMqConsumerServer) Stop(ctx context.Context) error {
	log.Info("mq code consumer closing")
	return s.c.Shutdown()
}
