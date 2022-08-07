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
	"net/http"
	"net/url"
	"time"
)

var ProviderSet = wire.NewSet(NewData, NewDB, NewRedis, NewTransaction, NewCommentRepo, NewRocketmqCommentProducer, NewRocketmqReviewProducer, NewRocketmqAchievementProducer, NewCosServiceClient)

type ReviewMqPro struct {
	producer rocketmq.Producer
}

type CommentMqPro struct {
	producer rocketmq.Producer
}

type AchievementMqPro struct {
	producer rocketmq.Producer
}

type Data struct {
	db               *gorm.DB
	log              *log.Helper
	cosCli           *cos.Client
	redisCli         redis.Cmdable
	reviewMqPro      *ReviewMqPro
	commonMqPro      *CommentMqPro
	achievementMqPro *AchievementMqPro
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
		DB:           3,
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

func NewCosServiceClient(conf *conf.Data, logger log.Logger) *cos.Client {
	l := log.NewHelper(log.With(logger, "module", "creation/data/new-cos-client"))
	u, err := url.Parse(conf.Cos.Url)
	if err != nil {
		l.Errorf("fail to init cos server, error: %v", err)
	}
	b := &cos.BaseURL{BucketURL: u}
	return cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  conf.Cos.SecretId,
			SecretKey: conf.Cos.SecretKey,
		},
	})
}

func NewRocketmqCommentProducer(conf *conf.Data, logger log.Logger) *CommentMqPro {
	l := log.NewHelper(log.With(logger, "module", "creation/data/rocketmq-comment-producer"))
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.CommentMq.ServerAddress})),
		producer.WithCredentials(primitive.Credentials{
			SecretKey: conf.CommentMq.SecretKey,
			AccessKey: conf.CommentMq.AccessKey,
		}),
		producer.WithInstanceName("comment"),
		producer.WithGroupName(conf.CommentMq.Comment.GroupName),
		producer.WithNamespace(conf.CommentMq.NameSpace),
	)

	if err != nil {
		l.Fatalf("init producer error: %v", err)
	}

	err = p.Start()
	if err != nil {
		l.Fatalf("start producer error: %v", err)
	}
	return &CommentMqPro{
		producer: p,
	}
}

func NewRocketmqAchievementProducer(conf *conf.Data, logger log.Logger) *AchievementMqPro {
	l := log.NewHelper(log.With(logger, "module", "creation/data/rocketmq-achievement-producer"))
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.AchievementMq.ServerAddress})),
		producer.WithCredentials(primitive.Credentials{
			SecretKey: conf.AchievementMq.SecretKey,
			AccessKey: conf.AchievementMq.AccessKey,
		}),
		producer.WithInstanceName("achievement"),
		producer.WithGroupName(conf.AchievementMq.Achievement.GroupName),
		producer.WithNamespace(conf.AchievementMq.NameSpace),
	)

	if err != nil {
		l.Fatalf("init producer error: %v", err)
	}

	err = p.Start()
	if err != nil {
		l.Fatalf("start producer error: %v", err)
	}

	return &AchievementMqPro{
		producer: p,
	}
}

func NewData(db *gorm.DB, cos *cos.Client, redisCmd redis.Cmdable, rm *ReviewMqPro, cm *CommentMqPro, ap *AchievementMqPro, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "comment/data/new-data"))

	d := &Data{
		db:               db,
		redisCli:         redisCmd,
		cosCli:           cos,
		reviewMqPro:      rm,
		achievementMqPro: ap,
		commonMqPro:      cm,
	}
	return d, func() {
		err := d.commonMqPro.producer.Shutdown()
		if err != nil {
			l.Errorf("shutdown comment producer error: %v", err.Error())
		}

		err = d.reviewMqPro.producer.Shutdown()
		if err != nil {
			l.Errorf("shutdown comment review producer error: %v", err.Error())
		}

		err = d.achievementMqPro.producer.Shutdown()
		if err != nil {
			l.Errorf("shutdown achievement producer error: %v", err.Error())
		}

		l.Info("closing the data resources")
	}, nil
}
