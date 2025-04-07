package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robertodauria/token-exchange/internal/auth"
	"github.com/robertodauria/token-exchange/internal/handler"
	"github.com/robertodauria/token-exchange/internal/store"
)

const (
	jwkPrivKeyPath = "/secrets/jwk-priv.json"
)

func main() {
	log.Printf("Starting token exchange service...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Port set to: %s", port)

	// Initialize JWT signer
	keyPath := os.Getenv("PRIVATE_KEY_PATH")
	if keyPath == "" {
		keyPath = jwkPrivKeyPath
	}
	log.Printf("Using private key from: %s", keyPath)

	jwtSigner, err := auth.NewJWTSigner(keyPath)
	if err != nil {
		log.Fatalf("Failed to initialize JWT signer: %v", err)
	}
	log.Printf("JWT signer initialized successfully")

	// Initialize Datastore client
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		log.Fatal("PROJECT_ID environment variable is required")
	}

	namespace := "credentials"

	datastoreClient, err := store.NewDatastoreClient(context.Background(), projectID, namespace)
	if err != nil {
		log.Fatalf("Failed to initialize Datastore client: %v", err)
	}
	defer datastoreClient.Close()

	mux := http.NewServeMux()

	// Register handlers with both JWT signer and Datastore client
	exchangeHandler := handler.NewExchangeHandler(jwtSigner, datastoreClient)
	jwksHandler := handler.NewJWKSHandler(jwtSigner)

	mux.HandleFunc("/token", exchangeHandler.Exchange)
	mux.HandleFunc("/.well-known/jwks.json", jwksHandler.ServeJWKS)

	// Add health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Printf("Received shutdown signal, gracefully shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()

	log.Printf("Server starting on port %s", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
	log.Printf("Server stopped")
}
