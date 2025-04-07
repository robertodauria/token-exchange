package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/uuid"
)

type JWTSigner struct {
	signer    jose.Signer
	publicJWK jose.JSONWebKey
}

type Claims struct {
	jwt.Claims
	Organization string `json:"org"`
}

// NewJWTSigner loads a private key from a JWK file and prepares a signer.
func NewJWTSigner(keyPath string) (*JWTSigner, error) {
	// Read private key JWK file
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file %s: %w", keyPath, err)
	}

	// Unmarshal the raw JSON into a jose.JSONWebKey
	var privateJWK jose.JSONWebKey
	if err := json.Unmarshal(keyBytes, &privateJWK); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JWK from %s: %w", keyPath, err)
	}

	// Check if the key is private (required for signing)
	if privateJWK.IsPublic() {
		return nil, fmt.Errorf("JWK in %s is not a private key", keyPath)
	}

	// Ensure the key has a Key ID (kid)
	if privateJWK.KeyID == "" {
		return nil, fmt.Errorf("JWK in %s must have a 'kid' (Key ID)", keyPath)
	}

	// Create the signer using the private key
	signerOpts := (&jose.SignerOptions{}).WithHeader(jose.HeaderKey("kid"), privateJWK.KeyID)

	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.EdDSA,
		Key:       privateJWK.Key,
	}, signerOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create jose signer: %w", err)
	}

	// Store the public part of the key for the JWKS endpoint
	publicJWK := privateJWK.Public()

	return &JWTSigner{
		signer:    signer,
		publicJWK: publicJWK,
	}, nil
}

func (s *JWTSigner) GenerateToken(org string) (string, error) {
	now := time.Now()
	expiry := now.Add(time.Hour) // Token expiry: 1 hour

	claims := Claims{
		Organization: org,
		Claims: jwt.Claims{
			ID:        uuid.New().String(),
			Issuer:    "token-exchange",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Expiry:    jwt.NewNumericDate(expiry),
		},
	}

	signedToken, err := jwt.Signed(s.signer).Claims(claims).Serialize()
	if err != nil {
		return "", fmt.Errorf("failed to sign claims: %w", err)
	}

	return signedToken, nil
}

// GetPublicJWK returns the public key in jose.JSONWebKey format.
func (s *JWTSigner) GetPublicJWK() jose.JSONWebKey {
	return s.publicJWK
}
