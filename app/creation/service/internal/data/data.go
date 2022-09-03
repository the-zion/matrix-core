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

var ProviderSet = wire.NewSet(NewData, NewDB, NewRedis, NewTransaction, NewRocketmqArticleProducer, NewRocketmqArticleReviewProducer, NewRocketmqAchievementProducer, NewRocketmqTalkReviewProducer, NewRocketmqTalkProducer, NewRocketmqColumnReviewProducer, NewRocketmqColumnProducer, NewRocketmqCollectionsReviewProducer, NewRocketmqCollectionsProducer, NewCosServiceClient, NewElasticsearch, NewNewsClient, NewArticleRepo, NewTalkRepo, NewCreationRepo, NewColumnRepo, NewNewsRepo, NewRecovery)

type ArticleReviewMqPro struct {
	producer rocketmq.Producer
}

type ArticleMqPro struct {
	producer rocketmq.Producer
}

type TalkReviewMqPro struct {
	producer rocketmq.Producer
}

type TalkMqPro struct {
	producer rocketmq.Producer
}

type ColumnReviewMqPro struct {
	producer rocketmq.Producer
}

type ColumnMqPro struct {
	producer rocketmq.Producer
}

type CollectionsReviewMqPro struct {
	producer rocketmq.Producer
}

type CollectionsMqPro struct {
	producer rocketmq.Producer
}

type AchievementMqPro struct {
	producer rocketmq.Producer
}

type ElasticSearch struct {
	es *elasticsearch.Client
}

type News struct {
	url string
}

type Data struct {
	db                     *gorm.DB
	log                    *log.Helper
	redisCli               redis.Cmdable
	articleMqPro           *ArticleMqPro
	articleReviewMqPro     *ArticleReviewMqPro
	talkMqPro              *TalkMqPro
	talkReviewMqPro        *TalkReviewMqPro
	columnReviewMqPro      *ColumnReviewMqPro
	columnMqPro            *ColumnMqPro
	collectionsReviewMqPro *CollectionsReviewMqPro
	collectionsMqPro       *CollectionsMqPro
	achievementMqPro       *AchievementMqPro
	cosCli                 *cos.Client
	elasticSearch          *ElasticSearch
	newsCli                *News
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
				log.Context(ctx).Errorf("%v: %s\n", rerr, buf)
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
				log.Context(ctx).Errorf("%v: %s\n", rerr, buf)
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

func NewRocketmqArticleReviewProducer(conf *conf.Data, logger log.Logger) *ArticleReviewMqPro {
	l := log.NewHelper(log.With(logger, "module", "creation/data/rocketmq-article-review-producer"))
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.CreationMq.ServerAddress})),
		producer.WithCredentials(primitive.Credentials{
			SecretKey: conf.CreationMq.SecretKey,
			AccessKey: conf.CreationMq.AccessKey,
		}),
		producer.WithInstanceName("creation"),
		producer.WithGroupName(conf.CreationMq.ArticleReview.GroupName),
		producer.WithNamespace(conf.CreationMq.NameSpace),
	)

	if err != nil {
		l.Fatalf("init producer error: %v", err)
	}

	err = p.Start()
	if err != nil {
		l.Fatalf("start producer error: %v", err)
	}
	return &ArticleReviewMqPro{
		producer: p,
	}
}

func NewRocketmqArticleProducer(conf *conf.Data, logger log.Logger) *ArticleMqPro {
	l := log.NewHelper(log.With(logger, "module", "creation/data/rocketmq-article-producer"))
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.CreationMq.ServerAddress})),
		producer.WithCredentials(primitive.Credentials{
			SecretKey: conf.CreationMq.SecretKey,
			AccessKey: conf.CreationMq.AccessKey,
		}),
		producer.WithInstanceName("creation"),
		producer.WithGroupName(conf.CreationMq.Article.GroupName),
		producer.WithNamespace(conf.CreationMq.NameSpace),
	)

	if err != nil {
		l.Fatalf("init producer error: %v", err)
	}

	err = p.Start()
	if err != nil {
		l.Fatalf("start producer error: %v", err)
	}
	return &ArticleMqPro{
		producer: p,
	}
}

