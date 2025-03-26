package handler

import (
	"encoding/json"
	"net/http"

	"github.com/m-lab/token-exchange/internal/auth"
)

type JWKSHandler struct {
	jwtSigner *auth.JWTSigner
}

func NewJWKSHandler(jwtSigner *auth.JWTSigner) *JWKSHandler {
	return &JWKSHandler{
		jwtSigner: jwtSigner,
	}
}

func (h *JWKSHandler) ServeJWKS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	publicKey := h.jwtSigner.GetPublicKey()
	jwk := auth.RSAPublicKeyToJWK(publicKey)
	jwks := auth.JWKS{Keys: []auth.JWK{jwk}}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwks)
}
