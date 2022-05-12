package fileWallet

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"heyuanlong/blockchain-step/accounts"
	"heyuanlong/blockchain-step/common"
	"heyuanlong/blockchain-step/core/tx"
	"heyuanlong/blockchain-step/crypto"
	"heyuanlong/blockchain-step/protocol"
	"io/ioutil"
	"math/big"
	"sync"
)

var deferFileWallet *FileWallet
var once sync.Once
func GetFileWallet()*FileWallet{
	once.Do(func() {
		deferFileWallet = NewFileWallet()
	})

	return deferFileWallet
}

type FileWallet struct {
	Scheme string

	isOpen     bool
	dir        string
	passPhrase string
	AddrMap    map[crypto.Address]*Key

	store StoreI
}

func NewFileWallet() *FileWallet {

	return &FileWallet{
		Scheme:  "file",
		store:   NewStoreFile(),
		AddrMap: make(map[crypto.Address]*Key),
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
	if err := w.loadAll(); err != nil {
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
	if !w.Contains(account) {
		return nil, errors.New("The wallet does not have the account")
	}
	v, ok := w.AddrMap[account.Address]
	if !ok {
		return nil, errors.New("The wallet does not have the account")
	}
	privByte, err := common.AESDecrypt(common.Hex2Bytes(v.PrivateKeyAes), []byte(w.passPhrase))
	if err != nil {
		log.Error(account.Address.String(), "AESDecrypt fail,may password error", err)
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
	key.Account.Address = crypto.PrivKeyToAddress(priv)
	key.URL = accounts.URL{Scheme: w.Scheme, Path: w.store.JoinPath(w.dir, key.Account.Address.String())}
	key.PrivateKeyAes = common.Bytes2Hex(common.AESEncrypt(crypto.PrivKeySerialize(priv), []byte(w.passPhrase)))

	w.AddrMap[key.Account.Address] = key
	w.store.StoreKey(key.URL.Path, key, w.passPhrase)

	return nil
}

//-创建账户
func (w *FileWallet) CreateAccount() (accounts.Account ,error ){
	if !w.isOpen {
		return accounts.Account{}, errors.New("The wallet was not opened")
	}

	priv, err := crypto.NewPrivateKey()
	if err != nil {
		log.Error("create account fail", err)
		return accounts.Account{},err
	}
	key := new(Key)
	key.Account.Address = crypto.PrivKeyToAddress(priv)
	key.URL = accounts.URL{Scheme: w.Scheme, Path: w.store.JoinPath(w.dir, key.Account.Address.String())}
	key.PrivateKeyAes = common.Bytes2Hex(common.AESEncrypt(priv.Serialize(), []byte(w.passPhrase)))

	w.AddrMap[key.Account.Address] = key
	w.store.StoreKey(key.URL.Path, key, w.passPhrase)

	return key.Account ,nil
}

//-
func (w *FileWallet) Accounts() []accounts.Account {
	acc := make([]accounts.Account, 0, len(w.AddrMap))
	for _, v := range w.AddrMap {
		acc = append(acc, v.Account)
	}
	return acc
}

//-
func (w *FileWallet) Contains(account accounts.Account) bool {
	if _, ok := w.AddrMap[account.Address]; ok {
		return true
	}
	return false
}

func (w *FileWallet) SignData(account accounts.Account, data []byte) ([]byte, error) {
	priv, err := w.Export(account)
	if err != nil {
		log.Error(account.Address.String(), "Export fail", err)
		return nil, err
	}

	return crypto.Sign(priv, data), nil
}
func (w *FileWallet) SignDataWithPassphrase(account accounts.Account, passphrase string, data []byte) ([]byte, error) {
	return []byte{}, nil
}

func (w *FileWallet) SignTx(account accounts.Account, txObj *protocol.Tx, chainID *big.Int) (*protocol.Tx, error) {
	hash := tx.DeferTxMgt.Hash(txObj)
	priv, err := w.Export(account)
	if err != nil {
		log.Error(account.Address.String(), "Export fail,maybe not in wallet:", err)
		return nil, err
	}
	txObj.Sign = crypto.Sign(priv, hash)
	return txObj, nil
}
func (w *FileWallet) SignTxWithPassphrase(account accounts.Account, passphrase string, txObj *protocol.Tx, chainID *big.Int) (*protocol.Tx, error) {
	return nil, nil

}

//-----------------------------------------------------------------
func (w *FileWallet) loadAll() error {
	if !w.isOpen {
		return errors.New("The wallet was not opened")
	}

	files, err := ioutil.ReadDir(w.dir)
	if err != nil {
		return err
	}
	for _, fi := range files {
		key, err := w.store.GetKey(crypto.Address{}, w.dir, fi.Name(), w.passPhrase)
		if err != nil {
			log.Trace("GetKey file ", w.dir, fi.Name())
			continue
		}
		w.AddrMap[key.Account.Address] = key
	}
	return nil
}
