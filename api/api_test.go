package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"heyuanlong/blockchain-step/common"
	"testing"
)

func Test_AccountCreate(t *testing.T) {

	url := common.BuildUrlParam("http://127.0.0.1:3000/account/create", map[string]interface{}{

	})
	body, err := common.HttpGet(url, map[string]string{}, 10)
	if err != nil {
		fmt.Println(body)
		t.Error(err)
	}

	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, body, "", "\t")
	fmt.Println("body:\n", string(prettyJSON.Bytes()))
}


func Test_TxSend(t *testing.T) {

	url := common.BuildUrlParam("http://127.0.0.1:3000/tx/send", map[string]interface{}{
		"from":"0x2C35a52A51B742C3c1218e6c053059268c85BaB7",
		"to":"0xe2cD88bE7757921451A90A2A06df3831a2B38698",
		"amount":1000,
		"password":"123456",
	})

	body, err := common.HttpGet(url, map[string]string{}, 10)
	if err != nil {
		fmt.Println(body)
		t.Error(err)
	}

	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, body, "", "\t")
	fmt.Println("body:\n", string(prettyJSON.Bytes()))
}


func Test_BlockGetByNumber(t *testing.T) {

	url := common.BuildUrlParam("http://127.0.0.1:3000/block/getByNumber", map[string]interface{}{
		"number":193,
	})

	body, err := common.HttpGet(url, map[string]string{}, 10)
	if err != nil {
		fmt.Println(body)
		t.Error(err)
	}

	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, body, "", "\t")
	fmt.Println("body:\n", string(prettyJSON.Bytes()))
}
