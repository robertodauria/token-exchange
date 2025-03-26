package store

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/datastore"
)

const (
	OrgKind    = "Organization"
	APIKeyKind = "APIKey"
)

type DatastoreClient struct {
	client    *datastore.Client
	namespace string
}

type Organization struct {
	Name                  string    `datastore:"name"`
	Email                 string    `datastore:"email"`
	CreatedAt             time.Time `datastore:"created_at"`
	ProbabilityMultiplier *float64  `datastore:"probability_multiplier"`
}

type APIKey struct {
	CreatedAt time.Time `datastore:"created_at"`
	Key       string    `datastore:"key"`
}

func NewDatastoreClient(ctx context.Context, projectID, namespace string) (*DatastoreClient, error) {
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create datastore client: %w", err)
	}

	return &DatastoreClient{
		client:    client,
		namespace: namespace,
	}, nil
}

func (d *DatastoreClient) VerifyAPIKey(ctx context.Context, apiKey string) (string, error) {
	log.Printf("Attempting to verify API key: %s in namespace: %s", apiKey, d.namespace)

	// Query to find the API key and its parent organization
	q := datastore.NewQuery(APIKeyKind).
		Namespace(d.namespace).
		FilterField("key", "=", apiKey).Limit(1)

	log.Printf("Query details:")
	log.Printf("  Kind: %s", APIKeyKind)
	log.Printf("  Namespace: %s", d.namespace)
	log.Printf("  Filter: key = %s", apiKey)

	var apiKeys []APIKey
	keys, err := d.client.GetAll(ctx, q, &apiKeys)
	if err != nil {
		log.Printf("Error querying API key: %v", err)
		return "", fmt.Errorf("failed to query API key: %w", err)
	}

	log.Printf("Query returned %d results", len(keys))
	if len(keys) == 0 {
		return "", fmt.Errorf("invalid API key")
	}

	log.Printf("Parent organization: %s", keys[0].Parent.Name)

	// Get the organization ID from the parent key
	orgID := keys[0].Parent.Name

	return orgID, nil
}

func (d *DatastoreClient) Close() {
	d.client.Close()
}
