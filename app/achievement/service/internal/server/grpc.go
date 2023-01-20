package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	v1 "github.com/the-zion/matrix-core/api/achievement/service/v1"
	"github.com/the-zion/matrix-core/app/achievement/service/internal/conf"
	"github.com/the-zion/matrix-core/app/achievement/service/internal/service"
	"github.com/the-zion/matrix-core/pkg/responce"
)

func NewGRPCServer(c *conf.Server, achievementService *service.AchievementService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(recovery.WithHandler(func(ctx context.Context, req, err interface{}) error {
				l := log.NewHelper(log.With(logger, "message", "panic"))
				l.Error(err)
				return nil
			})),
			ratelimit.Server(),
			tracing.Server(),
			responce.Server(),
			logging.Server(log.NewFilter(logger, log.FilterLevel(log.LevelError))),
			validate.Validator(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterAchievementServer(srv, achievementService)
	return srv
}
