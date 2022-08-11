package biz

import (
	"context"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewUserUseCase, NewArticleUseCase, NewCreationUseCase, NewTalkUseCase, NewColumnUseCase, NewAchievementUseCase, NewNewsUseCase, NewCommentUseCase)

type Recovery interface {
	GroupRecover(context.Context, func(ctx context.Context) error) func() error
}
