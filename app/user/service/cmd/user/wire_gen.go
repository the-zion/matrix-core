// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/the-zion/matrix-core/app/user/service/internal/biz"
	"github.com/the-zion/matrix-core/app/user/service/internal/conf"
	"github.com/the-zion/matrix-core/app/user/service/internal/data"
	"github.com/the-zion/matrix-core/app/user/service/internal/server"
	"github.com/the-zion/matrix-core/app/user/service/internal/service"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, auth *conf.Auth, logLogger log.Logger, registry *nacos.Registry) (*kratos.App, func(), error) {
	db := data.NewDB(confData)
	cmdable := data.NewRedis(confData)
	mqPro := data.NewRocketmqProducer(confData)
	elasticSearch := data.NewElasticsearch(confData)
	cos := data.NewCosClient(confData)
	client := data.NewCosServiceClient(confData)
	github := data.NewGithub(confData)
	wechat := data.NewWechat(confData)
	qq := data.NewQQ(confData)
	gitee := data.NewGitee(confData)
	aliCode := data.NewPhoneCodeClient(confData)
	mail := data.NewMail(confData)
	dataData, cleanup2, err := data.NewData(db, cmdable, mqPro, elasticSearch, cos, client, github, wechat, qq, gitee, aliCode, mail, logLogger)
	if err != nil {
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData, logLogger)
	recovery := data.NewRecovery(dataData)
	transaction := data.NewTransaction(dataData)
	userUseCase := biz.NewUserUseCase(userRepo, recovery, transaction, logLogger)
	authRepo := data.NewAuthRepo(dataData, logLogger)
	authUseCase := biz.NewAuthUseCase(auth, authRepo, recovery, userRepo, transaction, logLogger)
	userService := service.NewUserService(userUseCase, authUseCase, logLogger)
	httpServer := server.NewHTTPServer(confServer, userService, logLogger)
	grpcServer := server.NewGRPCServer(confServer, userService, logLogger)
	kratosApp := newApp(registry, httpServer, grpcServer)
	return kratosApp, func() {
		cleanup2()
	}, nil
}