func NewRocketmqTalkReviewProducer(conf *conf.Data, logger log.Logger) *TalkReviewMqPro {
	l := log.NewHelper(log.With(logger, "module", "creation/data/rocketmq-talk-review-producer"))
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.CreationMq.ServerAddress})),
		producer.WithCredentials(primitive.Credentials{
			SecretKey: conf.CreationMq.SecretKey,
			AccessKey: conf.CreationMq.AccessKey,
		}),
		producer.WithInstanceName("creation"),
		producer.WithGroupName(conf.CreationMq.TalkReview.GroupName),
		producer.WithNamespace(conf.CreationMq.NameSpace),
	)

	if err != nil {
		l.Fatalf("init producer error: %v", err)
	}

	err = p.Start()
	if err != nil {
		l.Fatalf("start producer error: %v", err)
	}
	return &TalkReviewMqPro{
		producer: p,
	}
}

func NewRocketmqTalkProducer(conf *conf.Data, logger log.Logger) *TalkMqPro {
	l := log.NewHelper(log.With(logger, "module", "creation/data/rocketmq-talk-producer"))
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.CreationMq.ServerAddress})),
		producer.WithCredentials(primitive.Credentials{
			SecretKey: conf.CreationMq.SecretKey,
			AccessKey: conf.CreationMq.AccessKey,
		}),
		producer.WithInstanceName("creation"),
		producer.WithGroupName(conf.CreationMq.Talk.GroupName),
		producer.WithNamespace(conf.CreationMq.NameSpace),
	)

	if err != nil {
		l.Fatalf("init producer error: %v", err)
	}

	err = p.Start()
	if err != nil {
		l.Fatalf("start producer error: %v", err)
	}
	return &TalkMqPro{
		producer: p,
	}
}

func NewRocketmqColumnReviewProducer(conf *conf.Data, logger log.Logger) *ColumnReviewMqPro {
	l := log.NewHelper(log.With(logger, "module", "creation/data/rocketmq-column-review-producer"))
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.CreationMq.ServerAddress})),
		producer.WithCredentials(primitive.Credentials{
			SecretKey: conf.CreationMq.SecretKey,
			AccessKey: conf.CreationMq.AccessKey,
		}),
		producer.WithInstanceName("creation"),
		producer.WithGroupName(conf.CreationMq.ColumnReview.GroupName),
		producer.WithNamespace(conf.CreationMq.NameSpace),
	)

	if err != nil {
		l.Fatalf("init producer error: %v", err)
	}

	err = p.Start()
	if err != nil {
		l.Fatalf("start producer error: %v", err)
	}
	return &ColumnReviewMqPro{
		producer: p,
	}
}

func NewRocketmqColumnProducer(conf *conf.Data, logger log.Logger) *ColumnMqPro {
	l := log.NewHelper(log.With(logger, "module", "creation/data/rocketmq-column-producer"))
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.CreationMq.ServerAddress})),
		producer.WithCredentials(primitive.Credentials{
			SecretKey: conf.CreationMq.SecretKey,
			AccessKey: conf.CreationMq.AccessKey,
		}),
		producer.WithInstanceName("creation"),
		producer.WithGroupName(conf.CreationMq.Column.GroupName),
		producer.WithNamespace(conf.CreationMq.NameSpace),
	)

	if err != nil {
		l.Fatalf("init producer error: %v", err)
	}

	err = p.Start()
	if err != nil {
		l.Fatalf("start producer error: %v", err)
	}
	return &ColumnMqPro{
		producer: p,
	}
}

func NewRocketmqCollectionsReviewProducer(conf *conf.Data, logger log.Logger) *CollectionsReviewMqPro {
	l := log.NewHelper(log.With(logger, "module", "creation/data/rocketmq-collections-review-producer"))
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.CreationMq.ServerAddress})),
		producer.WithCredentials(primitive.Credentials{
			SecretKey: conf.CreationMq.SecretKey,
			AccessKey: conf.CreationMq.AccessKey,
		}),
		producer.WithInstanceName("creation"),
		producer.WithGroupName(conf.CreationMq.CollectionsReview.GroupName),
		producer.WithNamespace(conf.CreationMq.NameSpace),
	)

	if err != nil {
		l.Fatalf("init producer error: %v", err)
	}

	err = p.Start()
	if err != nil {
		l.Fatalf("start producer error: %v", err)
	}
	return &CollectionsReviewMqPro{
		producer: p,
	}
}

