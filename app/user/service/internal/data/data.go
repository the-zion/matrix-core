package data

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"github.com/the-zion/matrix-core/app/user/service/internal/biz"
	"github.com/the-zion/matrix-core/app/user/service/internal/conf"
	"gopkg.in/gomail.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var ProviderSet = wire.NewSet(NewData, NewDB, NewTransaction, NewRedis, NewRocketmqProducer, NewPhoneCode, NewGoMail, NewUserRepo, NewAuthRepo, NewProfileRepo)

type TxCode struct {
	client  *sms.Client
	request *sms.SendSmsRequest
}

type GoMail struct {
	message *gomail.Message
	dialer  *gomail.Dialer
}

type Data struct {
	db           *gorm.DB
	redisCli     redis.Cmdable
	mqPro        rocketmq.Producer
	phoneCodeCli *TxCode
	goMailCli    *GoMail
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

func NewPhoneCode(conf *conf.Data) *TxCode {
	credential := common.NewCredential(
		conf.Code.SecretId,
		conf.Code.SecretKey,
	)
	cpf := profile.NewClientProfile()
	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)
	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr("1400590793")
	request.SignName = common.StringPtr("魔方技术")
	return &TxCode{
		client:  client,
		request: request,
	}
}

func NewGoMail(conf *conf.Data) *GoMail {
	m := gomail.NewMessage()
	m.SetHeader("From", "matrixtechnology@163.com")
	d := gomail.NewDialer("smtp.163.com", 465, "matrixtechnology@163.com", conf.Mail.Code)
	return &GoMail{
		message: m,
		dialer:  d,
	}
}

func NewRocketmqProducer(conf *conf.Data, logger log.Logger) rocketmq.Producer {
	l := log.NewHelper(log.With(logger, "module", "user/data/rocketmq-producer"))
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
	return p
}

func NewData(db *gorm.DB, redisCmd redis.Cmdable, p rocketmq.Producer, phoneCodeCli *TxCode, goMailCli *GoMail, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "user/data/new-data"))

	d := &Data{
		db:           db,
		mqPro:        p,
		redisCli:     redisCmd,
		phoneCodeCli: phoneCodeCli,
		goMailCli:    goMailCli,
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
