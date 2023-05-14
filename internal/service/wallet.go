package service

import (
	"context"

	v1 "github.com/harrison-minibucks/gethelp/api/wallet/v1"
	"github.com/harrison-minibucks/gethelp/internal/biz"
	"github.com/harrison-minibucks/gethelp/internal/util"
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
	return &v1.BalanceReply{
		Account:        b.Account,
		Balance:        util.FormatEth(b.Balance),
		PendingBalance: util.FormatEth(b.PendingBalance),
	}, nil
}

func (s *WalletService) SendTransaction(ctx context.Context, in *v1.TxRequest) (*v1.TxReply, error) {
	amount, err := util.ConvertAmount(in.Amount)
	if err != nil {
		return nil, err
	}
	res, err := s.uc.SendTransaction(ctx, &biz.SendTransaction{
		SenderAccount: in.SenderAccount,
		Password:      in.Password,
		Recipient:     in.RecipientAccount,
		Amount:        amount,
	})
	if err != nil {
		return nil, err
	}
	return &v1.TxReply{
		Success:         true,
		TransactionHash: res.TransactionHash,
	}, nil
}

func (s *WalletService) SuggestGasPrice(ctx context.Context, in *v1.Empty) (*v1.GasPrice, error) {
	gas, err := s.uc.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GasPrice{Gas: gas}, nil
}

func (s *WalletService) TxCost(ctx context.Context, in *v1.TxCostRequest) (*v1.TxCostReply, error) {
	res, err := s.uc.TxCost(ctx, in.TxHash)
	if err != nil {
		return nil, err
	}
	return &v1.TxCostReply{IsPending: res.IsPending, TxCost: res.TxCost}, nil
}

func (s *WalletService) WithdrawWallet(ctx context.Context, in *v1.Withdrawal) (*v1.WithdrawalResult, error) {
	res, err := s.uc.TransferAll(ctx, &biz.SendTransaction{
		SenderAccount: in.WalletAddress,
		Password:      in.Password,
		Recipient:     in.RecipientAccount,
	})
	if err != nil {
		return nil, err
	}
	return &v1.WithdrawalResult{Success: true, TransactionHash: res.TransactionHash}, nil
}

func (s *WalletService) DepositWallet(ctx context.Context, in *v1.Deposit) (*v1.DepositResult, error) {
	res, err := s.uc.TransferAll(ctx, &biz.SendTransaction{
		SenderAccount: in.SenderAccount,
		Password:      in.Password,
		Recipient:     in.WalletAddress,
	})
	if err != nil {
		return nil, err
	}
	return &v1.DepositResult{Success: true, TransactionHash: res.TransactionHash}, nil
}
