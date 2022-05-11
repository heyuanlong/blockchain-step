package api

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"heyuanlong/blockchain-step/accounts"
	"heyuanlong/blockchain-step/accounts/fileWallet"
	"heyuanlong/blockchain-step/common"
	"heyuanlong/blockchain-step/core/config"
	"heyuanlong/blockchain-step/core/tx"
	"heyuanlong/blockchain-step/core/types"
	"heyuanlong/blockchain-step/crypto"
	"heyuanlong/blockchain-step/protocol"
	"math/big"
	"path"
	"time"
)


type txSendBind struct {
	From     string `form:"from"  binding:"required"`
	To       string `form:"to"  binding:"required"`
	Amount   uint64 `form:"amount"  binding:"required"`
	Password string `form:"password" binding:"required"`
}

func (ts *ApiStruct) txSend(c *gin.Context) {
	var param txSendBind
	if err := c.Bind(&param); err != nil {
		log.Error(err)
		ReturnError(c, PARAM_WRONG, err)
		return
	}

	address := crypto.HexToAddress(param.From)

	txObj := &protocol.Tx{
		Sender:    &protocol.Address{Address: param.From},
		To:        &protocol.Address{Address: param.To},
		Amount:    param.Amount,
		Nonce:     0, //todo
		TimeStamp: uint64(time.Now().Unix()),
		Input:     []byte{},
	}

	w := fileWallet.GetFileWallet()
	w.Open(path.Join(config.Config.DataDir, types.WALLET_DIR), param.Password)

	//签名交易
	if _, err := w.SignTx(accounts.Account{Address: address}, txObj, big.NewInt(0)); err != nil {
		log.Error(err)
		ReturnError(c, OPERATION_WRONG, err)
		return
	}

	err := tx.DeferTxMgt.Add(txObj)
	if err != nil {
		log.Error(err)
		ReturnError(c, OPERATION_WRONG, err)
		return
	}

	hash ,_:= tx.DeferTxMgt.Hash(txObj)
	ReturnData(c, SUCCESS_STATUS, map[string]interface{}{
		"hash": common.Bytes2HexWithPrefix(hash),
	})
}

type txBroadcastBind struct {
	From      string `form:"from"  binding:"required"`
	To        string `form:"to"  binding:"required"`
	Amount    uint64 `form:"amount"  binding:"required"`
	Sign      string `form:"sign" binding:"required"`
	PublicKey string `form:"public_key" binding:"required"`
	Nonce     uint64 `form:"nonce" binding:"-"`
	Timestamp uint64 `form:"timestamp" binding:"required"`
}

func (ts *ApiStruct) txBroadcast(c *gin.Context) {
	var param txBroadcastBind
	if err := c.Bind(&param); err != nil {
		log.Error(err)
		ReturnError(c, PARAM_WRONG, err)
		return
	}

	txObj := &protocol.Tx{
		Sender:    &protocol.Address{Address: param.From},
		To:        &protocol.Address{Address: param.To},
		Amount:    param.Amount,
		Nonce:     param.Nonce,
		Sign:      common.FromHex(param.Sign),
		TimeStamp: param.Timestamp,
		Input:     []byte{},
	}

	tx.DeferTxMgt.Add(txObj)
}
