package api

import "github.com/gin-gonic/gin"

type account struct {
}

func NewAccount() *account {
	return &account{}
}
func (ts *account) Load() []RouteWrapStruct {
	m := make([]RouteWrapStruct, 0)

	m = append(m, Wrap("GET|POST", "/account/create", ts.create))

	return m
}

func (ts *account) create(c *gin.Context) {

}
