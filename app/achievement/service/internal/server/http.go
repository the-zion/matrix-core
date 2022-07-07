package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/the-zion/matrix-core/api/achievement/service/v1"
	"github.com/the-zion/matrix-core/app/achievement/service/internal/conf"
	"github.com/the-zion/matrix-core/app/achievement/service/internal/service"
)

// NewHTTPServer new a HTTP user.
func NewHTTPServer(c *conf.Server, achievementService *service.AchievementService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterAchievementHTTPServer(srv, achievementService)
	return srv
}
