package server

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/the-zion/matrix-core/app/user/service/internal/conf"
	"github.com/the-zion/matrix-core/app/user/service/internal/service"
)

type MqConsumerServer struct {
	c  rocketmq.PushConsumer
	us *service.UserService
}

func NewMqConsumerServer(conf *conf.Data, userService *service.UserService, logger log.Logger) *MqConsumerServer {
	l := log.NewHelper(log.With(logger, "server", "user/server/rocketmq-consumer"))
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(conf.Rocketmq.GroupName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.Rocketmq.ServerAddress})),
		consumer.WithCredentials(primitive.Credentials{
			SecretKey: conf.Rocketmq.SecretKey,
			AccessKey: conf.Rocketmq.AccessKey,
		}),
		consumer.WithNamespace(conf.Rocketmq.NameSpace),
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
		userService.SendCode(msgs...)
		return consumer.ConsumeRetryLater, nil
	})
	if err != nil {
		l.Fatalf("consumer subscribe error: %v", err)
	}

	return &MqConsumerServer{
		c: c,
	}
}

func (s *MqConsumerServer) Start(ctx context.Context) error {
	log.Info("mq consumer starting")
	return s.c.Start()
}

func (s *MqConsumerServer) Stop(ctx context.Context) error {
	log.Info("mq consumer closing")
	return s.c.Shutdown()
}
