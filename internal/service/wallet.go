package service

import (
	"context"

	v1 "github.com/harrison-minibucks/gethelp/api/wallet/v1"
	"github.com/harrison-minibucks/gethelp/internal/biz"
)

type WalletService struct {
	v1.UnimplementedWalletServer
	uc *biz.WalletUsecase
}

func NewWalletService(uc *biz.WalletUsecase) *WalletService {
	return &WalletService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *WalletService) GetBalance(ctx context.Context, in *v1.BalanceRequest) (*v1.BalanceReply, error) {
	b, err := s.uc.ReadBalance(in.Account)
	if err != nil {
		return nil, err
	}
	return &v1.BalanceReply{Account: b.Account, Balance: b.Balance, PendingBalance: b.PendingBalance}, nil
}

func (s *WalletService) SendTransaction(ctx context.Context, in *v1.TxRequest) (*v1.TxReply, error) {
	err := s.uc.SendTransaction(ctx, &biz.SendTransaction{
		SenderAccount: in.SenderAccount,
		Password:      in.Password,
		Recipient:     in.RecipientAccount,
	})
	if err != nil {
		return nil, err
	}
	return &v1.TxReply{Success: true}, nil
}

func (s *WalletService) SuggestGasPrice(ctx context.Context, in *v1.Empty) (*v1.GasPrice, error) {
	gas, err := s.uc.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GasPrice{Gas: gas}, nil
}
