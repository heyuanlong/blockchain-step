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

func PrivKeyToPubKey(p *PrivateKey) *PublicKey {
	return p.PubKey()
}

func PrivKeyToAddress(p *PrivateKey) Address {
	pub := PrivKeyToPubKey(p)
	ecdsaPub := pub.ToECDSA()
	pubBytes := PubkeyEcdsaToByte(ecdsaPub)
	return BytesToAddress(Keccak256(pubBytes[1:])[12:])
}


// PrivKeyFromBytes returns a private and public key for `curve' based on the
// private key passed as an argument as a byte slice.
func PrivKeyFromBytes(pk []byte) (*PrivateKey) {
	privKey := secp.PrivKeyFromBytes(pk)

	return privKey
}

func PrivKeySerialize(p *PrivateKey) []byte {
	return p.Serialize()
}
