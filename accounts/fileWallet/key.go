package fileWallet

import (
	"encoding/json"
	"heyuanlong/blockchain-step/accounts"
	"heyuanlong/blockchain-step/crypto"
)

type Key struct {
	Account       accounts.Account `json:"account"`
	URL           accounts.URL     `json:"url"` // Optional resource locator within a backend
	PrivateKeyAes string           `json:"privateKeyAes"`
}

//--------------------------------------------------------------
type plainKeyJSON struct {
	Address    string       `json:"address"`
	PrivateKey string       `json:"privatekey"`
	URL        accounts.URL `json:"url"` // Optional resource locator within a backend
}

func (k *Key) MarshalJSON() (j []byte, err error) {
	jStruct := plainKeyJSON{
		k.Account.Address.String(),
		k.PrivateKeyAes,
		k.URL,
	}
	j, err = json.Marshal(jStruct)
	return j, err
}

func (k *Key) UnmarshalJSON(j []byte) (err error) {
	keyJSON := new(plainKeyJSON)
	err = json.Unmarshal(j, &keyJSON)
	if err != nil {
		return err
	}

	k.Account.Address = crypto.HexToAddress(keyJSON.Address)
	k.PrivateKeyAes = keyJSON.PrivateKey
	k.URL = keyJSON.URL

	return nil
}
