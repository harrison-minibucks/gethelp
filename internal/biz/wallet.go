package biz

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/harrison-minibucks/gethelp/internal/util"

	"github.com/go-kratos/kratos/v2/log"
)

const MY_PUBLIC_ADDRESS = "0x929548598a3b93362c5aa2a24de190d18e657ae0"
const RECIPIENT_ADDRESS = "0xb02A2EdA1b317FBd16760128836B0Ac59B560e91"
const KEYSTORE_DIR = "%LocalAppData%\\Ethereum\\keystore" // TODO: Change with host os
const PASSWORD = "P@ssw0rd"                               // TODO: Remove

type KeystoreRepo interface {
	HasAddress(string) bool
	SignTransaction(accounts.Account, *types.Transaction, *big.Int) (*types.Transaction, error)
	SecretKeyOf(string, string) (string, error)
	Unlock(accounts.Account, string) error
}

type WalletUsecase struct {
	repo KeystoreRepo
	cl   *ethclient.Client
	log  *log.Helper
}

type AccountBalance struct {
	Account        string
	Balance        *big.Int
	PendingBalance *big.Int
}

type SendTransaction struct {
	SenderAccount string
	Password      string
	Recipient     string
	Amount        *big.Int
}

type TxCost struct {
	IsPending bool
	TxCost    string
}

type TransactionResult struct {
	TransactionHash string
}

func NewWalletUsecase(repo KeystoreRepo, logger log.Logger, cl *ethclient.Client) *WalletUsecase {
	return &WalletUsecase{repo: repo, log: log.NewHelper(logger), cl: cl}
}

// Prints the balance of an account
func (s *WalletUsecase) ReadBalance(accountAddress string) (*AccountBalance, error) {
	account := common.HexToAddress(accountAddress)
	balance, err := s.cl.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return nil, err
	}
	pendingBalance, err := s.cl.PendingBalanceAt(context.Background(), account)
	if err != nil {
		return nil, err
	}
	accBalance := &AccountBalance{
		Account: accountAddress,
		Balance: balance,
	}
	accBalance.PendingBalance = pendingBalance
	return accBalance, nil
}

func (s *WalletUsecase) SendTransaction(ctx context.Context, sendTx *SendTransaction) (*TransactionResult, error) {
	account := accounts.Account{Address: common.HexToAddress(sendTx.SenderAccount)}
	if !s.repo.HasAddress(sendTx.SenderAccount) {
		s.log.Error("Account ", sendTx.SenderAccount, " not found in keystore")
		return nil, fmt.Errorf("account not found in keystore")
	}

	err := s.repo.Unlock(account, sendTx.Password)
	if err != nil {
		s.log.Error("Failed to unlock account")
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	var (
		to       = common.HexToAddress(sendTx.Recipient)
		value    = sendTx.Amount
		gasLimit = uint64(21000)
	)

	chainID, err := s.cl.ChainID(ctx)
	if err != nil {
		s.log.Error("Fail to retrieve chainid")
		return nil, err
	}

	nonce, err := s.cl.PendingNonceAt(ctx, account.Address)
	if err != nil {
		s.log.Error("Fail to retrieve nonce")
		return nil, err
	}

	tipCap, _ := s.cl.SuggestGasTipCap(ctx)
	feeCap, _ := s.cl.SuggestGasPrice(ctx)

	tx := types.NewTx(
		&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			GasTipCap: tipCap,
			GasFeeCap: feeCap,
			Gas:       gasLimit,
			To:        &to,
			Value:     value,
			Data:      nil,
		})

	signedTx, err := s.repo.SignTransaction(account, tx, chainID)
	if err != nil {
		s.log.Error("Failed to sign transaction")
		return nil, err
	}
	if err := s.cl.SendTransaction(ctx, signedTx); err != nil {
		return nil, err
	} else {
		return &TransactionResult{
			TransactionHash: signedTx.Hash().Hex(),
		}, nil
	}
}

func (s *WalletUsecase) TxCost(ctx context.Context, txHashStr string) (*TxCost, error) {
	txHash := common.HexToHash(txHashStr)
	_, isPending, err := s.cl.TransactionByHash(ctx, txHash)
	if err != nil {
		return nil, err
	}
	if isPending {
		return &TxCost{IsPending: true}, nil
	}
	if receipt, err := s.cl.TransactionReceipt(ctx, txHash); err != nil {
		return nil, err
	} else {
		cost := new(big.Int).Mul(big.NewInt(int64(receipt.GasUsed)), receipt.EffectiveGasPrice)
		return &TxCost{TxCost: util.FormatEth(cost)}, err
	}
}

func (s *WalletUsecase) TransferAll(ctx context.Context, sendTx *SendTransaction) (*TransactionResult, error) {
	tipCap, _ := s.cl.SuggestGasTipCap(ctx)
	feeCap, _ := s.cl.SuggestGasPrice(ctx)
	gasLimit := uint64(21000)
	estimatedTxCost := new(big.Int)
	if tipCap.Cmp(feeCap) > 0 {
		estimatedTxCost.Mul(tipCap, big.NewInt(int64(gasLimit)))
	} else {
		estimatedTxCost.Mul(feeCap, big.NewInt(int64(gasLimit)))
	}
	accBalance, err := s.ReadBalance(sendTx.SenderAccount)
	if err != nil {
		return nil, err
	}
	return s.SendTransaction(ctx, &SendTransaction{
		SenderAccount: sendTx.SenderAccount,
		Password:      sendTx.Password,
		Recipient:     sendTx.Recipient,
		Amount:        new(big.Int).Sub(accBalance.PendingBalance, estimatedTxCost),
	})
}

func (s *WalletUsecase) SuggestGasPrice(ctx context.Context) (string, error) {
	gas, err := s.cl.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}
	return util.FormatEth(gas), nil
}
