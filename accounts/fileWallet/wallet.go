package fileWallet

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"heyuanlong/blockchain-step/accounts"
	"heyuanlong/blockchain-step/common"
	"heyuanlong/blockchain-step/core/types"
	"heyuanlong/blockchain-step/crypto"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
)

type FileWallet struct {
	Scheme string

	isOpen     bool
	dir        string
	passPhrase string
	AddrMap    map[common.Address]*Key
}

func NewFileWallet() *FileWallet {
	return &FileWallet{
		Scheme: "file",
	}
}

//-
func (w *FileWallet) Open(dir string, passphrase string) error {
	if passphrase == "" {
		return errors.New("passphrase cannot be empty")
	}
	w.dir = dir
	w.passPhrase = passphrase
	w.isOpen = true
	if err := w.load(); err != nil {
		return err
	}

	return nil
}

//-
func (w *FileWallet) Close() error {
	w.dir = ""
	w.passPhrase = ""
	w.isOpen = false

	return nil
}

//-备份钱包
func (w *FileWallet) BackUpWallet() error {
	return nil
}

//-导入钱包
func (w *FileWallet) ImportWallet(dir string, passphrase string) error {
	if !w.isOpen {
		return errors.New("The wallet was not opened")
	}
	return nil
}

//-导出账户
func (w *FileWallet) Export(account accounts.Account) (*crypto.PrivateKey, error) {
	if !w.isOpen {
		return nil, errors.New("The wallet was not opened")
	}
	if ! w.Contains(account){
		return nil, errors.New("The wallet does not have the account")
	}
	v,ok:= w.AddrMap[account.Address]
	if !ok{
		return nil, errors.New("The wallet does not have the account")
	}
	privByte ,err :=common.AESDecrypt( common.Hex2Bytes(v.PrivateKeyAes),[]byte(w.passPhrase))
	if err != nil{
		log.Error(account.Address.String(), "AESDecrypt fail",err)
		return nil, err
	}

	return crypto.PrivKeyFromBytes(privByte), nil
}

//-导入账户
func (w *FileWallet) Import(priv *crypto.PrivateKey) error {
	if !w.isOpen {
		return errors.New("The wallet was not opened")
	}

	key := new(Key)
	key.Account.Address = crypto.PubkeyToAddress2(priv.PubKey())
	key.URL = accounts.URL{Scheme: w.Scheme, Path: JoinPath(w.dir, keyFileName(key.Account.Address))}
	key.PrivateKeyAes = common.Bytes2Hex(common.AESEncrypt(crypto.Serialize(priv), []byte(w.passPhrase)))

	w.AddrMap[key.Account.Address] = key
	StoreKey(key.URL.Path,key)

	return nil
}

//-创建账户
func (w *FileWallet) CreateAccount() accounts.Account {
	priv, err := crypto.NewPrivateKey()
	if err != nil {
		log.Error("create account fail", err)
	}
	key := new(Key)
	key.Account.Address = crypto.PubkeyToAddress2(priv.PubKey())
	key.URL = accounts.URL{Scheme: w.Scheme, Path: JoinPath(w.dir, keyFileName(key.Account.Address))}
	key.PrivateKeyAes = common.Bytes2Hex(common.AESEncrypt(priv.Serialize(), []byte(w.passPhrase)))

	w.AddrMap[key.Account.Address] = key
	StoreKey(key.URL.Path,key)

	return key.Account
}

//-
func (w *FileWallet) Accounts() []accounts.Account {
	acc := make([]accounts.Account,0,len(w.AddrMap))
	for _, v := range w.AddrMap {
		acc =append(acc,v.Account)
	}
	return acc
}

//-
func (w *FileWallet) Contains(account accounts.Account) bool {
	if _,ok:= w.AddrMap[account.Address];ok{
		return true
	}
	return false
}

func (w *FileWallet) SignData(account accounts.Account, mimeType string, data []byte) ([]byte, error) {
	return []byte{}, nil
}
func (w *FileWallet) SignDataWithPassphrase(account accounts.Account, passphrase, mimeType string, data []byte) ([]byte, error) {
	return []byte{}, nil
}

func (w *FileWallet) SignText(account accounts.Account, text []byte) ([]byte, error) {
	return []byte{}, nil
}
func (w *FileWallet) SignTextWithPassphrase(account accounts.Account, passphrase string, hash []byte) ([]byte, error) {
	return []byte{}, nil
}

func (w *FileWallet) SignTx(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	return nil, nil
}
func (w *FileWallet) SignTxWithPassphrase(account accounts.Account, passphrase string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	return nil, nil
}

//-----------------------------------------------------------------
func (w *FileWallet) load() error {
	files, err := ioutil.ReadDir(w.dir)
	if err != nil {
		return err
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
		w.AddrMap[key.Account.Address] = key
	}
	return nil
}
