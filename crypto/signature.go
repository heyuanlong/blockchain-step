package crypto

import (
	"crypto/ecdsa"
	"errors"
	secp "github.com/decred/dcrd/dcrec/secp256k1/v4"
	secp_ecdsa "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
)

//
//// SignCompact produces a compact signature of the data in hash with the given
//// private key on the given koblitz curve. The isCompressed  parameter should
//// be used to detail if the given signature should reference a compressed
//// public key or not. If successful the bytes of the compact signature will be
//// returned in the format:
//// <(byte of 27+public key solution)+4 if compressed >< padded bytes for signature R><padded bytes for signature S>
//// where the R and S parameters are padde up to the bitlengh of the curve.
//func SignCompact(key *secp.PrivateKey, hash []byte,
//	isCompressedKey bool) ([]byte, error) {
//
//	return secp_ecdsa.SignCompact(key, hash, isCompressedKey), nil
//}

// RecoverCompact verifies the compact signature "signature" of "hash" for the
// Koblitz curve in "curve". If the signature matches then the recovered public
// key will be returned as well as a boolean if the original key was compressed
// or not, else an error will be returned.
func RecoverCompact(signature, hash []byte) (*secp.PublicKey, bool, error) {
	return secp_ecdsa.RecoverCompact(signature, hash)
}

//// Sign generates an ECDSA signature over the secp256k1 curve for the provided
//// hash (which should be the result of hashing a larger message) using the
//// given private key. The produced signature is deterministic (same message and
//// same key yield the same signature) and canonical in accordance with RFC6979
//// and BIP0062.
//func Signx(key *secp.PrivateKey, hash []byte) *secp_ecdsa.Signature {
//	return secp_ecdsa.Sign(key, hash)
//}
//
////Serialize returns the ECDSA signature in the Distinguished Encoding Rules
//func Serialize(sig *secp_ecdsa.Signature) []byte {
//	return sig.Serialize()
//}
//
////ParseDERSignature parses a signature in the Distinguished Encoding Rules
//func ParseDERSignature(sig []byte) (*secp_ecdsa.Signature, error) {
//	return secp_ecdsa.ParseDERSignature(sig)
//}
//

//--------------------------------------------------------------------------------------------

// Output <compactSigRecoveryCode><32-byte R><32-byte S>.
func Sign(priv *secp.PrivateKey, hash []byte) []byte {
	return secp_ecdsa.SignCompact(priv, hash, false) // ref uncompressed pubkey
}

func VerifySignature(pubkey, hash, signature []byte) bool {
	if len(signature) != 65 {
		return false
	}
	key, err := secp.ParsePubKey(pubkey)
	if err != nil {
		return false
	}
	var r, s secp.ModNScalar
	if r.SetByteSlice(signature[1:33]) {
		return false // overflow
	}
	if s.SetByteSlice(signature[33:65]) {
		return false
	}

	sig := secp_ecdsa.NewSignature(&r, &s)
	return sig.Verify(hash, key)
}

func Ecrecover(hash, sig []byte) ([]byte, error) {
	pub, err := sigToPub(hash, sig)
	if err != nil {
		return nil, err
	}
	bytes := pub.SerializeUncompressed()
	return bytes, err
}

// SigToPub returns the public key that created the given signature.
func SigToEcdsaPub(hash, sig []byte) (*ecdsa.PublicKey, error) {
	pub, err := sigToPub(hash, sig)
	if err != nil {
		return nil, err
	}
	return pub.ToECDSA(), nil
}

func sigToPub(hash, sig []byte) (*PublicKey, error) {
	if len(sig) != 65 {
		return nil, errors.New("invalid signature")
	}

	pub, _, err := secp_ecdsa.RecoverCompact(sig, hash)
	return pub, err
}