func NewRocketmqCollectionsProducer(conf *conf.Data, logger log.Logger) *CollectionsMqPro {
	l := log.NewHelper(log.With(logger, "module", "creation/data/rocketmq-collections-producer"))
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.CreationMq.ServerAddress})),
		producer.WithCredentials(primitive.Credentials{
			SecretKey: conf.CreationMq.SecretKey,
			AccessKey: conf.CreationMq.AccessKey,
		}),
		producer.WithInstanceName("creation"),
		producer.WithGroupName(conf.CreationMq.Collections.GroupName),
		producer.WithNamespace(conf.CreationMq.NameSpace),
	)

	if err != nil {
		l.Fatalf("init producer error: %v", err)
	}

	err = p.Start()
	if err != nil {
		l.Fatalf("start producer error: %v", err)
	}
	return &CollectionsMqPro{
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

func NewElasticsearch(conf *conf.Data, logger log.Logger) *ElasticSearch {
	l := log.NewHelper(log.With(logger, "module", "creation/data/elastic-search"))
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

func NewNewsClient(conf *conf.Data) *News {
	return &News{
		url: conf.News.Url,
	}
}

func NewData(db *gorm.DB, redisCmd redis.Cmdable, cos *cos.Client, es *ElasticSearch, amp *ArticleMqPro, arp *ArticleReviewMqPro, tmp *TalkMqPro, trp *TalkReviewMqPro, cmp *ColumnMqPro, crq *ColumnReviewMqPro, cormq *CollectionsReviewMqPro, cmq *CollectionsMqPro, ap *AchievementMqPro, news *News, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "creation/data/new-data"))

	d := &Data{
		db:                     db,
		cosCli:                 cos,
		redisCli:               redisCmd,
		articleMqPro:           amp,
		articleReviewMqPro:     arp,
		talkMqPro:              tmp,
		talkReviewMqPro:        trp,
		columnMqPro:            cmp,
		columnReviewMqPro:      crq,
		collectionsReviewMqPro: cormq,
		collectionsMqPro:       cmq,
		achievementMqPro:       ap,
		elasticSearch:          es,
		newsCli:                news,
	}
	return d, func() {
		l.Info("closing the data resources")

		err := d.articleMqPro.producer.Shutdown()
		if err != nil {
			l.Errorf("shutdown article producer error: %v", err.Error())
		}

		err = d.articleReviewMqPro.producer.Shutdown()
		if err != nil {
			l.Errorf("shutdown article review producer error: %v", err.Error())
		}

		err = d.talkMqPro.producer.Shutdown()
		if err != nil {
			l.Errorf("shutdown talk producer error: %v", err.Error())
		}

		err = d.talkReviewMqPro.producer.Shutdown()
		if err != nil {
			l.Errorf("shutdown talk review producer error: %v", err.Error())
		}

		err = d.columnMqPro.producer.Shutdown()
		if err != nil {
			l.Errorf("shutdown column producer error: %v", err.Error())
		}

		err = d.columnReviewMqPro.producer.Shutdown()
		if err != nil {
			l.Errorf("shutdown column review producer error: %v", err.Error())
		}

		err = d.collectionsMqPro.producer.Shutdown()
		if err != nil {
			l.Errorf("shutdown collections producer error: %v", err.Error())
		}

		err = d.collectionsReviewMqPro.producer.Shutdown()
		if err != nil {
			l.Errorf("shutdown collections review producer error: %v", err.Error())
		}

		err = d.achievementMqPro.producer.Shutdown()
		if err != nil {
			l.Errorf("shutdown achievement producer error: %v", err.Error())
		}
	}, nil
}
