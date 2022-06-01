package data

import (
	"context"
	"fmt"
	"github.com/Cube-v2/matrix-core/app/user/service/internal/conf"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"gopkg.in/gomail.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var ProviderSet = wire.NewSet(NewData, NewDB, NewRedis, NewRocketmqProducer, NewRocketmqConsumer, NewPhoneCode, NewGoMail, NewUserRepo, NewAuthRepo, NewProfileRepo)

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
	mqCum        rocketmq.PushConsumer
	phoneCodeCli *TxCode
	goMailCli    *GoMail
}

func NewDB(conf *conf.Data, logger log.Logger) *gorm.DB {
	l := log.NewHelper(log.With(logger, "module", "user/data/mysql"))

	db, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		l.Fatalf("failed opening connection to db: %v", err)
	}
	if err := db.AutoMigrate(&User{}, &Profile{}); err != nil {
		l.Fatalf("failed creat or update table resources: %v", err)
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
	m.SetHeader("From", "cube-technology@foxmail.com")
	d := gomail.NewDialer("smtp.qq.com", 465, "cube-technology@foxmail.com", conf.Mail.Code)
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

func NewRocketmqConsumer(conf *conf.Data, logger log.Logger) rocketmq.PushConsumer {
	l := log.NewHelper(log.With(logger, "module", "user/data/rocketmq-producer"))
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(conf.Rocketmq.GroupName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.Rocketmq.ServerAddress})),
		consumer.WithCredentials(primitive.Credentials{
			SecretKey: conf.Rocketmq.SecretKey,
			AccessKey: conf.Rocketmq.AccessKey,
		}),
		consumer.WithNamespace(conf.Rocketmq.NameSpace),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		consumer.WithConsumerModel(consumer.Clustering),
	)
	if err != nil {
		l.Fatalf("init consumer error: %v", err)
	}

	//err = c.Subscribe("code", consumer.MessageSelector{
	//	Type:       consumer.TAG,
	//	Expression: "Phone || Email",
	//}, func(ctx context.Context,
	//	msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	//	fmt.Printf("subscribe callback len: %d \n", len(msgs))
	//	return consumer.ConsumeSuccess, nil
	//})
	err = c.Subscribe("code", consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: "Phone || Email",
	}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		fmt.Printf("subscribe callback len: %d \n", len(msgs))
		// 设置下次消费的延迟级别
		concurrentCtx, _ := primitive.GetConcurrentlyCtx(ctx)
		concurrentCtx.DelayLevelWhenNextConsume = 1 // only run when return consumer.ConsumeRetryLater

		for _, msg := range msgs {
			// 模拟重试3次后消费成功
			if msg.ReconsumeTimes > 3 {
				fmt.Printf("msg ReconsumeTimes > 3. msg: %v", msg)
				return consumer.ConsumeSuccess, nil
			} else {
				fmt.Printf("subscribe callback: %v \n", msg)
			}
		}
		// 模拟消费失败，回复重试
		return consumer.ConsumeRetryLater, nil
	})
	if err != nil {
		l.Fatalf("consumer subscribe error: %v", err)
	}

	err = c.Start()
	if err != nil {
		l.Fatalf("start consumer error: %v", err)
	}

	return c
}

func NewData(db *gorm.DB, redisCmd redis.Cmdable, p rocketmq.Producer, c rocketmq.PushConsumer, phoneCodeCli *TxCode, goMailCli *GoMail, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "user/data/new-data"))

	d := &Data{
		db:           db,
		mqPro:        p,
		mqCum:        c,
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

		err = d.mqCum.Shutdown()
		if err != nil {
			l.Errorf("shutdown consumer error: %v", err.Error())
		}
	}, nil
}
