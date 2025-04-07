package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-jose/go-jose/v4"
	"github.com/robertodauria/token-exchange/internal/auth"
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

	publicJWK := h.jwtSigner.GetPublicJWK()
	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{publicJWK},
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600, must-revalidate")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(jwks); err != nil {
		// Log the internal error
		log.Printf("Error encoding JWKS response: %v", err)
		http.Error(w, "Failed to encode JWKS response", http.StatusInternalServerError)
		return
	}
}
