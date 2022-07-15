package biz

import (
	"context"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewArticleUseCase, NewTalkUseCase, NewCreationUseCase)

type Transaction interface {
	ExecTx(context.Context, func(ctx context.Context) error) error
}
