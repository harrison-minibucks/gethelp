package biz

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"

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
	Balance        string
	PendingBalance string
}

type SendTransaction struct {
	SenderAccount string
	Password      string
	Recipient     string
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
		Balance: formatEth(balance),
	}
	if balance.Cmp(pendingBalance) != 0 {
		accBalance.PendingBalance = formatEth(pendingBalance)
	}
	return accBalance, nil
}

func (s *WalletUsecase) SendTransaction(sendTx *SendTransaction) error {
	account := accounts.Account{Address: common.HexToAddress(sendTx.SenderAccount)}
	if !s.repo.HasAddress(sendTx.SenderAccount) {
		fmt.Println("Account", sendTx.SenderAccount, "not found in keystore")
		return fmt.Errorf("account not found in keystore")
	}

	err := s.repo.Unlock(account, sendTx.Password)
	if err != nil {
		fmt.Println("Failed to unlock account")
		return err
	}

	var (
		to       = common.HexToAddress(sendTx.Recipient)
		value    = new(big.Int).Mul(big.NewInt(1), big.NewInt(params.Ether))
		gasLimit = uint64(21000)
	)

	chainID, err := s.cl.ChainID(context.Background())
	if err != nil {
		fmt.Println("Fail to retrieve chainid")
		return err
	}

	nonce, err := s.cl.PendingNonceAt(context.Background(), account.Address)
	if err != nil {
		fmt.Println("Fail to retrieve nonce")
		return err
	}

	tipCap, _ := s.cl.SuggestGasTipCap(context.Background())
	feeCap, _ := s.cl.SuggestGasPrice(context.Background())

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
		fmt.Println("Failed to sign transaction")
		return err
	}

	return s.cl.SendTransaction(context.Background(), signedTx)
}

func formatEth(wei *big.Int) string {
	// Represent with Wei if smaller than Gwei
	if wei.Cmp(big.NewInt(params.GWei)) == -1 {
		return fmt.Sprintf("%d Wei", wei)
	}
	// Represent with Gwei if smaller than Ether
	if wei.Cmp(big.NewInt(params.Ether)) == -1 {
		return fmt.Sprintf("%d Gwei", wei)
	}
	// Show up to 5 decimals in Ether
	ether := new(big.Float).SetInt(wei)
	ether.Quo(ether, big.NewFloat(params.Ether))
	parts := strings.Split(ether.Text('f', 18), ".")
	// Hide decimals if first 5 decimals are 0
	if len(parts) < 2 || strings.HasPrefix(parts[1], strings.Repeat("0", 5)) {
		return fmt.Sprintf("%s Ether", parts[0])
	} else {
		return fmt.Sprintf("%s.%s Ether", parts[0], parts[1][:5])
	}
}
