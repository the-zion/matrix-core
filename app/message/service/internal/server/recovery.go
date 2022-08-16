package server

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-kratos/kratos/v2/log"
	"runtime"
)

func MqRecovery(fn func(ctx context.Context,
	msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error)) func(ctx context.Context,
	msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	return func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		defer func() {
			if rerr := recover(); rerr != nil {
				buf := make([]byte, 64<<10)
				n := runtime.Stack(buf, false)
				buf = buf[:n]
				log.Context(ctx).Errorf("%v: %+v\n%s\n", rerr, msgs, buf)
			}
		}()
		return fn(ctx, msgs...)
	}
}
