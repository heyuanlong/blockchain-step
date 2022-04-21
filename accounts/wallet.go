package accounts

import (
	"heyuanlong/blockchain-step/common"
	"heyuanlong/blockchain-step/crypto"
	"heyuanlong/blockchain-step/core/types"
	"math/big"
)


type Account struct {
	Address common.Address `json:"address"` // Ethereum account address derived from the key
	URL     URL            `json:"url"`     // Optional resource locator within a backend
}

// Wallet represents a software or hardware wallet that might contain one or more
// accounts (derived from the same seed).
type Wallet interface {
	//打开钱包
	Open(dir string,passphrase string) error
	//关闭
	Close() error

	//备份钱包
	BackUpWallet() error
	//导入钱包
	ImportWallet(dir string,passphrase string) error
	//导出账户
	Export(Account) (*crypto.PrivateKey,error)
	//导入账户
	Import(*crypto.PrivateKey) error

	//创建账户
	CreateAccount() Account
	//
	Accounts() []Account
	//
	Contains(account Account) bool

	SignData(account Account, mimeType string, data []byte) ([]byte, error)
	SignDataWithPassphrase(account Account, passphrase, mimeType string, data []byte) ([]byte, error)

	SignText(account Account, text []byte) ([]byte, error)
	SignTextWithPassphrase(account Account, passphrase string, hash []byte) ([]byte, error)

	SignTx(account Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)
	SignTxWithPassphrase(account Account, passphrase string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)
}
