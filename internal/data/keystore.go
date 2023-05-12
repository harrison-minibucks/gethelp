package data

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/harrison-minibucks/gethelp/internal/biz"
	"github.com/harrison-minibucks/gethelp/internal/conf"
)

type keystoreRepo struct {
	ks *keystore.KeyStore
	// accountIndex map[string]accounts.Account // TODO: Use map if large number of keystores
}

func NewKeystoreRepo(conf *conf.Config) biz.KeystoreRepo {
	keystoreDir := ""
	if strings.Compare(conf.KeystoreDir, "") != 0 {
		keystoreDir = conf.KeystoreDir
	} else {
		if runtime.GOOS == "windows" {
			keystoreDir = os.Getenv("LocalAppData") + "\\Ethereum\\keystore"
		} else {
			if home, err := os.UserHomeDir(); err != nil {
				panic(err)
			} else {
				keystoreDir = home + "/ethereum/keystore"
			}
		}
	}
	ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	fmt.Println(ks.Accounts())
	return &keystoreRepo{ks: ks}
}

func (s *keystoreRepo) HasAddress(acc string) bool {
	return s.ks.HasAddress(common.HexToAddress(acc))
}

func (s *keystoreRepo) SignTransaction(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	return s.ks.SignTx(account, tx, chainID)
}

// Need to be careful when using this method (TODO: Evaluate the need for retrieving SecretKey)
func (s *keystoreRepo) SecretKeyOf(acc string, pass string) (string, error) {
	account := accounts.Account{Address: common.HexToAddress(acc)}
	keyjson, err := s.ks.Export(account, pass, pass)
	if err != nil {
		return "", err
	}
	key, err := keystore.DecryptKey(keyjson, pass)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key.PrivateKey.D.Bytes()), nil
}

func (s *keystoreRepo) Unlock(account accounts.Account, password string) error {
	return s.ks.Unlock(account, password)
}
