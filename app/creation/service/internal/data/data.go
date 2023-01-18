package data

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/the-zion/matrix-core/app/creation/service/internal/biz"
	"github.com/the-zion/matrix-core/app/creation/service/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

var ProviderSet = wire.NewSet(NewData, NewDB, NewRedis, NewTransaction, NewRocketmqProducer, NewCosServiceClient, NewElasticsearch, NewNewsClient, NewArticleRepo, NewTalkRepo, NewCreationRepo, NewColumnRepo, NewNewsRepo, NewRecovery)

type MqPro struct {
	producer rocketmq.Producer
}

type ElasticSearch struct {
	es *elasticsearch.Client
}

type NewsClient struct {
	url string
}

type Data struct {
	db            *gorm.DB
	log           *log.Helper
	redisCli      redis.Cmdable
	mqPro         *MqPro
	cosCli        *cos.Client
	elasticSearch *ElasticSearch
	newsCli       *NewsClient
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

func NewDB(conf *conf.Data) *gorm.DB {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "creation/data/mysql"))

	db, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		l.Fatalf("failed opening connection to db: %v", err)
	}
	return db
}

func NewRedis(conf *conf.Data) redis.Cmdable {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "creation/data/redis"))
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

func NewCosServiceClient(conf *conf.Data) *cos.Client {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "creation/data/new-cos-client"))
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

func NewRocketmqProducer(conf *conf.Data) *MqPro {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "creation/data/rocketmq-producer"))
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.Rocketmq.ServerAddress})),
		producer.WithCredentials(primitive.Credentials{
			SecretKey: conf.Rocketmq.SecretKey,
			AccessKey: conf.Rocketmq.AccessKey,
		}),
		producer.WithGroupName(conf.Rocketmq.GroupName),
		producer.WithNamespace(conf.Rocketmq.NameSpace),
	)

	if err != nil {
		l.Fatalf("init producer error: %v", err)
	}

	err = p.Start()
	if err != nil {
		l.Fatalf("start producer error: %v", err)
	}
	return &MqPro{
		producer: p,
	}
}

func NewElasticsearch(conf *conf.Data) *ElasticSearch {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "creation/data/elastic-search"))
	cfg := elasticsearch.Config{
		Username: conf.ElasticSearch.User,
		Password: conf.ElasticSearch.Password,
		Addresses: []string{
			conf.ElasticSearch.Endpoint,
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		l.Fatalf("Error creating the es client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		l.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		l.Fatalf("Error: %s", res.String())
	}

	return &ElasticSearch{
		es: es,
	}
}

func NewNewsClient(conf *conf.Data) *NewsClient {
	return &NewsClient{
		url: conf.News.Url,
	}
}

func NewData(db *gorm.DB, redisCmd redis.Cmdable, cos *cos.Client, es *ElasticSearch, mq *MqPro, news *NewsClient, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "creation/data/new-data"))

	d := &Data{
		db:            db,
		log:           log.NewHelper(log.With(logger, "module", "creation/data")),
		cosCli:        cos,
		redisCli:      redisCmd,
		mqPro:         mq,
		elasticSearch: es,
		newsCli:       news,
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

		err = d.mqPro.producer.Shutdown()
		if err != nil {
			l.Errorf("shutdown mq producer error: %v", err.Error())
		}
	}, nil
}
