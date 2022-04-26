package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Start() {
	r := NewRouteStruct("0.0.0.0", 7091)
	//r.SetMiddleware(kroute.MiddlewareCrossDomain())
	r.SetMiddleware(MiddlewareLoggerWithWriter(log.New().Out))

	//开启prometheus监控
	r.StartPrometheus()

	r.Load(NewTx())
	r.Load(NewAccount())

	r.Run()
}

//-------------------------------------------------------------
const (
	SUCCESS_STATUS  = 200
	OPERATION_WRONG = 20001
	PARAM_WRONG     = 20004
)

type DataIStruct struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ReturnError(c *gin.Context, status int, errors error) []byte {
	v := DataIStruct{
		Status:  status,
		Message: errors.Error(),
		Data:    nil,
	}
	jsonStr, _ := json.Marshal(v)
	c.Data(http.StatusOK, "text/plain", jsonStr)
	return jsonStr
}

func ReturnData(c *gin.Context, status int, data interface{}) []byte {
	v := DataIStruct{
		Status:  status,
		Message: "",
		Data:    data,
	}
	jsonStr, _ := json.Marshal(v)
	c.Data(http.StatusOK, "text/plain", jsonStr)
	return jsonStr
}
