package tx

import (
	"crypto/sha256"
	"fmt"
	"errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"heyuanlong/blockchain-step/common"
	"heyuanlong/blockchain-step/crypto"
	"heyuanlong/blockchain-step/protocol"
	"sync"
	"time"
)

var DeferTxMgt TxMgt

func init() {
	DeferTxMgt.poolCap = 1000
	DeferTxMgt.txPool = make( map[string]*protocol.Tx)
}

type TxMgt struct {
	sync.RWMutex
	poolCap      int
	txPool map[string]*protocol.Tx
}

func (ts *TxMgt) Complete(tx *protocol.Tx){

	//block.Hash
	tx.Hash = common.Bytes2HexWithPrefix(ts.Hash(tx))
}


func (ts *TxMgt) Bytes(tx *protocol.Tx) ([]byte, error) {
	b, err := proto.Marshal(tx)
	if err != nil {
		log.Error("to bytes fail", err)
		return []byte{}, err
	}
	return b, nil
}

func (ts *TxMgt) Add(tx *protocol.Tx)error {

	hash := common.Bytes2HexWithPrefix(ts.Hash(tx))
	if hash != tx.Hash {
		log.Error("hash != txObj.Hash")
		return  errors.New("hash != txObj.Hash")
	}

	if tx.Sender == nil || tx.Sender.Address == "" {
		return fmt.Errorf("交易数据from地址为空")
	}
	n := time.Now().Unix()
	if n-int64(tx.TimeStamp) > 48*3600 || int64(tx.TimeStamp)-n > 5*60 {
		return fmt.Errorf("交易时间戳错误")
	}
	if len(tx.Sign) == 0 {
		return fmt.Errorf("交易数据未签名")
	}

	//todo Sender 是否在钱包里
	//todo 检验 nonce
	//todo 模拟检验amount交易

	//检验签名
	accountAddr,err:=ts.Sender(tx)
	if err != nil{
		return err
	}

	//以太坊的发送者地址是直接从Sender方法里返回的
	//而这里在交易里有发送者地址，所以必须得判断下
	if tx.Sender.Address != accountAddr.Hex() {
		return fmt.Errorf("公钥地址和sender不匹配 p: %s, sender: %s",  accountAddr.Hex(), tx.Sender.Address)
	}

	//加入交易池
	if err := ts.AddToPool(tx);err != nil{
		return err
	}

	//todo 广播

	return nil
}

func (ts *TxMgt) Hash(tx *protocol.Tx) ([]byte) {
	t := &protocol.Tx{
		To:        tx.To,
		Amount:    tx.Amount,
		Nonce:     tx.Nonce,
		TimeStamp: tx.TimeStamp,
		Input:     tx.Input,
	}
	b, _ := proto.Marshal(t)

	sh := sha256.New()
	sh.Write(b)
	hash := sh.Sum(nil)

	return hash
}

//检验并获取sender
func (ts *TxMgt) Sender(tx *protocol.Tx) (crypto.Address, error) {
	hash  := ts.Hash(tx)

	pub,err  :=crypto.Ecrecover(hash,tx.Sign)
	if err != nil{
		return crypto.Address{}, err
	}

	return  crypto.BytesToAddress(crypto.Keccak256(pub[1:])[12:]),nil
}

