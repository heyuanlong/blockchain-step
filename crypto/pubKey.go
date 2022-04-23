package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	secp "github.com/decred/dcrd/dcrec/secp256k1/v4"
)


type PublicKey = secp.PublicKey

// These constants define the lengths of serialized public keys.
const (
	PubKeyBytesLenCompressed = 33
)

const (
	pubkeyCompressed   byte = 0x2 // y_bit + x coord
	pubkeyUncompressed byte = 0x4 // x coord + y coord
	pubkeyHybrid       byte = 0x6 // y_bit + x coord + y coord
)

// IsCompressedPubKey returns true the the passed serialized public key has
// been encoded in compressed format, and false otherwise.
func IsCompressedPubKey(pubKey []byte) bool {
	// The public key is only compressed if it is the correct length and
	// the format (first byte) is one of the compressed pubkey values.
	return len(pubKey) == PubKeyBytesLenCompressed &&
		(pubKey[0]&^byte(0x1) == pubkeyCompressed)
}

// ParsePubKey parses a public key for a koblitz curve from a bytestring into a
// ecdsa.Publickey, verifying that it is valid. It supports compressed,
// uncompressed and hybrid signature formats.
func ParsePubKey(pubKeyStr []byte) (*PublicKey, error) {
	return secp.ParsePubKey(pubKeyStr)
}

// SerializeUncompressed serializes a public key in the 65-byte uncompressed
// format.
func SerializeUncompressed(pub *PublicKey) []byte {
	return pub.SerializeUncompressed()
}

// SerializeCompressed serializes a public key in the 33-byte compressed format.
func SerializeCompressed(pub *PublicKey) []byte {
	return pub.SerializeCompressed()
}

func PubkeyEcdsaToByte(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(S256(), pub.X, pub.Y)
}

func PubkeyEcdsaToAddress(p ecdsa.PublicKey) Address {
	pubBytes := PubkeyEcdsaToByte(&p)
	return BytesToAddress(Keccak256(pubBytes[1:])[12:])
}

func PubkeyToAddress2(pub *PublicKey) Address {
	p := pub.ToECDSA()
	pubBytes := PubkeyEcdsaToByte(p)
	return BytesToAddress(Keccak256(pubBytes[1:])[12:])
}