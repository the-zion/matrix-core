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
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/tencentyun/cos-go-sdk-v5"
	_ "github.com/tencentyun/cos-go-sdk-v5"
	achievementv1 "github.com/the-zion/matrix-core/api/achievement/service/v1"
	commentv1 "github.com/the-zion/matrix-core/api/comment/service/v1"
	creationv1 "github.com/the-zion/matrix-core/api/creation/service/v1"
	userv1 "github.com/the-zion/matrix-core/api/user/service/v1"
	"github.com/the-zion/matrix-core/app/message/service/internal/biz"
	"github.com/the-zion/matrix-core/app/message/service/internal/conf"
	"github.com/the-zion/matrix-core/pkg/trace"
	"go.opentelemetry.io/otel/propagation"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewCreationRepo, NewCommentRepo, NewMessageRepo, NewAchievementRepo, NewUserServiceClient, NewCreationServiceClient, NewAchievementServiceClient, NewCommentServiceClient, NewCosUserClient, NewCosCreationClient, NewCosCommentClient, NewJwtClient, NewJwt, NewRecovery, NewTransaction, NewRedis, NewDB)

type CosUser struct {
	cos *cos.Client
}

type CosCreation struct {
	cos      *cos.Client
	callback map[string]string
}

type CosComment struct {
	cos      *cos.Client
	callback map[string]string
}

type Jwt struct {
	key string
}

type Data struct {
	db             *gorm.DB
	log            *log.Helper
	redisCli       redis.Cmdable
	uc             userv1.UserClient
	cc             creationv1.CreationClient
	commc          commentv1.CommentClient
	ac             achievementv1.AchievementClient
	jwt            Jwt
	cosUserCli     *CosUser
	cosCreationCli *CosCreation
	cosCommentCli  *CosComment
}

type contextTxKey struct{}

func (d *Data) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

func (d *Data) GroupRecover(ctx context.Context, fn func(ctx context.Context) error) func() error {
	return func() error {
		defer func() {
			if rerr := recover(); rerr != nil {
				buf := make([]byte, 64<<10)
				n := runtime.Stack(buf, false)
				buf = buf[:n]
				d.log.Errorf("%v: %s\n", rerr, buf)
			}
		}()
		return fn(ctx)
	}
}

func (d *Data) Recover(ctx context.Context, fn func(ctx context.Context)) func() {
	return func() {
		defer func() {
			if rerr := recover(); rerr != nil {
				buf := make([]byte, 64<<10)
				n := runtime.Stack(buf, false)
				buf = buf[:n]
				d.log.Errorf("%v: %s\n", rerr, buf)
			}
		}()
		fn(ctx)
	}
}

func NewRecovery(d *Data) biz.Recovery {
	return d
}

func NewTransaction(d *Data) biz.Transaction {
	return d
}

func NewUserServiceClient(r *nacos.Registry) userv1.UserClient {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "message/data/new-user-client"))
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

func NewCreationServiceClient(r *nacos.Registry) creationv1.CreationClient {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "message/data/new-creation-client"))
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

func NewCommentServiceClient(r *nacos.Registry) commentv1.CommentClient {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "message/data/new-comment-client"))
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

func NewAchievementServiceClient(r *nacos.Registry) achievementv1.AchievementClient {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "message/data/new-achievement-client"))
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

func NewCosUserClient(conf *conf.Data) *CosUser {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "message/data/new-cos-user-client"))
	u, err := url.Parse(conf.Cos.BucketUser.BucketUrl)
	if err != nil {
		l.Errorf("fail to init cos server, error: %v", err)
	}
	b := &cos.BaseURL{BucketURL: u}
	return &CosUser{
		cos: cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  conf.Cos.BucketUser.SecretId,
				SecretKey: conf.Cos.BucketUser.SecretKey,
			},
		}),
	}
}

func NewCosCreationClient(conf *conf.Data) *CosCreation {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "message/data/new-cos-creation-client"))
	bu, err := url.Parse(conf.Cos.BucketCreation.BucketUrl)
	if err != nil {
		l.Errorf("fail to init cos server, error: %v", err)
	}
	cu, err := url.Parse(conf.Cos.BucketCreation.CiUrl)
	if err != nil {
		l.Errorf("fail to init cos server, error: %v", err)
	}
	b := &cos.BaseURL{BucketURL: bu, CIURL: cu}
	return &CosCreation{
		callback: conf.Cos.BucketCreation.Callback,
		cos: cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  conf.Cos.BucketCreation.SecretId,
				SecretKey: conf.Cos.BucketCreation.SecretKey,
			},
		}),
	}
}

func NewCosCommentClient(conf *conf.Data) *CosComment {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "message/data/new-cos-comment-client"))
	bu, err := url.Parse(conf.Cos.BucketComment.BucketUrl)
	if err != nil {
		l.Errorf("fail to init cos server, error: %v", err)
	}
	cu, err := url.Parse(conf.Cos.BucketComment.CiUrl)
	if err != nil {
		l.Errorf("fail to init cos server, error: %v", err)
	}
	b := &cos.BaseURL{BucketURL: bu, CIURL: cu}
	return &CosComment{
		callback: conf.Cos.BucketComment.Callback,
		cos: cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  conf.Cos.BucketComment.SecretId,
				SecretKey: conf.Cos.BucketComment.SecretKey,
			},
		}),
	}
}

func NewDB(conf *conf.Data) *gorm.DB {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "message/data/mysql"))

	db, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		l.Fatalf("failed opening connection to db: %v", err)
	}
	return db
}

func NewRedis(conf *conf.Data) redis.Cmdable {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "message/data/redis"))
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Addr,
		DB:           4,
		ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
		DialTimeout:  time.Second * 2,
		PoolSize:     10,
		Password:     conf.Redis.Password,
	})
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()
	err := client.Ping(timeout).Err()
	if err != nil {
		l.Fatalf("redis connect error: %v", err)
	}
	return client
}

func NewData(db *gorm.DB, redisCmd redis.Cmdable, uc userv1.UserClient, cc creationv1.CreationClient, commc commentv1.CommentClient, ac achievementv1.AchievementClient, jwt Jwt, cosUser *CosUser, cosCreation *CosCreation, cosComment *CosComment, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "message/data"))
	selector.SetGlobalSelector(p2c.NewBuilder())
	d := &Data{
		db:             db,
		log:            log.NewHelper(log.With(logger, "module", "creation/data")),
		redisCli:       redisCmd,
		uc:             uc,
		cc:             cc,
		commc:          commc,
		ac:             ac,
		jwt:            jwt,
		cosUserCli:     cosUser,
		cosCreationCli: cosCreation,
		cosCommentCli:  cosComment,
	}
	return d, func() {
		l.Info("closing the data resources")

		sqlDB, err := db.DB()
		if err != nil {
			l.Errorf("close db err: %v", err.Error())
		}

		err = sqlDB.Close()
		if err != nil {
			l.Errorf("close db err: %v", err.Error())
		}

		err = redisCmd.(*redis.Client).Close()
		if err != nil {
			l.Errorf("close redis err: %v", err.Error())
		}

	}, nil
}
