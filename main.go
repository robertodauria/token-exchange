package tokenexchange

import (
	"context"
	"encoding/json"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("TokenExchange", TokenExchange)
}

type RequestBody struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type ResponseBody struct {
	Token string `json:"token"`
	Error string `json:"error,omitempty"`
}

func TokenExchange(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		http.Error(w, "Error initializing Firebase", http.StatusInternalServerError)
		return
	}

	auth, err := app.Auth(ctx)
	if err != nil {
		http.Error(w, "Error getting Auth client", http.StatusInternalServerError)
		return
	}

	// Just create token with client ID as UID
	token, err := auth.CustomToken(ctx, req.ClientID)
	if err != nil {
		http.Error(w, "Error creating custom token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResponseBody{Token: token})
}
