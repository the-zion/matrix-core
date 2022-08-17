package biz

import (
	"context"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewUserUseCase, NewAuthUseCase)

type Transaction interface {
	ExecTx(context.Context, func(ctx context.Context) error) error
}

type Recovery interface {
	GroupRecover(context.Context, func(ctx context.Context) error) func() error
}
