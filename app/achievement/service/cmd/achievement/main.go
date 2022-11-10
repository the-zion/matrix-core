package main

import (
	"flag"
	nc "github.com/go-kratos/kratos/contrib/config/nacos/v2"
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
	"github.com/the-zion/matrix-core/app/achievement/service/internal/conf"
	"github.com/the-zion/matrix-core/pkg/trace"
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
	Name           = "matrix.achievement.service"
	id, _          = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&traceUrl, "trace", "127.0.0.1:14268/api/traces", "trace report path, eg: -trace 127.0.0.1:14268/api/traces")
	flag.StringVar(&traceToken, "token", "", "trace token, eg: -token xxx")
	flag.StringVar(&nacosUrl, "nacos", "127.0.0.1:30000", "nacos path, eg: -nacos 127.0.0.1:30000")
	flag.StringVar(&nacosNameSpace, "namespace", "public", "nacos namespace, eg: -namespace xxx")
	flag.StringVar(&nacosGroup, "group", "DEFAULT_GROUP", "nacos config group, eg: -group xxx")
	flag.StringVar(&nacosConfig, "config", "matrix.achievement.service", "nacos config name, eg: -config xxx")
}

func newApp(logger log.Logger, r *nacos.Registry, hs *http.Server, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Registrar(r),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	err := trace.SetTracerProvider(traceUrl, traceToken, Name, id)
	if err != nil {
		log.Error(err)
	}

	url := strings.Split(nacosUrl, ":")[0]
	port, err := strconv.Atoi(strings.Split(nacosUrl, ":")[1])
	if err != nil {
		panic(err)
	}

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(url, uint64(port)),
	}
	cc := &constant.ClientConfig{
		NamespaceId:          nacosNameSpace,
		TimeoutMs:            5000,
		NotLoadCacheAtStart:  true,
		UpdateCacheWhenEmpty: true,
		RotateTime:           "1h",
		MaxAge:               3,
		LogLevel:             "warn",
	}

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
	_, err = client.PublishConfig(vo.ConfigParam{DataId: dataID, Group: group, Content: `
logger:
  level: warn
`})
	if err != nil {
		panic(err)
	}

	c := config.New(
		config.WithSource(
			nc.NewConfigSource(client, nc.WithGroup(group), nc.WithDataID(dataID)),
		),
		config.WithDecoder(func(kv *config.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	rclient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	r := nacos.New(rclient)

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger, r)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
