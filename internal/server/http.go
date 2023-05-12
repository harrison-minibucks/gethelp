package server

import (
	greeterV1 "github.com/harrison-minibucks/gethelp/api/helloworld/v1"
	walletV1 "github.com/harrison-minibucks/gethelp/api/wallet/v1"
	"github.com/harrison-minibucks/gethelp/internal/conf"
	"github.com/harrison-minibucks/gethelp/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, wallet *service.WalletService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	greeterV1.RegisterGreeterHTTPServer(srv, greeter)
	walletV1.RegisterWalletHTTPServer(srv, wallet)
	return srv
}
