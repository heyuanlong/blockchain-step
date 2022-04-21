package fileWallet

import (
	"encoding/json"
	"errors"
	"heyuanlong/blockchain-step/accounts"
	"heyuanlong/blockchain-step/common"
	"heyuanlong/blockchain-step/core/types"
	"heyuanlong/blockchain-step/crypto"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	log "github.com/sirupsen/logrus"
)



type Key struct {
	Account accounts.Account `json:"account"`
	PrivateKeyAes string  `json:"privateKeyAes"`
}

type FileWallet struct {
	Scheme string

	isOpen bool
	dir string
	passPhrase string
	AddrMap map[common.Address]*Key
}

func NewFileWallet() *FileWallet{
	return &FileWallet{
		Scheme:"file",
	}
}

func (w *FileWallet) Open(dir string,passphrase string) error {
	if passphrase == ""{
		return errors.New("passphrase cannot be empty")
	}
	w.dir = dir
	w.passPhrase = passphrase
	w.isOpen = true
	if err := w.load();err != nil{
		return err
	}

	return nil
}

func (w *FileWallet) Close() error {
	w.dir = ""
	w.passPhrase = ""
	w.isOpen = false

	return nil
}

//备份钱包
func (w *FileWallet) BackUpWallet() error {
	return nil
}

//导入钱包
func (w *FileWallet) ImportWallet(dir string,passphrase string) error {
	if !w.isOpen{
		return errors.New("The wallet was not opened")
	}
	return nil
}
//导出账户
func (w *FileWallet) Export(accounts.Account) (*crypto.PrivateKey,error) {
	if !w.isOpen{
		return nil,errors.New("The wallet was not opened")
	}
	return nil,nil
}
//导入账户
func (w *FileWallet) Import(*crypto.PrivateKey) error{
	if !w.isOpen{
		return errors.New("The wallet was not opened")
	}
	return nil
}

//创建账户
func (w *FileWallet) CreateAccount() accounts.Account{
	priv,err:=crypto.NewPrivateKey()
	if err != nil{
		log.Error("create account fail",err)
	}
	key := new(Key)
	key.Account.Address = crypto.PubkeyToAddress2(priv.PubKey())
	key.Account.URL =  accounts.URL{Scheme: w.Scheme, Path: JoinPath(w.dir, keyFileName(key.Account.Address))}

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


//-----------------------------------------------------------------
func (w *FileWallet) load() error {
	files, err := ioutil.ReadDir(w.dir)
	if err != nil {
		return  err
	}
	for _, fi := range files {
		path := filepath.Join(w.dir, fi.Name())
		// Skip any non-key files from the folder
		if nonKeyFile(fi) {
			log.Trace("Ignoring file on account scan", "path", path)
			continue
		}
		fd, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fd.Close()
		key := new(Key)
		if err := json.NewDecoder(fd).Decode(key); err != nil {
			return err
		}
		w.AddrMap[key.Account.Address] =key
	}
	return nil
}