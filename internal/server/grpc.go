package server

import (
	greeterV1 "github.com/harrison-minibucks/gethelp/api/helloworld/v1"
	walletV1 "github.com/harrison-minibucks/gethelp/api/wallet/v1"
	"github.com/harrison-minibucks/gethelp/internal/conf"
	"github.com/harrison-minibucks/gethelp/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, greeter *service.GreeterService, wallet *service.WalletService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	greeterV1.RegisterGreeterServer(srv, greeter)
	walletV1.RegisterWalletServer(srv, wallet)
	return srv
}
