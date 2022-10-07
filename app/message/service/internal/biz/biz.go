package biz

import (
	"context"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewUserUseCase, NewCreationUseCase, NewAchievementUseCase, NewCommentUseCase, NewMessageUseCase)

type Jwt interface {
	JwtCheck(token string) (string, error)
}

type Recovery interface {
	GroupRecover(context.Context, func(ctx context.Context) error) func() error
}
