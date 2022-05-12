package api

import (
	"errors"
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

	account ,_ := ts.db.GetAccount(param.From)
	if account == nil || account.Id.Address == "" {
		ReturnError(c, OPERATION_WRONG, errors.New("account info get fail"))
		return
	}



	txObj := &protocol.Tx{
		Sender:    &protocol.Address{Address: param.From},
		To:        &protocol.Address{Address: param.To},
		Amount:    param.Amount,
		Nonce:     account.Nonce, //todo
		TimeStamp: uint64(time.Now().Unix()),
		Input:     []byte{},
	}
	tx.DeferTxMgt.Complete(txObj)

	w := fileWallet.GetFileWallet()
	w.Open(path.Join(config.Config.DataDir, types.WALLET_DIR), param.Password)


	address := crypto.HexToAddress(param.From)
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

	ReturnData(c, SUCCESS_STATUS, map[string]interface{}{
		"hash": txObj.Hash,
	})
}

type txBroadcastBind struct {
	From      string `form:"from"  binding:"required"`
	To        string `form:"to"  binding:"required"`
	Amount    uint64 `form:"amount"  binding:"required"`
	Sign      string `form:"sign" binding:"required"`
	Nonce     uint64 `form:"nonce" binding:"-"`
	Timestamp uint64 `form:"timestamp" binding:"required"`
	Hash 	  string `form:"hash" binding:"required"`
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
		Hash:     param.Hash,
	}



	err := tx.DeferTxMgt.Add(txObj)
	if err != nil {
		log.Error(err)
		ReturnError(c, OPERATION_WRONG, err)
		return
	}

	ReturnData(c, SUCCESS_STATUS, map[string]interface{}{
		"hash": txObj.Hash,
	})
}
