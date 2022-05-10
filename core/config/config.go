package config

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

var Config = &jsonConfig{}

type jsonConfig struct {
	Node nodeConfig `json:"node"`
	DataDir string `json:"data_dir"`
}

type nodeConfig struct {
	Local string `json:"local"`
	NodeId    string `json:"node_id"`
	NodeAddrs    []string `json:"node_addrs"`
}

func InitConfParamByFile(file string) {
	cBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("Read Conf File fail:", err)
		return
	}

	err = json.Unmarshal(cBytes, &Config)
	if err != nil {
		log.Fatal("conf parse fail:", err)
		return
	}
	log.Printf("%#v\n", Config)

}

