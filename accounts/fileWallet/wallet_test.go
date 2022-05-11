package fileWallet

import (
	"heyuanlong/blockchain-step/accounts"
	"heyuanlong/blockchain-step/common"
	"heyuanlong/blockchain-step/crypto"
	"testing"
)

func TestWallet(t *testing.T) {
	w := NewFileWallet()
	w.Open("./walletdir", "123456")
	defer w.Close()

	p, _ := crypto.NewPrivateKey()
	if err := w.Import(p); err != nil {
		t.Fatal(err)
	}
	t.Log("Import:", crypto.PrivKeyToAddress(p))
	pe, err := w.Export(accounts.Account{Address: crypto.PrivKeyToAddress(p)})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Export:",  crypto.PrivKeyToAddress(pe))

	acc,_ := w.CreateAccount()
	t.Log("CreateAccount:", acc.Address)

	accs := w.Accounts()
	for _, acc := range accs {
		t.Log("accs:", acc.Address)
	}

	b := w.Contains(accounts.Account{Address: crypto.PrivKeyToAddress(p)})
	t.Log("b:", b)

	sign, err := w.SignData(accounts.Account{Address: crypto.PrivKeyToAddress(p)}, []byte("123456"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(common.Bytes2Hex(sign))

}
