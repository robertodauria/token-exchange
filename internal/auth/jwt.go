package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTSigner struct {
	privateKey *rsa.PrivateKey
}

type Claims struct {
	Organization string `json:"org"`
	jwt.RegisteredClaims
}

func NewJWTSigner(keyPath string) (*JWTSigner, error) {
	// Read private key file
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	// Parse PEM block
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}

	// Parse private key
	pkcs8Key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	var ok bool
	privateKey, ok := pkcs8Key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not RSA")
	}

	return &JWTSigner{privateKey: privateKey}, nil

}

func (s *JWTSigner) GenerateToken(org string) (string, error) {
	now := time.Now()

	claims := Claims{
		Organization: org,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "token-exchange-service",
			Audience:  jwt.ClaimStrings{"token-exchange-service"},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "1"

	signedToken, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}
