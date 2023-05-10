package main

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
)

type KeystoreService struct {
	Ks           *keystore.KeyStore          // TODO: Make field private
	AccountIndex map[string]accounts.Account // TODO: Use map if large number of keystores
}

func NewKeystoreService(keystoreDir string) (*KeystoreService, error) {
	ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	return &KeystoreService{Ks: ks}, nil
}

func (s *KeystoreService) HasAddress(acc string) bool {
	return s.Ks.HasAddress(common.HexToAddress(acc))
}

// Need to be careful when using this method (TODO: Evaluate the need for retrieving SecretKey)
func (s *KeystoreService) SecreyKeyOf(acc string, pass string) (string, error) {
	account := accounts.Account{Address: common.HexToAddress(acc)}
	keyjson, err := s.Ks.Export(account, pass, pass)
	if err != nil {
		return "", err
	}
	key, err := keystore.DecryptKey(keyjson, pass)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key.PrivateKey.D.Bytes()), nil
}
