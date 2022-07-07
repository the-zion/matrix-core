//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/the-zion/matrix-core/app/achievement/service/internal/biz"
	"github.com/the-zion/matrix-core/app/achievement/service/internal/conf"
	"github.com/the-zion/matrix-core/app/achievement/service/internal/data"
	"github.com/the-zion/matrix-core/app/achievement/service/internal/server"
	"github.com/the-zion/matrix-core/app/achievement/service/internal/service"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger, *nacos.Registry) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
