package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"heyuanlong/blockchain-step/storage/cache"
	"net/http"
)

type ApiStruct struct {
	port int
	db *cache.DBCache
}

func NewApi(port int ,db *cache.DBCache) *ApiStruct{
	return &ApiStruct{
		port:port,
		db:db,
	}

}



func (ts * ApiStruct)Run() {
	r := NewRouteStruct("0.0.0.0", ts.port)
	//r.SetMiddleware(kroute.MiddlewareCrossDomain())
	r.SetMiddleware(MiddlewareLoggerWithWriter(log.New().Out))

	//开启prometheus监控
	r.StartPrometheus()

	r.Load(ts)



	r.Run()
}


func (ts *ApiStruct) Load() []RouteWrapStruct {
	m := make([]RouteWrapStruct, 0)

	m = append(m, Wrap("GET|POST", "/account/create", ts.accountCreate))
	m = append(m, Wrap("GET|POST", "/account/info", ts.accountInfo))

	m = append(m, Wrap("GET|POST", "/block/getByHash", ts.blockGetByHash))
	m = append(m, Wrap("GET|POST", "/block/getByNumber", ts.blockGetByNumber))

	m = append(m, Wrap("GET|POST", "/tx/send", ts.txSend))
	m = append(m, Wrap("GET|POST", "/tx/broadcast", ts.txBroadcast))

	return m
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
	c.Data(http.StatusOK, "application/json; charset=utf-8", jsonStr)
	return jsonStr
}

func ReturnData(c *gin.Context, status int, data interface{}) []byte {
	v := DataIStruct{
		Status:  status,
		Message: "",
		Data:    data,
	}
	jsonStr, _ := json.Marshal(v)
	c.Data(http.StatusOK, "application/json; charset=utf-8", jsonStr)
	return jsonStr
}
