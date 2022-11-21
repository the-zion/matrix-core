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

type CommentReviewMqConsumerServer struct {
	c  rocketmq.PushConsumer
	ms *service.MessageService
}

func NewCommentReviewMqConsumerServer(conf *conf.Server, messageService *service.MessageService, logger log.Logger) *CommentReviewMqConsumerServer {
	l := log.NewHelper(log.With(logger, "server", "message/server/rocketmq-comment-review-consumer"))
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(conf.CommentMq.CommentReview.GroupName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.CommentMq.ServerAddress})),
		consumer.WithCredentials(primitive.Credentials{
			SecretKey: conf.CommentMq.SecretKey,
			AccessKey: conf.CommentMq.AccessKey,
		}),
		consumer.WithInstance("comment"),
		consumer.WithConsumeMessageBatchMaxSize(1),
		consumer.WithNamespace(conf.CommentMq.NameSpace),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		consumer.WithConsumerModel(consumer.Clustering),
	)
	if err != nil {
		l.Fatalf("init consumer error: %v", err)
	}

	delayLevel := 3
	err = c.Subscribe("review", consumer.MessageSelector{}, MqRecovery(func(ctx context.Context,
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
		case "comment_review":
			err = messageService.ToReviewCreateComment(int32(m["id"].(float64)), m["uuid"].(string))
		case "sub_comment_review":
			err = messageService.ToReviewCreateSubComment(int32(m["id"].(float64)), m["uuid"].(string))
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

	return &CommentReviewMqConsumerServer{
		c: c,
	}
}

func (s *CommentReviewMqConsumerServer) Start(_ context.Context) error {
	log.Info("mq talk review consumer starting")
	return s.c.Start()
}

func (s *CommentReviewMqConsumerServer) Stop(_ context.Context) error {
	log.Info("mq talk review consumer closing")
	return s.c.Shutdown()
}

type CommentMqConsumerServer struct {
	c  rocketmq.PushConsumer
	ms *service.MessageService
}

func NewCommentMqConsumerServer(conf *conf.Server, messageService *service.MessageService, logger log.Logger) *CommentMqConsumerServer {
	l := log.NewHelper(log.With(logger, "server", "message/server/rocketmq-comment-consumer"))
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(conf.CommentMq.Comment.GroupName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.CommentMq.ServerAddress})),
		consumer.WithCredentials(primitive.Credentials{
			SecretKey: conf.CommentMq.SecretKey,
			AccessKey: conf.CommentMq.AccessKey,
		}),
		consumer.WithInstance("comment"),
		consumer.WithConsumeMessageBatchMaxSize(1),
		consumer.WithNamespace(conf.CommentMq.NameSpace),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		consumer.WithConsumerModel(consumer.Clustering),
	)
	if err != nil {
		l.Fatalf("init consumer error: %v", err)
	}

	delayLevel := 3
	err = c.Subscribe("comment", consumer.MessageSelector{}, MqRecovery(func(ctx context.Context,
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
		case "create_comment_db_and_cache":
			err = messageService.CreateCommentDbAndCache(ctx, int32(m["id"].(float64)), int32(m["creationId"].(float64)), int32(m["creationType"].(float64)), m["uuid"].(string))
		case "create_sub_comment_db_and_cache":
			err = messageService.CreateSubCommentDbAndCache(ctx, int32(m["id"].(float64)), int32(m["rootId"].(float64)), int32(m["parentId"].(float64)), m["uuid"].(string))
		case "remove_comment_db_and_cache":
			err = messageService.RemoveCommentDbAndCache(ctx, int32(m["id"].(float64)), m["uuid"].(string))
		case "remove_sub_comment_db_and_cache":
			err = messageService.RemoveSubCommentDbAndCache(ctx, int32(m["id"].(float64)), m["uuid"].(string))
		case "set_comment_agree_db_and_cache":
			err = messageService.SetCommentAgreeDbAndCache(ctx, int32(m["id"].(float64)), int32(m["creationId"].(float64)), int32(m["creationType"].(float64)), m["uuid"].(string), m["userUuid"].(string))
		case "set_sub_comment_agree_db_and_cache":
			err = messageService.SetSubCommentAgreeDbAndCache(ctx, int32(m["id"].(float64)), m["uuid"].(string), m["userUuid"].(string))
		case "cancel_comment_agree_db_and_cache":
			err = messageService.CancelCommentAgreeDbAndCache(ctx, int32(m["id"].(float64)), int32(m["creationId"].(float64)), int32(m["creationType"].(float64)), m["uuid"].(string), m["userUuid"].(string))
		case "cancel_sub_comment_agree_db_and_cache":
			err = messageService.CancelSubCommentAgreeDbAndCache(ctx, int32(m["id"].(float64)), m["uuid"].(string), m["userUuid"].(string))
		case "add_comment_content_review_db_and_cache":
			err = messageService.AddCommentContentReviewDbAndCache(ctx, int32(m["commentId"].(float64)), int32(m["result"].(float64)), m["uuid"].(string), m["jobId"].(string), m["label"].(string), m["comment"].(string), m["kind"].(string), m["section"].(string))
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

	return &CommentMqConsumerServer{
		c: c,
	}
}

func (s *CommentMqConsumerServer) Start(_ context.Context) error {
	log.Info("mq comment consumer starting")
	return s.c.Start()
}

func (s *CommentMqConsumerServer) Stop(_ context.Context) error {
	log.Info("mq comment consumer closing")
	return s.c.Shutdown()
}
