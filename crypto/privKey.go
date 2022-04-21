package crypto

import (
	secp "github.com/decred/dcrd/dcrec/secp256k1/v4"
)

type PrivateKey = secp.PrivateKey

// S256 returns a Curve which implements secp256k1.
func S256() *secp.KoblitzCurve {
	return secp.S256()
}


// NewPrivateKey is a wrapper for ecdsa.GenerateKey that returns a PrivateKey
// instead of the normal ecdsa.PrivateKey.
func NewPrivateKey() (*PrivateKey, error) {
	return secp.GeneratePrivateKey()
}

func  PubKey(p *PrivateKey) *secp.PublicKey {
	return p.PubKey()
}

// PrivKeyFromBytes returns a private and public key for `curve' based on the
// private key passed as an argument as a byte slice.
func PrivKeyFromBytes(pk []byte) (*PrivateKey, *secp.PublicKey) {
	privKey := secp.PrivKeyFromBytes(pk)

	return privKey, privKey.PubKey()
}

func  Serialize(p *PrivateKey) []byte {
	return p.Serialize()
}

