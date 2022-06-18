package data

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	userv1 "github.com/the-zion/matrix-core/api/user/service/v1"
)

var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewUserServiceClient)

type Data struct {
	log *log.Helper
	uc  userv1.UserClient
}

func NewData(logger log.Logger, uc userv1.UserClient) (*Data, error) {
	l := log.NewHelper(log.With(logger, "module", "bff/data"))
	d := &Data{
		log: l,
		uc:  uc,
	}
	return d, nil
}

func NewUserServiceClient(r *nacos.Registry, logger log.Logger) userv1.UserClient {
	l := log.NewHelper(log.With(logger, "module", "bff/data/new-user-client"))
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///matrix.user.service.grpc"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			//tracing.Client(tracing.WithTracerProvider(tp)),
			recovery.Recovery(),
		),
	)
	if err != nil {
		l.Fatalf(err.Error())
	}
	c := userv1.NewUserClient(conn)
	return c
}
