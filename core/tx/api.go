package tx

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"heyuanlong/blockchain-step/api"
	"heyuanlong/blockchain-step/common"
	"heyuanlong/blockchain-step/protocol"
)

func (ts *TxMgt) Load() []api.RouteWrapStruct {
	m := make([]api.RouteWrapStruct, 0)

	m = append(m, api.Wrap("GET|POST", "/tx/send", ts.send))

	return m
}

type sendBind struct {
	From      string `form:"from"  binding:"required"`
	To        string `form:"to"  binding:"required"`
	Amount    uint64 `form:"amount"  binding:"required"`
	Sign      string `form:"sign" binding:"required"`
	PublicKey string `form:"public_key" binding:"required"`
	Nonce     uint64 `form:"nonce" binding:"-"`
	Timestamp uint64 `form:"timestamp" binding:"required"`
}

func (ts *TxMgt) send(c *gin.Context) {
	var param sendBind
	if err := c.Bind(&param); err != nil {
		log.Error(err)
		api.ReturnError(c, api.PARAM_WRONG, err)
		return
	}

	tx := &protocol.Tx{
		Sender:    &protocol.Address{Address: param.From},
		To:        &protocol.Address{Address: param.To},
		Amount:    param.Amount,
		Nonce:     param.Nonce,
		Sign:      common.FromHex(param.Sign),
		PublicKey: common.FromHex(param.PublicKey),
		TimeStamp: param.Timestamp,
		Input:     []byte{},
	}
	ts.Add(tx)
}
