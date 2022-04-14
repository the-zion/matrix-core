//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/Cube-v2/cube-core/app/user/service/internal/biz"
	"github.com/Cube-v2/cube-core/app/user/service/internal/conf"
	"github.com/Cube-v2/cube-core/app/user/service/internal/data"
	"github.com/Cube-v2/cube-core/app/user/service/internal/server"
	"github.com/Cube-v2/cube-core/app/user/service/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Auth, *validator.Validate, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
