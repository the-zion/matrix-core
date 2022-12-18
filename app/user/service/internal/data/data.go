package data

import (
	"context"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"github.com/the-zion/matrix-core/app/user/service/internal/biz"
	"github.com/the-zion/matrix-core/app/user/service/internal/conf"
	"gopkg.in/gomail.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

var ProviderSet = wire.NewSet(NewData, NewDB, NewTransaction, NewRedis, NewRocketmqProducer, NewCosClient, NewCosServiceClient, NewUserRepo, NewAuthRepo, NewElasticsearch, NewGithub, NewWechat, NewQQ, NewGitee, NewPhoneCodeClient, NewMail, NewRecovery)

type Cos struct {
	client *sts.Client
	opt    *sts.CredentialOptions
}

type MqPro struct {
	producer rocketmq.Producer
}

type ElasticSearch struct {
	es *elasticsearch.Client
}

type Github struct {
	accessTokenUrl string
	userInfoUrl    string
	clientId       string
	clientSecret   string
}

type Wechat struct {
	accessTokenUrl string
	userInfoUrl    string
	appid          string
	secret         string
	grantType      string
}

type QQ struct {
	accessTokenUrl string
	openIdUrl      string
	userInfoUrl    string
	clientId       string
	clientSecret   string
	grantType      string
	redirectUri    string
}

type Gitee struct {
	accessTokenUrl string
	userInfoUrl    string
	clientId       string
	clientSecret   string
	grantType      string
	redirectUri    string
}

type AliCode struct {
	client   *dysmsapi20170525.Client
	signName string
}

type Mail struct {
	message *gomail.Message
	dialer  *gomail.Dialer
}

type Data struct {
	log           *log.Helper
	db            *gorm.DB
	redisCli      redis.Cmdable
	cosCli        *cos.Client
	mqPro         *MqPro
	elasticSearch *ElasticSearch
	cos           *Cos
	github        *Github
	gitee         *Gitee
	wechat        *Wechat
	qq            *QQ
	aliCode       *AliCode
	mail          *Mail
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

func NewRocketmqProducer(conf *conf.Data) *MqPro {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "user/data/rocketmq-producer"))
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

func NewCosServiceClient(conf *conf.Data) *cos.Client {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "user/data/new-cos-client"))
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

func NewGithub(conf *conf.Data) *Github {
	return &Github{
		accessTokenUrl: conf.Github.AccessTokenUrl,
		userInfoUrl:    conf.Github.UserInfoUrl,
		clientId:       conf.Github.ClientId,
		clientSecret:   conf.Github.ClientSecret,
	}
}

func NewWechat(conf *conf.Data) *Wechat {
	return &Wechat{
		accessTokenUrl: conf.Wechat.AccessTokenUrl,
		userInfoUrl:    conf.Wechat.UserInfoUrl,
		appid:          conf.Wechat.Appid,
		secret:         conf.Wechat.Secret,
		grantType:      conf.Wechat.GrantType,
	}
}

func NewQQ(conf *conf.Data) *QQ {
	return &QQ{
		accessTokenUrl: conf.Qq.AccessTokenUrl,
		openIdUrl:      conf.Qq.OpenIdUrl,
		userInfoUrl:    conf.Qq.UserInfoUrl,
		clientId:       conf.Qq.ClientId,
		clientSecret:   conf.Qq.ClientSecret,
		grantType:      conf.Qq.GrantType,
		redirectUri:    conf.Qq.RedirectUri,
	}
}

func NewGitee(conf *conf.Data) *Gitee {
	return &Gitee{
		accessTokenUrl: conf.Gitee.AccessTokenUrl,
		userInfoUrl:    conf.Gitee.UserInfoUrl,
		clientId:       conf.Gitee.ClientId,
		clientSecret:   conf.Gitee.ClientSecret,
		grantType:      conf.Gitee.GrantType,
		redirectUri:    conf.Gitee.RedirectUri,
	}
}

func NewPhoneCodeClient(conf *conf.Data) *AliCode {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "user/data/message"))
	config := &openapi.Config{
		AccessKeyId:     &conf.AliCode.AccessKeyId,
		AccessKeySecret: &conf.AliCode.AccessKeySecret,
	}
	config.Endpoint = tea.String(conf.AliCode.DomainUrl)
	client, err := dysmsapi20170525.NewClient(config)
	if err != nil {
		l.Fatalf("error creating the msg client: %s", err)
	}
	return &AliCode{
		client:   client,
		signName: conf.AliCode.SignName,
	}
}

func NewMail(conf *conf.Data) *Mail {
	m := gomail.NewMessage()
	m.SetHeader("From", "matrixtechnology@163.com")
	d := gomail.NewDialer("smtp.163.com", 465, "matrixtechnology@163.com", conf.Mail.Code)
	return &Mail{
		message: m,
		dialer:  d,
	}
}

func NewData(db *gorm.DB, redisCmd redis.Cmdable, mp *MqPro, es *ElasticSearch, cos *Cos, cosCli *cos.Client, github *Github, wechat *Wechat, qq *QQ, gitee *Gitee, code *AliCode, mailCli *Mail, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(log.GetLogger(), "module", "user/data/new-data"))

	d := &Data{
		log:           log.NewHelper(log.With(logger, "module", "creation/data")),
		db:            db,
		mqPro:         mp,
		redisCli:      redisCmd,
		elasticSearch: es,
		cos:           cos,
		cosCli:        cosCli,
		github:        github,
		wechat:        wechat,
		qq:            qq,
		gitee:         gitee,
		aliCode:       code,
		mail:          mailCli,
	}
	return d, func() {
		var err error
		l.Info("closing the data resources")

		mailCli.message.Reset()
		mail, err := mailCli.dialer.Dial()
		if err != nil {
			l.Errorf("close goMail err: %v", err.Error())
		}

		err = mail.Close()
		if err != nil {
			l.Errorf("close goMail err: %v", err.Error())
		}

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
