package api

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type tx struct {
}

func NewTx() *tx {
	return &tx{}
}
func (ts *tx) Load() []RouteWrapStruct {
	m := make([]RouteWrapStruct, 0)

	m = append(m, Wrap("GET|POST", "/tx/send", ts.send))

	return m
}

type sendBind struct {
	From   string  `form:"from"  binding:"required"`
	To     string  `form:"to"  binding:"required"`
	Amount float64 `form:"amount"  binding:"required"`
}

func (ts *tx) send(c *gin.Context) {
	var param sendBind
	if err := c.Bind(&param); err != nil {
		log.Error(err)
		ReturnError(c, PARAM_WRONG, err)
		return
	}
	//todo

	//检查时间
	//from是否在钱包里
	//from nonce
	//检验签名

	//todo
	//加入交易池
}
