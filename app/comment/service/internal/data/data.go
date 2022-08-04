package data

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/tencentyun/cos-go-sdk-v5"
	_ "github.com/tencentyun/cos-go-sdk-v5"
	"github.com/the-zion/matrix-core/app/comment/service/internal/biz"
	"github.com/the-zion/matrix-core/app/comment/service/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var ProviderSet = wire.NewSet(NewData, NewDB, NewRedis, NewTransaction, NewCommentRepo, NewRocketmqReviewProducer)

type ReviewMqPro struct {
	producer rocketmq.Producer
}

type Data struct {
	db          *gorm.DB
	log         *log.Helper
	cosCli      *cos.Client
	redisCli    redis.Cmdable
	reviewMqPro *ReviewMqPro
}

type contextTxKey struct{}

func (d *Data) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.db
}

func NewTransaction(d *Data) biz.Transaction {
	return d
}

func NewDB(conf *conf.Data, logger log.Logger) *gorm.DB {
	l := log.NewHelper(log.With(logger, "module", "comment/data/mysql"))

	db, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		l.Fatalf("failed opening connection to db: %v", err)
	}
	return db
}

func NewRedis(conf *conf.Data, logger log.Logger) redis.Cmdable {
	l := log.NewHelper(log.With(logger, "module", "comment/data/redis"))
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Addr,
		DB:           2,
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

func NewRocketmqReviewProducer(conf *conf.Data, logger log.Logger) *ReviewMqPro {
	l := log.NewHelper(log.With(logger, "module", "creation/data/rocketmq-review-producer"))
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.CommentMq.ServerAddress})),
		producer.WithCredentials(primitive.Credentials{
			SecretKey: conf.CommentMq.SecretKey,
			AccessKey: conf.CommentMq.AccessKey,
		}),
		producer.WithInstanceName("comment"),
		producer.WithGroupName(conf.CommentMq.CommentReview.GroupName),
		producer.WithNamespace(conf.CommentMq.NameSpace),
	)

	if err != nil {
		l.Fatalf("init producer error: %v", err)
	}

	err = p.Start()
	if err != nil {
		l.Fatalf("start producer error: %v", err)
	}
	return &ReviewMqPro{
		producer: p,
	}
}

func NewData(db *gorm.DB, cos *cos.Client, redisCmd redis.Cmdable, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "comment/data/new-data"))

	d := &Data{
		db:       db,
		redisCli: redisCmd,
		cosCli:   cos,
	}
	return d, func() {
		l.Info("closing the data resources")
	}, nil
}
