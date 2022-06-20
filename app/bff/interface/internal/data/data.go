package data

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	messagev1 "github.com/the-zion/matrix-core/api/message/service/v1"
	userv1 "github.com/the-zion/matrix-core/api/user/service/v1"
)

var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewMessageRepo, NewUserServiceClient, NewMessageServiceClient)

type Data struct {
	log *log.Helper
	uc  userv1.UserClient
	mc  messagev1.MessageClient
}

func NewData(logger log.Logger, uc userv1.UserClient, mc messagev1.MessageClient) (*Data, error) {
	l := log.NewHelper(log.With(logger, "module", "bff/data"))
	d := &Data{
		log: l,
		uc:  uc,
		mc:  mc,
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

func NewMessageServiceClient(r *nacos.Registry, logger log.Logger) messagev1.MessageClient {
	l := log.NewHelper(log.With(logger, "module", "bff/data/new-message-client"))
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///matrix.message.service.grpc"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			//tracing.Client(tracing.WithTracerProvider(tp)),
			recovery.Recovery(),
		),
	)
	if err != nil {
		l.Fatalf(err.Error())
	}
	c := messagev1.NewMessageClient(conn)
	return c
}
