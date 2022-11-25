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
	"github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"github.com/the-zion/matrix-core/app/user/service/internal/biz"
	"github.com/the-zion/matrix-core/app/user/service/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"runtime"
	"time"
)

var ProviderSet = wire.NewSet(NewData, NewDB, NewTransaction, NewRedis, NewRocketmqCodeProducer, NewRocketmqProfileProducer, NewRocketmqFollowProducer, NewRocketmqPictureProducer, NewRocketmqAchievementProducer, NewCosClient, NewUserRepo, NewAuthRepo, NewElasticsearch, NewRecovery)

type Cos struct {
	client *sts.Client
	opt    *sts.CredentialOptions
}

type CodeMqPro struct {
	producer rocketmq.Producer
}

type ProfileMqPro struct {
	producer rocketmq.Producer
}

type FollowMqPro struct {
	producer rocketmq.Producer
}

type PictureMqPro struct {
	producer rocketmq.Producer
}

type AchievementMqPro struct {
	producer rocketmq.Producer
}

type ElasticSearch struct {
	es *elasticsearch.Client
}

type Data struct {
	log              *log.Helper
	db               *gorm.DB
	redisCli         redis.Cmdable
	codeMqPro        *CodeMqPro
	profileMqPro     *ProfileMqPro
	followMqPro      *FollowMqPro
	pictureMqPro     *PictureMqPro
	achievementMqPro *AchievementMqPro
	elasticSearch    *ElasticSearch
	cos              *Cos
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

func NewDB(conf *conf.Data) *gorm.DB {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "user/data/mysql"))

	db, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		l.Fatalf("failed opening connection to db: %v", err)
	}
	return db
}

func NewRedis(conf *conf.Data) redis.Cmdable {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "user/data/redis"))
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Addr,
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

func NewRocketmqCodeProducer(conf *conf.Data) *CodeMqPro {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "user/data/rocketmq-code-producer"))
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.Rocketmq.ServerAddress})),
		producer.WithCredentials(primitive.Credentials{
			SecretKey: conf.Rocketmq.SecretKey,
			AccessKey: conf.Rocketmq.AccessKey,
		}),
		producer.WithInstanceName("user"),
		producer.WithGroupName(conf.Rocketmq.Code.GroupName),
		producer.WithNamespace(conf.Rocketmq.NameSpace),
	)
	if err != nil {
		l.Fatalf("init producer error: %v", err)
	}

	err = p.Start()
	if err != nil {
		l.Fatalf("start producer error: %v", err)
	}
	return &CodeMqPro{
		producer: p,
	}
}

func NewRocketmqProfileProducer(conf *conf.Data) *ProfileMqPro {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "user/data/rocketmq-profile-producer"))
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.Rocketmq.ServerAddress})),
		producer.WithCredentials(primitive.Credentials{
			SecretKey: conf.Rocketmq.SecretKey,
			AccessKey: conf.Rocketmq.AccessKey,
		}),
		producer.WithInstanceName("user"),
		producer.WithGroupName(conf.Rocketmq.Profile.GroupName),
		producer.WithNamespace(conf.Rocketmq.NameSpace),
	)

	if err != nil {
		l.Fatalf("init producer error: %v", err)
	}

	err = p.Start()
	if err != nil {
		l.Fatalf("start producer error: %v", err)
	}
	return &ProfileMqPro{
		producer: p,
	}
}

func NewRocketmqFollowProducer(conf *conf.Data) *FollowMqPro {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "user/data/rocketmq-follow-producer"))
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.Rocketmq.ServerAddress})),
		producer.WithCredentials(primitive.Credentials{
			SecretKey: conf.Rocketmq.SecretKey,
			AccessKey: conf.Rocketmq.AccessKey,
		}),
		producer.WithInstanceName("user"),
		producer.WithGroupName(conf.Rocketmq.Follow.GroupName),
		producer.WithNamespace(conf.Rocketmq.NameSpace),
	)

	if err != nil {
		l.Fatalf("init follow error: %v", err)
	}

	err = p.Start()
	if err != nil {
		l.Fatalf("start follow error: %v", err)
	}
	return &FollowMqPro{
		producer: p,
	}
}

func NewRocketmqPictureProducer(conf *conf.Data) *PictureMqPro {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "user/data/rocketmq-picture-producer"))
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.Rocketmq.ServerAddress})),
		producer.WithCredentials(primitive.Credentials{
			SecretKey: conf.Rocketmq.SecretKey,
			AccessKey: conf.Rocketmq.AccessKey,
		}),
		producer.WithInstanceName("user"),
		producer.WithGroupName(conf.Rocketmq.Picture.GroupName),
		producer.WithNamespace(conf.Rocketmq.NameSpace),
	)

	if err != nil {
		l.Fatalf("init picture error: %v", err)
	}

	err = p.Start()
	if err != nil {
		l.Fatalf("start picture error: %v", err)
	}
	return &PictureMqPro{
		producer: p,
	}
}

func NewRocketmqAchievementProducer(conf *conf.Data) *AchievementMqPro {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "creation/data/rocketmq-achievement-producer"))
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

func NewCosClient(conf *conf.Data) *Cos {
	c := sts.NewClient(
		conf.Cos.SecretId,
		conf.Cos.SecretKey,
		nil,
	)
	opt := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          conf.Cos.Region,
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{},
		},
	}
	for _, item := range conf.Cos.Policy.Statement {
		opt.Policy.Statement = append(opt.Policy.Statement, sts.CredentialPolicyStatement{
			Action:   item.Action,
			Effect:   item.Effect,
			Resource: item.Resource,
		})
	}
	return &Cos{
		client: c,
		opt:    opt,
	}
}

func NewElasticsearch(conf *conf.Data) *ElasticSearch {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "user/data/elastic-search"))
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

func NewData(db *gorm.DB, redisCmd redis.Cmdable, cp *CodeMqPro, es *ElasticSearch, pp *ProfileMqPro, fp *FollowMqPro, pip *PictureMqPro, aq *AchievementMqPro, cos *Cos, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "user/data/new-data"))

	d := &Data{
		log:              log.NewHelper(log.With(logger, "module", "creation/data")),
		db:               db,
		codeMqPro:        cp,
		profileMqPro:     pp,
		achievementMqPro: aq,
		followMqPro:      fp,
		pictureMqPro:     pip,
		redisCli:         redisCmd,
		elasticSearch:    es,
		cos:              cos,
	}
	return d, func() {
		var err error
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

		err = d.codeMqPro.producer.Shutdown()
		if err != nil {
			l.Errorf("shutdown code producer error: %v", err.Error())
		}

		err = d.profileMqPro.producer.Shutdown()
		if err != nil {
			l.Errorf("shutdown profile producer error: %v", err.Error())
		}

		err = d.followMqPro.producer.Shutdown()
		if err != nil {
			l.Errorf("shutdown follow producer error: %v", err.Error())
		}

		err = d.pictureMqPro.producer.Shutdown()
		if err != nil {
			l.Errorf("shutdown picture producer error: %v", err.Error())
		}

		err = d.achievementMqPro.producer.Shutdown()
		if err != nil {
			l.Errorf("shutdown achievement producer error: %v", err.Error())
		}
	}, nil
}
