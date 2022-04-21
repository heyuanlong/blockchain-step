package crypto

import (
	secp "github.com/decred/dcrd/dcrec/secp256k1/v4"
	secp_ecdsa "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
)


// SignCompact produces a compact signature of the data in hash with the given
// private key on the given koblitz curve. The isCompressed  parameter should
// be used to detail if the given signature should reference a compressed
// public key or not. If successful the bytes of the compact signature will be
// returned in the format:
// <(byte of 27+public key solution)+4 if compressed >< padded bytes for signature R><padded bytes for signature S>
// where the R and S parameters are padde up to the bitlengh of the curve.
func SignCompact(key *secp.PrivateKey, hash []byte,
	isCompressedKey bool) ([]byte, error) {

	return secp_ecdsa.SignCompact(key, hash, isCompressedKey), nil
}

// RecoverCompact verifies the compact signature "signature" of "hash" for the
// Koblitz curve in "curve". If the signature matches then the recovered public
// key will be returned as well as a boolean if the original key was compressed
// or not, else an error will be returned.
func RecoverCompact(signature, hash []byte) (*secp.PublicKey, bool, error) {
	return secp_ecdsa.RecoverCompact(signature, hash)
}

// Sign generates an ECDSA signature over the secp256k1 curve for the provided
// hash (which should be the result of hashing a larger message) using the
// given private key. The produced signature is deterministic (same message and
// same key yield the same signature) and canonical in accordance with RFC6979
// and BIP0062.
func Sign(key *secp.PrivateKey, hash []byte) *secp_ecdsa.Signature {
	return secp_ecdsa.Sign(key, hash)
}

//Serialize returns the ECDSA signature in the Distinguished Encoding Rules
func Serialize(sig *secp_ecdsa.Signature) []byte {
	return sig.Serialize()
}

//ParseDERSignature parses a signature in the Distinguished Encoding Rules
func ParseDERSignature(sig []byte) (*secp_ecdsa.Signature, error) {
	return secp_ecdsa.ParseDERSignature(sig)
}

func Verify(sig *secp_ecdsa.Signature,hash []byte, pubKey *secp.PublicKey) bool {
	return sig.Verify(hash,pubKey)
}