package crypto

import "testing"

func TestPriv(t *testing.T) {
	priv, _ := NewPrivateKey()
	t.Log(PrivKeyToAddress(priv))

}
