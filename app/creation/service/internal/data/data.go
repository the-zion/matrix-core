package data

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/tencentyun/cos-go-sdk-v5"
	_ "github.com/tencentyun/cos-go-sdk-v5"
	"github.com/the-zion/matrix-core/app/creation/service/internal/biz"
	"github.com/the-zion/matrix-core/app/creation/service/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/url"
)

var ProviderSet = wire.NewSet(NewData, NewDB, NewTransaction, NewRocketmqArticleProducer, NewCosServiceClient, NewArticleRepo)

type ArticleMqPro struct {
	producer rocketmq.Producer
}

type Data struct {
	db           *gorm.DB
	log          *log.Helper
	articleMqPro *ArticleMqPro
	cosCli       *cos.Client
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

func NewRocketmqArticleProducer(conf *conf.Data, logger log.Logger) *ArticleMqPro {
	l := log.NewHelper(log.With(logger, "module", "user/data/rocketmq-article-producer"))
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.Rocketmq.ServerAddress})),
		producer.WithCredentials(primitive.Credentials{
			SecretKey: conf.Rocketmq.SecretKey,
			AccessKey: conf.Rocketmq.AccessKey,
		}),
		producer.WithRetry(2),
		producer.WithGroupName(conf.Rocketmq.Article.GroupName),
		producer.WithNamespace(conf.Rocketmq.NameSpace),
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

func NewData(db *gorm.DB, cos *cos.Client, ap *ArticleMqPro, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "creation/data/new-data"))

	d := &Data{
		db:           db,
		cosCli:       cos,
		articleMqPro: ap,
	}
	return d, func() {
		l.Info("closing the data resources")

		err := d.articleMqPro.producer.Shutdown()
		if err != nil {
			l.Errorf("shutdown article producer error: %v", err.Error())
		}
	}, nil
}
