package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	secp "github.com/decred/dcrd/dcrec/secp256k1/v4"
	"golang.org/x/crypto/sha3"
	"hash"
	"heyuanlong/blockchain-step/common"
)

type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	b := make([]byte, 32)
	d := sha3.NewLegacyKeccak256().(KeccakState)
	for _, b := range data {
		d.Write(b)
	}
	d.Read(b)
	return b
}

func FromECDSAPub(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(S256(), pub.X, pub.Y)
}

func PubkeyToAddress(p ecdsa.PublicKey) common.Address {
	pubBytes := FromECDSAPub(&p)
	return common.BytesToAddress(Keccak256(pubBytes[1:])[12:])
}

func PubkeyToAddress2(pub *secp.PublicKey) common.Address {
	p := pub.ToECDSA()
	pubBytes := FromECDSAPub(p)
	return common.BytesToAddress(Keccak256(pubBytes[1:])[12:])
}
