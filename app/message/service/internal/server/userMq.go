package server

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
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
	rlog.SetLogLevel("warn")
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(conf.UserMq.Code.GroupName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.UserMq.ServerAddress})),
		consumer.WithCredentials(primitive.Credentials{
			SecretKey: conf.UserMq.SecretKey,
			AccessKey: conf.UserMq.AccessKey,
		}),
		consumer.WithInstance("user"),
		consumer.WithNamespace(conf.UserMq.NameSpace),
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
		return consumer.ConsumeSuccess, nil
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

type ProfileMqConsumerServer struct {
	c  rocketmq.PushConsumer
	ms *service.MessageService
}

func NewProfileMqConsumerServer(conf *conf.Server, messageService *service.MessageService, logger log.Logger) *ProfileMqConsumerServer {
	l := log.NewHelper(log.With(logger, "server", "message/server/rocketmq-profile-consumer"))
	rlog.SetLogLevel("warn")
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(conf.UserMq.Profile.GroupName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.UserMq.ServerAddress})),
		consumer.WithCredentials(primitive.Credentials{
			SecretKey: conf.UserMq.SecretKey,
			AccessKey: conf.UserMq.AccessKey,
		}),
		consumer.WithInstance("user"),
		consumer.WithNamespace(conf.UserMq.NameSpace),
		consumer.WithConsumeMessageBatchMaxSize(1),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		consumer.WithConsumerModel(consumer.Clustering),
	)
	if err != nil {
		l.Fatalf("init consumer error: %v", err)
	}

	delayLevel := 3
	err = c.Subscribe("profile", consumer.MessageSelector{},
		func(ctx context.Context,
			msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

			concurrentCtx, _ := primitive.GetConcurrentlyCtx(ctx)
			concurrentCtx.DelayLevelWhenNextConsume = delayLevel

			err := messageService.UploadProfileToCos(msgs[0])
			if err != nil {
				l.Error(err.Error())
				return consumer.ConsumeRetryLater, nil
			}
			return consumer.ConsumeSuccess, nil
		})
	if err != nil {
		l.Fatalf("consumer subscribe error: %v", err)
	}

	return &ProfileMqConsumerServer{
		c: c,
	}
}

func (s *ProfileMqConsumerServer) Start(ctx context.Context) error {
	log.Info("mq profile consumer starting")
	return s.c.Start()
}

func (s *ProfileMqConsumerServer) Stop(ctx context.Context) error {
	log.Info("mq profile consumer closing")
	return s.c.Shutdown()
}
