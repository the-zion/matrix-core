package server

import (
	"context"
	"encoding/json"
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
	}, MqRecovery(func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

		messageService.SendCode(msgs...)
		return consumer.ConsumeSuccess, nil
	}))
	if err != nil {
		l.Fatalf("consumer subscribe error: %v", err)
	}

	return &CodeMqConsumerServer{
		c: c,
	}
}

func (s *CodeMqConsumerServer) Start(_ context.Context) error {
	log.Info("mq code consumer starting")
	return s.c.Start()
}

func (s *CodeMqConsumerServer) Stop(_ context.Context) error {
	log.Info("mq code consumer closing")
	return s.c.Shutdown()
}

type ProfileMqConsumerServer struct {
	c  rocketmq.PushConsumer
	ms *service.MessageService
}

func NewProfileMqConsumerServer(conf *conf.Server, messageService *service.MessageService, logger log.Logger) *ProfileMqConsumerServer {
	l := log.NewHelper(log.With(logger, "server", "message/server/rocketmq-profile-consumer"))
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
		MqRecovery(func(ctx context.Context,
			msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

			concurrentCtx, _ := primitive.GetConcurrentlyCtx(ctx)
			concurrentCtx.DelayLevelWhenNextConsume = delayLevel

			err := messageService.UploadProfileToCos(msgs[0])
			if err != nil {
				l.Error(err.Error())
				return consumer.ConsumeRetryLater, nil
			}
			return consumer.ConsumeSuccess, nil
		}))
	if err != nil {
		l.Fatalf("consumer subscribe error: %v", err)
	}

	return &ProfileMqConsumerServer{
		c: c,
	}
}

func (s *ProfileMqConsumerServer) Start(_ context.Context) error {
	log.Info("mq profile consumer starting")
	return s.c.Start()
}

func (s *ProfileMqConsumerServer) Stop(_ context.Context) error {
	log.Info("mq profile consumer closing")
	return s.c.Shutdown()
}

type FollowMqConsumerServer struct {
	c  rocketmq.PushConsumer
	ms *service.MessageService
}

func NewFollowMqConsumerServer(conf *conf.Server, messageService *service.MessageService, logger log.Logger) *FollowMqConsumerServer {
	l := log.NewHelper(log.With(logger, "server", "message/server/rocketmq-follow-consumer"))
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(conf.UserMq.Follow.GroupName),
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
	err = c.Subscribe("follow", consumer.MessageSelector{},
		MqRecovery(func(ctx context.Context,
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
			case "set_follow_db_and_cache":
				err = messageService.SetFollowDbAndCache(ctx, m["uuid"].(string), m["userId"].(string))
			case "cancel_follow_db_and_cache":
				err = messageService.CancelFollowDbAndCache(ctx, m["uuid"].(string), m["userId"].(string))
			}

			if err != nil {
				l.Error(err.Error())
				return consumer.ConsumeRetryLater, nil
			}
			return consumer.ConsumeSuccess, nil
		}))
	if err != nil {
		l.Fatalf("consumer subscribe error: %v", err)
	}

	return &FollowMqConsumerServer{
		c: c,
	}
}

func (s *FollowMqConsumerServer) Start(_ context.Context) error {
	log.Info("mq follow consumer starting")
	return s.c.Start()
}

func (s *FollowMqConsumerServer) Stop(_ context.Context) error {
	log.Info("mq follow consumer closing")
	return s.c.Shutdown()
}

type PictureMqConsumerServer struct {
	c  rocketmq.PushConsumer
	ms *service.MessageService
}

func NewPictureMqConsumerServer(conf *conf.Server, messageService *service.MessageService, logger log.Logger) *PictureMqConsumerServer {
	l := log.NewHelper(log.With(logger, "server", "message/server/rocketmq-picture-consumer"))
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(conf.UserMq.Picture.GroupName),
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
	err = c.Subscribe("picture", consumer.MessageSelector{},
		MqRecovery(func(ctx context.Context,
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
			case "add_avatar_review_db_and_cache":
				err = messageService.AddAvatarReviewDbAndCache(ctx, int32(m["score"].(float64)), int32(m["result"].(float64)), m["uuid"].(string), m["jobId"].(string), m["label"].(string), m["category"].(string), m["subLabel"].(string))
			case "add_cover_review_db_and_cache":
				err = messageService.AddCoverReviewDbAndCache(ctx, int32(m["score"].(float64)), int32(m["result"].(float64)), m["uuid"].(string), m["jobId"].(string), m["label"].(string), m["category"].(string), m["subLabel"].(string))
			}

			if err != nil {
				l.Error(err.Error())
				return consumer.ConsumeRetryLater, nil
			}
			return consumer.ConsumeSuccess, nil
		}))
	if err != nil {
		l.Fatalf("consumer subscribe error: %v", err)
	}

	return &PictureMqConsumerServer{
		c: c,
	}
}

func (s *PictureMqConsumerServer) Start(_ context.Context) error {
	log.Info("mq picture consumer starting")
	return s.c.Start()
}

func (s *PictureMqConsumerServer) Stop(_ context.Context) error {
	log.Info("mq picture consumer closing")
	return s.c.Shutdown()
}
