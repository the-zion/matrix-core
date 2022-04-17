package data

import (
	"context"
	"github.com/Cube-v2/cube-core/app/user/service/internal/conf"
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

var ProviderSet = wire.NewSet(NewData, NewDB, NewRedis, NewPhoneCode, NewGoMail, NewUserRepo, NewAuthRepo, NewProfileRepo)

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
	phoneCodeCli *TxCode
	goMailCli    *GoMail
}

func NewDB(conf *conf.Data, logger log.Logger) *gorm.DB {
	l := log.NewHelper(log.With(logger, "module", "user/data/new-mysql"))

	db, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		l.Fatalf("failed opening connection to db: %v", err)
	}
	if err := db.AutoMigrate(&User{}, &Profile{}, &Achievement{}, &Follower{}); err != nil {
		l.Fatalf("failed creat or update table resources: %v", err)
	}
	return db
}

func NewRedis(conf *conf.Data, logger log.Logger) redis.Cmdable {
	l := log.NewHelper(log.With(logger, "module", "user/data/new-redis"))
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

func NewData(db *gorm.DB, redisCmd redis.Cmdable, phoneCodeCli *TxCode, goMailCli *GoMail, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "user/data/new-data"))

	d := &Data{
		db:           db,
		redisCli:     redisCmd,
		phoneCodeCli: phoneCodeCli,
		goMailCli:    goMailCli,
	}
	return d, func() {
		l.Info("closing the data resources")
	}, nil
}
