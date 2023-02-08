package main

import (
	"context"
	"flag"
	nc "github.com/go-kratos/kratos/contrib/config/nacos/v2"
	"github.com/go-kratos/kratos/contrib/log/tencent/v2"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/the-zion/matrix-core/app/user/service/internal/conf"
	"github.com/the-zion/matrix-core/pkg/kube"
	"github.com/the-zion/matrix-core/pkg/trace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf       string
	traceUrl       string
	traceToken     string
	nacosUrl       string
	nacosNameSpace string
	nacosGroup     string
	nacosConfig    string
	nacosUserName  string
	nacosPassword  string
	logSelect      string
	traceProvider  *tracesdk.TracerProvider
	bootstrap      conf.Bootstrap
	servieConfig   config.Config
	rclient        *nacos.Registry
	logger         log.Logger
	tencentLogger  tencent.Logger
	app            *kratos.App
	sc             []constant.ServerConfig
	cc             *constant.ClientConfig
	cleanup        func()
	Name           = "matrix.user.service"
	id, _          = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&traceUrl, "trace", "127.0.0.1:14268/api/traces", "trace report path, eg: -trace 127.0.0.1:14268/api/traces")
	flag.StringVar(&traceToken, "token", "", "trace token, eg: -token xxx")
	flag.StringVar(&nacosUrl, "nacos", "127.0.0.1:30000", "nacos path, eg: -nacos 127.0.0.1:30000")
	flag.StringVar(&nacosNameSpace, "namespace", "public", "nacos namespace, eg: -namespace xxx")
	flag.StringVar(&nacosGroup, "group", "DEFAULT_GROUP", "nacos config group, eg: -group xxx")
	flag.StringVar(&nacosConfig, "config", "matrix.user.service", "nacos config name, eg: -config xxx")
	flag.StringVar(&nacosUserName, "username", "nacos", "nacos username, eg: -username nacos")
	flag.StringVar(&nacosPassword, "password", "nacos", "nacos password, eg: -password nacos")
	flag.StringVar(&logSelect, "log", "default", "log select, eg: -log default")
}

func newApp(r *nacos.Registry, hs *http.Server, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Registrar(r),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Server(
			hs,
			gs,
		),
	)
}

func traceInit() {
	var err error
	traceProvider, err = trace.SetTracerProvider(traceUrl, traceToken, Name, id)
	if err != nil {
		log.Error(err)
	}
}

func nacosInit() {
	url := strings.Split(nacosUrl, ":")[0]
	port, err := strconv.Atoi(strings.Split(nacosUrl, ":")[1])
	if err != nil {
		panic(err)
	}

	sc = []constant.ServerConfig{
		*constant.NewServerConfig(url, uint64(port)),
	}
	cc = &constant.ClientConfig{
		NamespaceId:          nacosNameSpace,
		Username:             nacosUserName,
		Password:             nacosPassword,
		TimeoutMs:            5000,
		NotLoadCacheAtStart:  true,
		UpdateCacheWhenEmpty: true,
		RotateTime:           "1h",
		MaxAge:               3,
		LogLevel:             "warn",
	}
	configNew()
	clientNew()
}

func configNew() {
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	dataID := nacosConfig
	group := nacosGroup

	servieConfig = config.New(
		config.WithSource(
			nc.NewConfigSource(client, nc.WithGroup(group), nc.WithDataID(dataID)),
		),
		config.WithDecoder(func(kv *config.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}),
	)

	if err = servieConfig.Load(); err != nil {
		panic(err)
	}

	if err = servieConfig.Scan(&bootstrap); err != nil {
		panic(err)
	}

	if err = servieConfig.Watch("config", func(s string, value config.Value) {
		kubeClient, err := kube.NewKubeClient()
		if err != nil {
			log.Error(err)
			return
		}
		err = kubeClient.Update("matrix", "user")
		if err != nil {
			log.Error(err)
			return
		}
	}); err != nil {
		panic(err)
	}
}

func clientNew() {
	r, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	rclient = nacos.New(r)
}

func loggerInit() {
	var err error
	logKv := []interface{}{
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	}
	switch logSelect {
	case "tencent":
		tencentLogger, err = tencent.NewLogger(
			tencent.WithEndpoint(bootstrap.Config.Log.Host),
			tencent.WithAccessKey(bootstrap.Config.Log.AccessKeyID),
			tencent.WithAccessSecret(bootstrap.Config.Log.AccessKeySecret),
			tencent.WithTopicID(bootstrap.Config.Log.TopicID),
		)
		if err != nil {
			panic(err)
		}

		tencentLogger.GetProducer().Start()
		logger = log.With(tencentLogger, logKv...)
	default:
		logger = log.With(log.NewStdLogger(os.Stdout), logKv...)
	}
}

func appInit() {
	var err error
	app, cleanup, err = wireApp(bootstrap.Config.Server, bootstrap.Config.Data, bootstrap.Config.Auth, logger, rclient)
	if err != nil {
		panic(err)
	}
}

func appRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func appClean() {
	traceProvider.Shutdown(context.Background())
	cleanup()
	if tencentLogger != nil {
		tencentLogger.Close()
	}
	servieConfig.Close()
}

func matrixRun() {
	appInit()
	appRun()
	appClean()
}

func main() {
	flag.Parse()
	nacosInit()
	traceInit()
	loggerInit()
	matrixRun()
}
