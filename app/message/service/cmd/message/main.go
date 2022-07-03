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
	"github.com/the-zion/matrix-core/app/message/service/internal/conf"
	"github.com/the-zion/matrix-core/app/message/service/internal/server"
	"gopkg.in/yaml.v3"
	"os"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	Name  = "matrix.message.service"
	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, r *nacos.Registry, hs *http.Server, gs *grpc.Server, cmcs *server.CodeMqConsumerServer, pmcs *server.ProfileMqConsumerServer, adms *server.ArticleDraftMqConsumerServer, amcs *server.ArticleMqConsumerServer) *kratos.App {
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
			cmcs,
			pmcs,
			adms,
			amcs,
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

	sc := []constant.ServerConfig{
		*constant.NewServerConfig("124.223.94.61", 30000),
	}
	cc := &constant.ClientConfig{
		NamespaceId:          "public",
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

	dataID := Name
	group := "DEFAULT_GROUP"
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
