package request

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

func Server() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if header, ok := transport.FromServerContext(ctx); ok {
				ctx = context.WithValue(ctx, "uuid", header.RequestHeader().Get("uuid"))
				ctx = context.WithValue(ctx, "realIp", header.RequestHeader().Get("realIp"))
			}
			reply, err = handler(ctx, req)
			return
		}
	}
}
