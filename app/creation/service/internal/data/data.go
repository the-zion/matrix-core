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
	"github.com/the-zion/matrix-core/app/creation/service/internal/biz"
	"github.com/the-zion/matrix-core/app/creation/service/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"time"
)

var ProviderSet = wire.NewSet(NewData, NewDB, NewRedis, NewTransaction, NewRocketmqArticleDraftProducer, NewRocketmqAchievementProducer, NewCosServiceClient, NewArticleRepo, NewCreationRepo)

type ArticleDraftMqPro struct {
	producer rocketmq.Producer
}

type AchievementMqPro struct {
	producer rocketmq.Producer
}

type Data struct {
	db                *gorm.DB
	log               *log.Helper
	redisCli          redis.Cmdable
	articleDraftMqPro *ArticleDraftMqPro
	achievementMqPro  *AchievementMqPro
	cosCli            *cos.Client
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
	l := log.NewHelper(log.With(logger, "module", "creation/data/mysql"))

	db, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		l.Fatalf("failed opening connection to db: %v", err)
	}
	return db
}

func NewRedis(conf *conf.Data, logger log.Logger) redis.Cmdable {
	l := log.NewHelper(log.With(logger, "module", "creation/data/redis"))
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Addr,
		DB:           1,
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

func NewRocketmqArticleDraftProducer(conf *conf.Data, logger log.Logger) *ArticleDraftMqPro {
	l := log.NewHelper(log.With(logger, "module", "creation/data/rocketmq-article-draft-producer"))
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.CreationMq.ServerAddress})),
		producer.WithCredentials(primitive.Credentials{
			SecretKey: conf.CreationMq.SecretKey,
			AccessKey: conf.CreationMq.AccessKey,
		}),
		producer.WithInstanceName("creation"),
		producer.WithGroupName(conf.CreationMq.ArticleDraft.GroupName),
		producer.WithNamespace(conf.CreationMq.NameSpace),
	)

	if err != nil {
		l.Fatalf("init producer error: %v", err)
	}

	err = p.Start()
	if err != nil {
		l.Fatalf("start producer error: %v", err)
	}
	return &ArticleDraftMqPro{
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

func NewData(db *gorm.DB, redisCmd redis.Cmdable, cos *cos.Client, adp *ArticleDraftMqPro, ap *AchievementMqPro, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "creation/data/new-data"))

	d := &Data{
		db:                db,
		cosCli:            cos,
		redisCli:          redisCmd,
		articleDraftMqPro: adp,
		achievementMqPro:  ap,
	}
	return d, func() {
		l.Info("closing the data resources")

		err := d.articleDraftMqPro.producer.Shutdown()
		if err != nil {
			l.Errorf("shutdown article draft producer error: %v", err.Error())
		}

		err = d.achievementMqPro.producer.Shutdown()
		if err != nil {
			l.Errorf("shutdown achievement producer error: %v", err.Error())
		}
	}, nil
}
