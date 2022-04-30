package tx

import (
	"crypto/sha256"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"heyuanlong/blockchain-step/common"
	"heyuanlong/blockchain-step/crypto"
	"heyuanlong/blockchain-step/protocol"
	"sync"
	"time"
)

type TxMgt struct {
	Tx protocol.Tx

	sync.RWMutex
	poolCap      int
	txPool map[string]*protocol.Tx
}

func (ts TxMgt) Bytes() ([]byte, error) {
	b, err := proto.Marshal(&ts.Tx)
	if err != nil {
		log.Error("to bytes fail", err)
		return []byte{}, err
	}
	return b, nil
}

func (ts TxMgt) SetSign(sign []byte) {
	ts.Tx.Sign = sign
}

func (ts TxMgt) Add(tx *protocol.Tx)error {

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
	if len(tx.PublicKey) == 0 {
		return fmt.Errorf("交易数据公钥为空")
	}

	//
	pub ,err  := crypto.ParsePubKey( tx.PublicKey)
	if err != nil{
		return err
	}
	accountAddr := crypto.PubkeyToAddress(pub)
	if tx.Sender.Address != accountAddr.Hex() {
		return fmt.Errorf("公钥地址和sender不匹配 p: %s, sender: %s",  accountAddr.Hex(), tx.Sender.Address)
	}

	//todo Sender 是否在钱包里
	//todo 检验 nonce

	//检验 签名
	b,err:=ts.VerifySignedTx(tx)
	if err != nil{
		return err
	}
	if !b{
		return fmt.Errorf("验签不通过")
	}


	//todo
	//加入交易池
	ts.AddToPool(tx)

	return nil
}

func (ts *TxMgt) Hash(tx *protocol.Tx) ([]byte,error) {
	t := &protocol.Tx{
		To:        tx.To,
		Amount:    tx.Amount,
		Nonce:     tx.Nonce,
		TimeStamp: tx.TimeStamp,
		Input:     tx.Input,
	}
	b, err := proto.Marshal(t)
	if err != nil {
		return []byte{}, err
	}
	sh := sha256.New()
	sh.Write(b)
	hash := sh.Sum(nil)

	return hash,nil
}

func (ts *TxMgt) VerifySignedTx(tx *protocol.Tx) (bool, error) {
	hash ,err := ts.Hash(tx)
	if err != nil {
		return false, err
	}

	ret :=crypto.VerifySignature(tx.PublicKey,hash,tx.Sign)
	return ret,nil
}
