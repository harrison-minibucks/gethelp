package client

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/harrison-minibucks/gethelp/internal/conf"
)

func NewEthClient(conf *conf.Config) *ethclient.Client {
	cl, err := ethclient.Dial("\\\\.\\pipe\\geth.ipc")
	if err != nil {
		panic(err)
	}
	_ = cl
	return cl
}
