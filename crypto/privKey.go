package crypto

import (
	secp "github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// S256 returns a Curve which implements secp256k1.
func S256() *secp.KoblitzCurve {
	return secp.S256()
}


// PrivKeyFromBytes returns a private and public key for `curve' based on the
// private key passed as an argument as a byte slice.
func PrivKeyFromBytes(pk []byte) (*secp.PrivateKey, *secp.PublicKey) {
	privKey := secp.PrivKeyFromBytes(pk)

	return privKey, privKey.PubKey()
}

// NewPrivateKey is a wrapper for ecdsa.GenerateKey that returns a PrivateKey
// instead of the normal ecdsa.PrivateKey.
func NewPrivateKey() (*secp.PrivateKey, error) {
	return secp.GeneratePrivateKey()
}

func  PubKey(p *secp.PrivateKey) *secp.PublicKey {
	return p.PubKey()
}