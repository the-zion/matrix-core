package data

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/p2c"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	achievementv1 "github.com/the-zion/matrix-core/api/achievement/service/v1"
	commentv1 "github.com/the-zion/matrix-core/api/comment/service/v1"
	creationv1 "github.com/the-zion/matrix-core/api/creation/service/v1"
	messagev1 "github.com/the-zion/matrix-core/api/message/service/v1"
	userv1 "github.com/the-zion/matrix-core/api/user/service/v1"
	"github.com/the-zion/matrix-core/app/bff/interface/internal/biz"
	"github.com/the-zion/matrix-core/pkg/trace"
	"go.opentelemetry.io/otel/propagation"
	"runtime"
)

var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewCreationRepo, NewArticleRepo, NewTalkRepo, NewColumnRepo, NewNewsRepo, NewAchievementRepo, NewCommentRepo, NewMessageRepo, NewUserServiceClient, NewCreationServiceClient, NewMessageServiceClient, NewAchievementServiceClient, NewCommentServiceClient, NewRecovery)

type Data struct {
	log   *log.Helper
	uc    userv1.UserClient
	cc    creationv1.CreationClient
	mc    messagev1.MessageClient
	ac    achievementv1.AchievementClient
	commc commentv1.CommentClient
}

func (d *Data) GroupRecover(ctx context.Context, fn func(ctx context.Context) error) func() error {
	return func() error {
		defer func() {
			if rerr := recover(); rerr != nil {
				buf := make([]byte, 64<<10)
				n := runtime.Stack(buf, false)
				buf = buf[:n]
				log.Context(ctx).Errorf("%v: %s\n", rerr, buf)
			}
		}()
		return fn(ctx)
	}
}

func NewRecovery(d *Data) biz.Recovery {
	return d
}

func NewData(logger log.Logger, uc userv1.UserClient, cc creationv1.CreationClient, mc messagev1.MessageClient, ac achievementv1.AchievementClient, commc commentv1.CommentClient) (*Data, error) {
	l := log.NewHelper(log.With(logger, "module", "bff/data"))
	selector.SetGlobalSelector(p2c.NewBuilder())
	d := &Data{
		log:   l,
		uc:    uc,
		cc:    cc,
		mc:    mc,
		ac:    ac,
		commc: commc,
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
			recovery.Recovery(),
			circuitbreaker.Client(),
			tracing.Client(tracing.WithPropagator(propagation.NewCompositeTextMapPropagator(trace.Metadata{}, propagation.Baggage{}, propagation.TraceContext{}))),
		),
	)
	if err != nil {
		l.Fatalf(err.Error())
	}
	c := userv1.NewUserClient(conn)
	return c
}

func NewCreationServiceClient(r *nacos.Registry, logger log.Logger) creationv1.CreationClient {
	l := log.NewHelper(log.With(logger, "module", "bff/data/new-creation-client"))
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///matrix.creation.service.grpc"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			recovery.Recovery(),
			circuitbreaker.Client(),
			tracing.Client(tracing.WithPropagator(propagation.NewCompositeTextMapPropagator(trace.Metadata{}, propagation.Baggage{}, propagation.TraceContext{}))),
		),
	)
	if err != nil {
		l.Fatalf(err.Error())
	}
	c := creationv1.NewCreationClient(conn)
	return c
}

func NewMessageServiceClient(r *nacos.Registry, logger log.Logger) messagev1.MessageClient {
	l := log.NewHelper(log.With(logger, "module", "bff/data/new-message-client"))
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///matrix.message.service.grpc"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			recovery.Recovery(),
			circuitbreaker.Client(),
			tracing.Client(tracing.WithPropagator(propagation.NewCompositeTextMapPropagator(trace.Metadata{}, propagation.Baggage{}, propagation.TraceContext{}))),
		),
	)
	if err != nil {
		l.Fatalf(err.Error())
	}
	c := messagev1.NewMessageClient(conn)
	return c
}

func NewAchievementServiceClient(r *nacos.Registry, logger log.Logger) achievementv1.AchievementClient {
	l := log.NewHelper(log.With(logger, "module", "bff/data/new-achievement-client"))
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///matrix.achievement.service.grpc"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			recovery.Recovery(),
			circuitbreaker.Client(),
			tracing.Client(tracing.WithPropagator(propagation.NewCompositeTextMapPropagator(trace.Metadata{}, propagation.Baggage{}, propagation.TraceContext{}))),
		),
	)
	if err != nil {
		l.Fatalf(err.Error())
	}
	c := achievementv1.NewAchievementClient(conn)
	return c
}

func NewCommentServiceClient(r *nacos.Registry, logger log.Logger) commentv1.CommentClient {
	l := log.NewHelper(log.With(logger, "module", "bff/data/new-comment-client"))
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///matrix.comment.service.grpc"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			recovery.Recovery(),
			circuitbreaker.Client(),
			tracing.Client(tracing.WithPropagator(propagation.NewCompositeTextMapPropagator(trace.Metadata{}, propagation.Baggage{}, propagation.TraceContext{}))),
		),
	)
	if err != nil {
		l.Fatalf(err.Error())
	}
	c := commentv1.NewCommentClient(conn)
	return c
}
