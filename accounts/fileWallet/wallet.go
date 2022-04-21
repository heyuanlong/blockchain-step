package fileWallet

import (
	"heyuanlong/blockchain-step/core/types"
	"heyuanlong/blockchain-step/accounts"
	"heyuanlong/blockchain-step/crypto"
	"math/big"
)

type FileWallet struct {

}


func (w *FileWallet) Open(dir string,passphrase string) error {
	return nil
}



func (w *FileWallet) Close() error {
	return nil
}

//备份钱包
func (w *FileWallet) BackUp() error {
	return nil
}
//导出账户
func (w *FileWallet) Export(accounts.Account) (*crypto.PrivateKey,error) {
	return nil,nil
}
//导入账户
func (w *FileWallet) mport(*crypto.PrivateKey) error{
	return nil
}

//创建账户
func (w *FileWallet) CreateAccount() accounts.Account{
	return accounts.Account{}
}
//
func (w *FileWallet) Accounts() []accounts.Account{
	return []accounts.Account{}
}
//
func (w *FileWallet) Contains(account accounts.Account) bool{
	return true
}


func (w *FileWallet) SignData(account accounts.Account, mimeType string, data []byte) ([]byte, error){
	return []byte{},nil
}
func (w *FileWallet) SignDataWithPassphrase(account accounts.Account, passphrase, mimeType string, data []byte) ([]byte, error){
	return []byte{},nil
}


func (w *FileWallet) SignText(account accounts.Account, text []byte) ([]byte, error){
	return []byte{},nil
}
func (w *FileWallet) SignTextWithPassphrase(account accounts.Account, passphrase string, hash []byte) ([]byte, error){
	return []byte{},nil
}


func (w *FileWallet) SignTx(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error){
	return nil,nil
}
func (w *FileWallet) SignTxWithPassphrase(account accounts.Account, passphrase string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error){
	return nil,nil
}