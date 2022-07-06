package biz

import (
	"context"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewArticleUseCase, NewCreationUseCase)

type Transaction interface {
	ExecTx(context.Context, func(ctx context.Context) error) error
}
