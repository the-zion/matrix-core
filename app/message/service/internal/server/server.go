package server

import (
	"github.com/google/wire"
)

// ProviderSet is user providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer, NewCodeMqConsumerServer, NewProfileMqConsumerServer, NewArticleReviewMqConsumerServer, NewArticleMqConsumerServer, NewTalkReviewMqConsumerServer, NewTalkMqConsumerServer, NewColumnReviewMqConsumerServer, NewColumnMqConsumerServer,
	NewAchievementMqConsumerServer, NewCommentReviewMqConsumerServer, NewCommentMqConsumerServer, NewFollowMqConsumerServer, NewCollectionsReviewMqConsumerServer, NewCollectionsMqConsumerServer)
