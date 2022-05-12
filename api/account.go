package api

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)



func (ts *ApiStruct) accountCreate(c *gin.Context) {

}


type accountInfoBind struct {
	Address     string `form:"address"  binding:"required"`
}

func (ts *ApiStruct) accountInfo(c *gin.Context) {
	var param accountInfoBind
	if err := c.Bind(&param); err != nil {
		log.Error(err)
		ReturnError(c, PARAM_WRONG, err)
		return
	}

	account ,_:=ts.db.GetAccount(param.Address)

	ReturnData(c, SUCCESS_STATUS, map[string]interface{}{
		"account": account,
	})
}
