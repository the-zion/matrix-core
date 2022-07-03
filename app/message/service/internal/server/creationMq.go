package server

import (
	"context"
	"encoding/json"
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
		consumer.WithConsumeMessageBatchMaxSize(1),
		consumer.WithNamespace(conf.CreationMq.NameSpace),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		consumer.WithConsumerModel(consumer.Clustering),
	)
	if err != nil {
		l.Fatalf("init consumer error: %v", err)
	}

	delayLevel := 3
	err = c.Subscribe("article_draft", consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

		concurrentCtx, _ := primitive.GetConcurrentlyCtx(ctx)
		concurrentCtx.DelayLevelWhenNextConsume = delayLevel

		err := messageService.ToReviewArticleDraft(msgs[0])
		if err != nil {
			l.Error(err.Error())
			return consumer.ConsumeRetryLater, nil
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		l.Fatalf("consumer subscribe error: %v", err)
	}

	return &ArticleDraftMqConsumerServer{
		c: c,
	}
}

func (s *ArticleDraftMqConsumerServer) Start(_ context.Context) error {
	log.Info("mq article draft consumer starting")
	return s.c.Start()
}

func (s *ArticleDraftMqConsumerServer) Stop(_ context.Context) error {
	log.Info("mq article draft consumer closing")
	return s.c.Shutdown()
}

type ArticleMqConsumerServer struct {
	c  rocketmq.PushConsumer
	ms *service.MessageService
}

func NewArticleMqConsumerServer(conf *conf.Server, messageService *service.MessageService, logger log.Logger) *ArticleMqConsumerServer {
	l := log.NewHelper(log.With(logger, "server", "message/server/rocketmq-article-consumer"))
	rlog.SetLogLevel("warn")
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(conf.CreationMq.Article.GroupName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.CreationMq.ServerAddress})),
		consumer.WithCredentials(primitive.Credentials{
			SecretKey: conf.CreationMq.SecretKey,
			AccessKey: conf.CreationMq.AccessKey,
		}),
		consumer.WithInstance("creation"),
		consumer.WithConsumeMessageBatchMaxSize(1),
		consumer.WithNamespace(conf.CreationMq.NameSpace),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		consumer.WithConsumerModel(consumer.Clustering),
	)
	if err != nil {
		l.Fatalf("init consumer error: %v", err)
	}

	delayLevel := 3
	err = c.Subscribe("article", consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

		concurrentCtx, _ := primitive.GetConcurrentlyCtx(ctx)
		concurrentCtx.DelayLevelWhenNextConsume = delayLevel

		var err error
		fmt.Println(msgs[0])
		m := map[string]interface{}{}
		err = json.Unmarshal(msgs[0].Body, &m)
		if err != nil {
			l.Errorf("fail to unmarshal msg: %s", err.Error())
			return consumer.ConsumeRetryLater, nil
		}

		switch m["mode"].(string) {
		case "create_article_cache_and_search":
			err = messageService.CreateArticleCacheAndSearch(ctx, m["Uuid"].(string), int32(m["Id"].(float64)))
		}

		//err = messageService.ToReviewArticleDraft(msgs[0])
		if err != nil {
			l.Error(err.Error())
			return consumer.ConsumeRetryLater, nil
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		l.Fatalf("consumer subscribe error: %v", err)
	}

	return &ArticleMqConsumerServer{
		c: c,
	}
}

func (s *ArticleMqConsumerServer) Start(_ context.Context) error {
	log.Info("mq article draft consumer starting")
	return s.c.Start()
}

func (s *ArticleMqConsumerServer) Stop(_ context.Context) error {
	log.Info("mq article draft consumer closing")
	return s.c.Shutdown()
}
