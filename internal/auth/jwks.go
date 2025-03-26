package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"math/big"
)

type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid,omitempty"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
	Alg string `json:"alg"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

// RSAPublicKeyToJWK converts an RSA public key to JWK format
func RSAPublicKeyToJWK(publicKey *rsa.PublicKey) JWK {
	// Convert modulus and exponent to base64url encoding
	n := base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes())

	return JWK{
		Kty: "RSA",
		Kid: "1",
		Use: "sig",
		N:   n,
		E:   e,
		Alg: "RS256",
	}
}

// GetPublicKey extracts the public key from the private key
func (s *JWTSigner) GetPublicKey() *rsa.PublicKey {
	return &s.privateKey.PublicKey
}
