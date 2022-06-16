// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/the-zion/matrix-core/app/bff/interface/internal/biz"
	"github.com/the-zion/matrix-core/app/bff/interface/internal/conf"
	"github.com/the-zion/matrix-core/app/bff/interface/internal/data"
	"github.com/the-zion/matrix-core/app/bff/interface/internal/server"
	"github.com/the-zion/matrix-core/app/bff/interface/internal/service"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger, registry *nacos.Registry) (*kratos.App, func(), error) {
	userClient := data.NewUserServiceClient(registry, logger)
	dataData, err := data.NewData(logger, userClient)
	if err != nil {
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData, logger)
	userUseCase := biz.NewUserUseCase(userRepo, logger)
	bffService := service.NewUserService(userUseCase, logger)
	httpServer := server.NewHTTPServer(confServer, bffService, logger)
	grpcServer := server.NewGRPCServer(confServer, bffService, logger)
	app := newApp(logger, registry, httpServer, grpcServer)
	return app, func() {
	}, nil
}