package crypto

import (
	"heyuanlong/blockchain-step/common"
	"testing"
)

func TestPubKey1( t *testing.T)   {
	priv,_:=NewPrivateKey()

	secpPub:=priv.PubKey()
	pub := SerializeUncompressed(secpPub)
	t.Log(common.Bytes2Hex(pub))
	t.Log(BytesToAddress(Keccak256(pub[1:])[12:]))


}
