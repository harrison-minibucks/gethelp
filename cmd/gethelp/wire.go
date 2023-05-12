//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/harrison-minibucks/gethelp/internal/biz"
	"github.com/harrison-minibucks/gethelp/internal/client"
	"github.com/harrison-minibucks/gethelp/internal/conf"
	"github.com/harrison-minibucks/gethelp/internal/data"
	"github.com/harrison-minibucks/gethelp/internal/server"
	"github.com/harrison-minibucks/gethelp/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Config, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, client.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
