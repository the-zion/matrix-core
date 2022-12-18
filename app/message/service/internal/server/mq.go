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

type RocketMqConsumerServer struct {
	c  rocketmq.PushConsumer
	ms *service.MessageService
}

func NewRocketMqConsumerServer(conf *conf.Server, messageService *service.MessageService, logger log.Logger) *RocketMqConsumerServer {
	l := log.NewHelper(log.With(logger, "server", "message/server/rocketmq-consumer"))
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(conf.Rocketmq.GroupName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.Rocketmq.ServerAddress})),
		consumer.WithCredentials(primitive.Credentials{
			SecretKey: conf.Rocketmq.SecretKey,
			AccessKey: conf.Rocketmq.AccessKey,
		}),
		consumer.WithConsumeMessageBatchMaxSize(1),
		consumer.WithNamespace(conf.Rocketmq.NameSpace),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		consumer.WithConsumerModel(consumer.Clustering),
	)
	if err != nil {
		l.Fatalf("init consumer error: %v", err)
	}

	delayLevel := 3
	err = c.Subscribe("matrix", consumer.MessageSelector{}, MqRecovery(func(ctx context.Context,
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
		//achievement
		case "agree":
			err = messageService.SetAchievementAgree(ctx, m["uuid"].(string), m["userUuid"].(string))
		case "agree_cancel":
			err = messageService.CancelAchievementAgree(ctx, m["uuid"].(string), m["userUuid"].(string))
		case "view":
			err = messageService.SetAchievementView(ctx, m["uuid"].(string))
		case "collect":
			err = messageService.SetAchievementCollect(ctx, m["uuid"].(string))
		case "collect_cancel":
			err = messageService.CancelAchievementCollect(ctx, m["uuid"].(string))
		case "follow":
			err = messageService.SetAchievementFollow(ctx, m["follow"].(string), m["followed"].(string))
		case "follow_cancel":
			err = messageService.CancelAchievementFollow(ctx, m["follow"].(string), m["followed"].(string))
		case "add_score":
			err = messageService.AddAchievementScore(ctx, m["uuid"].(string), int32(m["score"].(float64)))
		case "set_user_medal_db_and_cache":
			err = messageService.SetUserMedalDbAndCache(ctx, m["medal"].(string), m["uuid"].(string))
		case "cancel_user_medal_db_and_cache":
			err = messageService.CancelUserMedalDbAndCache(ctx, m["medal"].(string), m["uuid"].(string))
		case "access_user_medal_db_and_cache":
			err = messageService.AccessUserMedalDbAndCache(ctx, m["medal"].(string), m["uuid"].(string))

		//CommentReview
		case "comment_review":
			err = messageService.ToReviewCreateComment(int32(m["id"].(float64)), m["uuid"].(string))
		case "sub_comment_review":
			err = messageService.ToReviewCreateSubComment(int32(m["id"].(float64)), m["uuid"].(string))

		//Comment
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

		//ArticleReview
		case "article_create_review":
			err = messageService.ToReviewCreateArticle(ctx, int32(m["id"].(float64)), m["uuid"].(string))
		case "article_edit_review":
			err = messageService.ToReviewEditArticle(ctx, int32(m["id"].(float64)), m["uuid"].(string))

		//Article
		case "create_article_db_cache_and_search":
			err = messageService.CreateArticleDbCacheAndSearch(ctx, int32(m["id"].(float64)), int32(m["auth"].(float64)), m["uuid"].(string))
		case "edit_article_cos_and_search":
			err = messageService.EditArticleCosAndSearch(ctx, int32(m["id"].(float64)), int32(m["auth"].(float64)), m["uuid"].(string))
		case "delete_article_cache_and_search":
			err = messageService.DeleteArticleCacheAndSearch(ctx, int32(m["id"].(float64)), m["uuid"].(string))
		case "set_article_view_db_and_cache":
			err = messageService.SetArticleViewDbAndCache(ctx, int32(m["id"].(float64)), m["uuid"].(string))
		case "set_article_agree_db_and_cache":
			err = messageService.SetArticleAgreeDbAndCache(ctx, int32(m["id"].(float64)), m["uuid"].(string), m["userUuid"].(string))
		case "cancel_article_agree_db_and_cache":
			err = messageService.CancelArticleAgreeDbAndCache(ctx, int32(m["id"].(float64)), m["uuid"].(string), m["userUuid"].(string))
		case "set_article_collect_db_and_cache":
			err = messageService.SetArticleCollectDbAndCache(ctx, int32(m["id"].(float64)), int32(m["collectionsId"].(float64)), m["uuid"].(string), m["userUuid"].(string))
		case "cancel_article_collect_db_and_cache":
			err = messageService.CancelArticleCollectDbAndCache(ctx, int32(m["id"].(float64)), m["uuid"].(string), m["userUuid"].(string))
		case "add_article_image_review_db_and_cache":
			err = messageService.AddArticleImageReviewDbAndCache(ctx, int32(m["creationId"].(float64)), int32(m["score"].(float64)), int32(m["result"].(float64)), m["kind"].(string), m["uid"].(string), m["uuid"].(string), m["jobId"].(string), m["label"].(string), m["category"].(string), m["subLabel"].(string))
		case "add_article_content_review_db_and_cache":
			err = messageService.AddArticleContentReviewDbAndCache(ctx, int32(m["creationId"].(float64)), int32(m["result"].(float64)), m["uuid"].(string), m["jobId"].(string), m["label"].(string), m["title"].(string), m["kind"].(string), m["section"].(string))

		//TalkReview
		case "talk_create_review":
			err = messageService.ToReviewCreateTalk(int32(m["id"].(float64)), m["uuid"].(string))
		case "talk_edit_review":
			err = messageService.ToReviewEditTalk(int32(m["id"].(float64)), m["uuid"].(string))

		//Talk
		case "create_talk_db_cache_and_search":
			err = messageService.CreateTalkDbCacheAndSearch(ctx, int32(m["id"].(float64)), int32(m["auth"].(float64)), m["uuid"].(string))
		case "edit_talk_cos_and_search":
			err = messageService.EditTalkCosAndSearch(ctx, int32(m["id"].(float64)), int32(m["auth"].(float64)), m["uuid"].(string))
		case "delete_talk_cache_and_search":
			err = messageService.DeleteTalkCacheAndSearch(ctx, int32(m["id"].(float64)), m["uuid"].(string))
		case "set_talk_view_db_and_cache":
			err = messageService.SetTalkViewDbAndCache(ctx, int32(m["id"].(float64)), m["uuid"].(string))
		case "set_talk_agree_db_and_cache":
			err = messageService.SetTalkAgreeDbAndCache(ctx, int32(m["id"].(float64)), m["uuid"].(string), m["userUuid"].(string))
		case "cancel_talk_agree_db_and_cache":
			err = messageService.CancelTalkAgreeDbAndCache(ctx, int32(m["id"].(float64)), m["uuid"].(string), m["userUuid"].(string))
		case "set_talk_collect_db_and_cache":
			err = messageService.SetTalkCollectDbAndCache(ctx, int32(m["id"].(float64)), int32(m["collectionsId"].(float64)), m["uuid"].(string), m["userUuid"].(string))
		case "cancel_talk_collect_db_and_cache":
			err = messageService.CancelTalkCollectDbAndCache(ctx, int32(m["id"].(float64)), m["uuid"].(string), m["userUuid"].(string))
		case "add_talk_image_review_db_and_cache":
			err = messageService.AddTalkImageReviewDbAndCache(ctx, int32(m["creationId"].(float64)), int32(m["score"].(float64)), int32(m["result"].(float64)), m["kind"].(string), m["uid"].(string), m["uuid"].(string), m["jobId"].(string), m["label"].(string), m["category"].(string), m["subLabel"].(string))
		case "add_talk_content_review_db_and_cache":
			err = messageService.AddTalkContentReviewDbAndCache(ctx, int32(m["creationId"].(float64)), int32(m["result"].(float64)), m["uuid"].(string), m["jobId"].(string), m["label"].(string), m["title"].(string), m["kind"].(string), m["section"].(string))

		//ColumnReview
		case "column_create_review":
			err = messageService.ToReviewCreateColumn(int32(m["id"].(float64)), m["uuid"].(string))
		case "column_edit_review":
			err = messageService.ToReviewEditColumn(int32(m["id"].(float64)), m["uuid"].(string))

		//Column
		case "create_column_db_cache_and_search":
			err = messageService.CreateColumnDbCacheAndSearch(ctx, int32(m["id"].(float64)), int32(m["auth"].(float64)), m["uuid"].(string))
		case "edit_column_cos_and_search":
			err = messageService.EditColumnCosAndSearch(ctx, int32(m["id"].(float64)), int32(m["auth"].(float64)), m["uuid"].(string))
		case "delete_column_cache_and_search":
			err = messageService.DeleteColumnCacheAndSearch(ctx, int32(m["id"].(float64)), m["uuid"].(string))
		case "set_column_view_db_and_cache":
			err = messageService.SetColumnViewDbAndCache(ctx, int32(m["id"].(float64)), m["uuid"].(string))
		case "set_column_agree_db_and_cache":
			err = messageService.SetColumnAgreeDbAndCache(ctx, int32(m["id"].(float64)), m["uuid"].(string), m["userUuid"].(string))
		case "cancel_column_agree_db_and_cache":
			err = messageService.CancelColumnAgreeDbAndCache(ctx, int32(m["id"].(float64)), m["uuid"].(string), m["userUuid"].(string))
		case "set_column_collect_db_and_cache":
			err = messageService.SetColumnCollectDbAndCache(ctx, int32(m["id"].(float64)), int32(m["collectionsId"].(float64)), m["uuid"].(string), m["userUuid"].(string))
		case "cancel_column_collect_db_and_cache":
			err = messageService.CancelColumnCollectDbAndCache(ctx, int32(m["id"].(float64)), m["uuid"].(string), m["userUuid"].(string))
		case "add_column_includes_db_and_cache":
			err = messageService.AddColumnIncludesDbAndCache(ctx, int32(m["id"].(float64)), int32(m["articleId"].(float64)), m["uuid"].(string))
		case "delete_column_includes_db_and_cache":
			err = messageService.DeleteColumnIncludesDbAndCache(ctx, int32(m["id"].(float64)), int32(m["articleId"].(float64)), m["uuid"].(string))
		case "set_column_subscribe_db_and_cache":
			err = messageService.SetColumnSubscribeDbAndCache(ctx, int32(m["id"].(float64)), m["uuid"].(string))
		case "cancel_column_subscribe_db_and_cache":
			err = messageService.CancelColumnSubscribeDbAndCache(ctx, int32(m["id"].(float64)), m["uuid"].(string))
		case "add_column_image_review_db_and_cache":
			err = messageService.AddColumnImageReviewDbAndCache(ctx, int32(m["creationId"].(float64)), int32(m["score"].(float64)), int32(m["result"].(float64)), m["kind"].(string), m["uid"].(string), m["uuid"].(string), m["jobId"].(string), m["label"].(string), m["category"].(string), m["subLabel"].(string))
		case "add_column_content_review_db_and_cache":
			err = messageService.AddColumnContentReviewDbAndCache(ctx, int32(m["creationId"].(float64)), int32(m["result"].(float64)), m["uuid"].(string), m["jobId"].(string), m["label"].(string), m["title"].(string), m["kind"].(string), m["section"].(string))

		//CollectionsReview
		case "collections_create_review":
			err = messageService.ToReviewCreateCollections(int32(m["id"].(float64)), m["uuid"].(string))
		case "collections_edit_review":
			err = messageService.ToReviewEditCollections(int32(m["id"].(float64)), m["uuid"].(string))

		//Collections
		case "create_collections_db_and_cache":
			err = messageService.CreateCollectionsDbAndCache(ctx, int32(m["id"].(float64)), int32(m["auth"].(float64)), m["uuid"].(string))
		case "edit_collections_cos":
			err = messageService.EditCollectionsCos(ctx, int32(m["id"].(float64)), int32(m["auth"].(float64)), m["uuid"].(string))
		case "delete_collections_cache":
			err = messageService.DeleteCollectionsCache(ctx, int32(m["id"].(float64)), m["uuid"].(string))
		case "add_collections_content_review_db_and_cache":
			err = messageService.AddCollectionsContentReviewDbAndCache(ctx, int32(m["creationId"].(float64)), int32(m["result"].(float64)), m["uuid"].(string), m["jobId"].(string), m["label"].(string), m["title"].(string), m["kind"].(string), m["section"].(string))

		//Follow
		case "set_follow_db_and_cache":
			err = messageService.SetFollowDbAndCache(ctx, m["uuid"].(string), m["userId"].(string))
		case "cancel_follow_db_and_cache":
			err = messageService.CancelFollowDbAndCache(ctx, m["uuid"].(string), m["userId"].(string))

		//Picture
		case "add_avatar_review_db_and_cache":
			err = messageService.AddAvatarReviewDbAndCache(ctx, int32(m["score"].(float64)), int32(m["result"].(float64)), m["uuid"].(string), m["jobId"].(string), m["label"].(string), m["category"].(string), m["subLabel"].(string))
		case "add_cover_review_db_and_cache":
			err = messageService.AddCoverReviewDbAndCache(ctx, int32(m["score"].(float64)), int32(m["result"].(float64)), m["uuid"].(string), m["jobId"].(string), m["label"].(string), m["category"].(string), m["subLabel"].(string))

		//Profile:
		case "user_profile_update":
			err = messageService.UploadProfileToCos(m)
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

	return &RocketMqConsumerServer{
		c: c,
	}
}

func (s *RocketMqConsumerServer) Start(_ context.Context) error {
	log.Info("mq consumer starting")
	return s.c.Start()
}

func (s *RocketMqConsumerServer) Stop(_ context.Context) error {
	log.Info("mq consumer closing")
	return s.c.Shutdown()
}
