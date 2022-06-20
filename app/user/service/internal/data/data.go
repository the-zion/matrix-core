package data

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"github.com/the-zion/matrix-core/app/user/service/internal/biz"
	"github.com/the-zion/matrix-core/app/user/service/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var ProviderSet = wire.NewSet(NewData, NewDB, NewTransaction, NewRedis, NewRocketmqCodeProducer, NewCosClient, NewUserRepo, NewAuthRepo, NewProfileRepo)

type Cos struct {
	client *sts.Client
	opt    *sts.CredentialOptions
}

type Data struct {
	db       *gorm.DB
	redisCli redis.Cmdable
	mqPro    rocketmq.Producer
	cos      *Cos
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
	l := log.NewHelper(log.With(logger, "module", "user/data/mysql"))

	db, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		l.Fatalf("failed opening connection to db: %v", err)
	}
	return db
}

func NewRedis(conf *conf.Data, logger log.Logger) redis.Cmdable {
	l := log.NewHelper(log.With(logger, "module", "user/data/redis"))
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

func NewRocketmqCodeProducer(conf *conf.Data, logger log.Logger) rocketmq.Producer {
	l := log.NewHelper(log.With(logger, "module", "user/data/rocketmq-code-producer"))
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.Rocketmq.Code.ServerAddress})),
		producer.WithCredentials(primitive.Credentials{
			SecretKey: conf.Rocketmq.Code.SecretKey,
			AccessKey: conf.Rocketmq.Code.AccessKey,
		}),
		producer.WithGroupName(conf.Rocketmq.Code.GroupName),
		producer.WithNamespace(conf.Rocketmq.Code.NameSpace),
	)
	if err != nil {
		l.Fatalf("init producer error: %v", err)
	}

	err = p.Start()
	if err != nil {
		l.Fatalf("start producer error: %v", err)
	}
	return p
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
			Statement: []sts.CredentialPolicyStatement{
				{
					Action: []string{
						"name/cos:PostObject",
						"name/cos:PutObject",
						"name/cos:InitiateMultipartUpload",
						"name/cos:ListMultipartUploads",
						"name/cos:ListParts",
						"name/cos:UploadPart",
						"name/cos:CompleteMultipartUpload",
					},
					Effect: "allow",
					Resource: []string{
						"qcs::cos:" + conf.Cos.Region + ":uid/" + conf.Cos.Appid + ":" + conf.Cos.Bucket + "/avatar/*",
					},
				},
			},
		},
	}
	return &Cos{
		client: c,
		opt:    opt,
	}
}

func NewData(db *gorm.DB, redisCmd redis.Cmdable, p rocketmq.Producer, cos *Cos, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "user/data/new-data"))

	d := &Data{
		db:       db,
		mqPro:    p,
		redisCli: redisCmd,
		cos:      cos,
	}
	return d, func() {
		var err error
		l.Info("closing the data resources")

		err = d.mqPro.Shutdown()
		if err != nil {
			l.Errorf("shutdown producer error: %v", err.Error())
		}
	}, nil
}
