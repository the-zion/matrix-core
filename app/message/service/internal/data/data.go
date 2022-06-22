package data

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"github.com/tencentyun/cos-go-sdk-v5"
	_ "github.com/tencentyun/cos-go-sdk-v5"
	userv1 "github.com/the-zion/matrix-core/api/user/service/v1"
	"github.com/the-zion/matrix-core/app/message/service/internal/conf"
	"gopkg.in/gomail.v2"
	"net/http"
	"net/url"
)

var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewPhoneCode, NewGoMail, NewUserServiceClient, NewCosServiceClient)

type TxCode struct {
	client  *sms.Client
	request *sms.SendSmsRequest
}

type GoMail struct {
	message *gomail.Message
	dialer  *gomail.Dialer
}

type Data struct {
	log          *log.Helper
	uc           userv1.UserClient
	phoneCodeCli *TxCode
	goMailCli    *GoMail
	cosCli       *cos.Client
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

func NewUserServiceClient(r *nacos.Registry, logger log.Logger) userv1.UserClient {
	l := log.NewHelper(log.With(logger, "module", "message/data/new-user-client"))
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///matrix.user.service.grpc"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			//tracing.Client(tracing.WithTracerProvider(tp)),
			recovery.Recovery(),
		),
	)
	if err != nil {
		l.Fatalf(err.Error())
	}
	c := userv1.NewUserClient(conn)
	return c
}

func NewCosServiceClient(conf *conf.Data, logger log.Logger) *cos.Client {
	l := log.NewHelper(log.With(logger, "module", "message/data/new-cos-client"))
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

func NewData(logger log.Logger, uc userv1.UserClient, cos *cos.Client, phoneCodeCli *TxCode, goMailCli *GoMail) (*Data, error) {
	l := log.NewHelper(log.With(logger, "module", "message/data"))
	d := &Data{
		log:          l,
		uc:           uc,
		phoneCodeCli: phoneCodeCli,
		goMailCli:    goMailCli,
		cosCli:       cos,
	}
	return d, nil
}
