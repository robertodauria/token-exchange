package handler

import (
	"encoding/json"
	"net/http"

	"github.com/m-lab/token-exchange/internal/auth"
	"github.com/m-lab/token-exchange/internal/store"
)

type ExchangeHandler struct {
	jwtSigner *auth.JWTSigner
	store     *store.DatastoreClient
}

type TokenRequest struct {
	APIKey string `json:"api_key"`
}

type TokenResponse struct {
	Token string `json:"token"`
	Error string `json:"error,omitempty"`
}

func NewExchangeHandler(jwtSigner *auth.JWTSigner, store *store.DatastoreClient) *ExchangeHandler {
	return &ExchangeHandler{
		jwtSigner: jwtSigner,
		store:     store,
	}
}

func (h *ExchangeHandler) Exchange(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify API key and get organization ID
	orgID, err := h.store.VerifyAPIKey(r.Context(), req.APIKey)
	if err != nil {
		http.Error(w, "Invalid API key", http.StatusUnauthorized)
		return
	}

	// Generate JWT using organization ID
	token, err := h.jwtSigner.GenerateToken(orgID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	resp := TokenResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
