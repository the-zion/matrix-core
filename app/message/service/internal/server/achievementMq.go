package server

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/the-zion/matrix-core/app/message/service/internal/conf"
	"github.com/the-zion/matrix-core/app/message/service/internal/service"
)

type AchievementMqConsumerServer struct {
	c  rocketmq.PushConsumer
	ms *service.MessageService
}

func NewAchievementMqConsumerServer(conf *conf.Server, messageService *service.MessageService, logger log.Logger) *AchievementMqConsumerServer {
	l := log.NewHelper(log.With(logger, "server", "message/server/rocketmq-achievement-consumer"))
	rlog.SetLogLevel("warn")
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(conf.AchievementMq.Achievement.GroupName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.AchievementMq.ServerAddress})),
		consumer.WithCredentials(primitive.Credentials{
			SecretKey: conf.AchievementMq.SecretKey,
			AccessKey: conf.AchievementMq.AccessKey,
		}),
		consumer.WithInstance("achievement"),
		consumer.WithConsumeMessageBatchMaxSize(1),
		consumer.WithNamespace(conf.AchievementMq.NameSpace),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		consumer.WithConsumerModel(consumer.Clustering),
	)
	if err != nil {
		l.Fatalf("init consumer error: %v", err)
	}

	delayLevel := 3
	err = c.Subscribe("achievement", consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

		concurrentCtx, _ := primitive.GetConcurrentlyCtx(ctx)
		concurrentCtx.DelayLevelWhenNextConsume = delayLevel

		var err error
		m := map[string]interface{}{}
		err = json.Unmarshal(msgs[0].Body, &m)
		if err != nil {
			l.Errorf("fail to unmarshal msg: %s", err.Error())
			return consumer.ConsumeRetryLater, nil
		}

		mode := m["mode"].(string)
		switch mode {
		case "agree":
			err = messageService.SetAchievementAgree(ctx, m["uuid"].(string))
		case "agree_cancel":
			err = messageService.CancelAchievementAgree(ctx, m["uuid"].(string))
		case "view":
			err = messageService.SetAchievementView(ctx, m["uuid"].(string))
		case "collect":
			err = messageService.SetAchievementCollect(ctx, m["uuid"].(string))
		case "collect_cancel":
			err = messageService.CancelAchievementCollect(ctx, m["uuid"].(string))
		}

		if err != nil {
			l.Error(err.Error())
			return consumer.ConsumeRetryLater, nil
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		l.Fatalf("consumer subscribe error: %v", err)
	}

	return &AchievementMqConsumerServer{
		c: c,
	}
}

func (s *AchievementMqConsumerServer) Start(_ context.Context) error {
	log.Info("mq achievement consumer starting")
	return s.c.Start()
}

func (s *AchievementMqConsumerServer) Stop(_ context.Context) error {
	log.Info("mq achievement consumer closing")
	return s.c.Shutdown()
}
