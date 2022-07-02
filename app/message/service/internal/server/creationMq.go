package server

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/the-zion/matrix-core/app/message/service/internal/conf"
	"github.com/the-zion/matrix-core/app/message/service/internal/service"
)

type ArticleDraftMqConsumerServer struct {
	c  rocketmq.PushConsumer
	ms *service.MessageService
}

func NewArticleDraftMqConsumerServer(conf *conf.Server, messageService *service.MessageService, logger log.Logger) *ArticleDraftMqConsumerServer {
	l := log.NewHelper(log.With(logger, "server", "message/server/rocketmq-article-draft-consumer"))
	rlog.SetLogLevel("warn")
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(conf.CreationMq.ArticleDraft.GroupName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.CreationMq.ServerAddress})),
		consumer.WithCredentials(primitive.Credentials{
			SecretKey: conf.CreationMq.SecretKey,
			AccessKey: conf.CreationMq.AccessKey,
		}),
		consumer.WithInstance("creation"),
		consumer.WithMaxReconsumeTimes(5),
		consumer.WithNamespace(conf.CreationMq.NameSpace),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		consumer.WithConsumerModel(consumer.Clustering),
	)
	if err != nil {
		l.Fatalf("init consumer error: %v", err)
	}

	//delayLevel := 2
	err = c.Subscribe("article_draft", consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		messageService.ToReviewArticleDraft(msgs...)
		fmt.Println(msgs)
		//concurrentCtx, _ := primitive.GetConcurrentlyCtx(ctx)
		//concurrentCtx.DelayLevelWhenNextConsume = delayLevel
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		l.Fatalf("consumer subscribe error: %v", err)
	}

	return &ArticleDraftMqConsumerServer{
		c: c,
	}
}

func (s *ArticleDraftMqConsumerServer) Start(ctx context.Context) error {
	log.Info("mq article draft consumer starting")
	return s.c.Start()
}

func (s *ArticleDraftMqConsumerServer) Stop(ctx context.Context) error {
	log.Info("mq article draft consumer closing")
	return s.c.Shutdown()
}
