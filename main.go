package main

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

const MY_PUBLIC_ADDRESS = "0x929548598a3b93362c5aa2a24de190d18e657ae0"
const RECIPIENT_ADDRESS = "0xb02A2EdA1b317FBd16760128836B0Ac59B560e91"
const KEYSTORE_DIR = "C:\\Users\\Shiro\\AppData\\Local\\Ethereum\\keystore"
const PASSWORD = "P@ssw0rd" // TODO: Remove

func main() {
	// Connect to the Ethereum node (Linux - $HOME/.ethereum/geth.ipc)
	cl, err := ethclient.Dial("\\\\.\\pipe\\geth.ipc")
	if err != nil {
		panic(err)
	}
	_ = cl

	ksService, err := NewKeystoreService(KEYSTORE_DIR)
	if err != nil {
		panic(err)
	}
	if err := sendTransaction(cl, ksService); err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully sent transaction", time.Now())
		// A delay in state update may occur (prior to commit), where pending balance is not reflected correctly
		time.Sleep(time.Millisecond * 1000)
	}
	printBalance(cl, MY_PUBLIC_ADDRESS)
	printBalance(cl, RECIPIENT_ADDRESS)

	fmt.Println("Done")
}

// Prints the balance of an account
func printBalance(cl *ethclient.Client, accountAddress string) error {
	account := common.HexToAddress(accountAddress)
	balance, err := cl.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return err
	}
	pendingBalance, err := cl.PendingBalanceAt(context.Background(), account)
	if err != nil {
		return err
	}
	fmt.Println("Account:", accountAddress)
	fmt.Println("Balance:", formatEth(balance))
	if balance.Cmp(pendingBalance) != 0 {
		fmt.Println("Pending Balance:", formatEth(pendingBalance))
	}
	fmt.Println()
	return nil
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

func sendTransaction(cl *ethclient.Client, ks *KeystoreService) error {
	accountAddress := MY_PUBLIC_ADDRESS
	account := accounts.Account{Address: common.HexToAddress(accountAddress)}
	if !ks.HasAddress(accountAddress) {
		fmt.Println("Account not found in keystore")
		return fmt.Errorf("account not found in keystore")
	}

	err := ks.Ks.Unlock(account, PASSWORD)
	if err != nil {
		fmt.Println("Failed to unlock account")
		return err
	}

	var (
		to       = common.HexToAddress(RECIPIENT_ADDRESS)
		value    = new(big.Int).Mul(big.NewInt(1), big.NewInt(params.Ether))
		gasLimit = uint64(21000)
	)

	chainID, err := cl.ChainID(context.Background())
	if err != nil {
		fmt.Println("Fail to retrieve chainid")
		return err
	}

	nonce, err := cl.PendingNonceAt(context.Background(), account.Address)
	if err != nil {
		fmt.Println("Fail to retrieve nonce")
		return err
	}

	tipCap, _ := cl.SuggestGasTipCap(context.Background())
	feeCap, _ := cl.SuggestGasPrice(context.Background())

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

	signedTx, err := ks.Ks.SignTx(account, tx, chainID)
	if err != nil {
		fmt.Println("Failed to sign transaction")
		return err
	}

	return cl.SendTransaction(context.Background(), signedTx)
}
