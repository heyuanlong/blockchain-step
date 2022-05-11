package api

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)




type blockGetByHashBind struct {
	Hash     string `form:"hash"  binding:"required"`
}

func (ts *ApiStruct) blockGetByHash(c *gin.Context) {
	var param blockGetByHashBind
	if err := c.Bind(&param); err != nil {
		log.Error(err)
		ReturnError(c, PARAM_WRONG, err)
		return
	}

	block ,_:=ts.db.GetBlockByHash(param.Hash)

	ReturnData(c, SUCCESS_STATUS, map[string]interface{}{
		"block": block,
	})
}


type blockGetByNumberBind struct {
	Number  uint64 `form:"number"  binding:"required"`
}

func (ts *ApiStruct) blockGetByNumber(c *gin.Context) {
	var param blockGetByNumberBind
	if err := c.Bind(&param); err != nil {
		log.Error(err)
		ReturnError(c, PARAM_WRONG, err)
		return
	}

	log.Println(param.Number)
	block ,_:=ts.db.GetBlockByNumber(param.Number)

	ReturnData(c, SUCCESS_STATUS, map[string]interface{}{
		"block": block,
	})
}
