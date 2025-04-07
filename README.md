# Token Exchange Service

A Cloud Run service that exchanges API keys for signed JWTs to be used with M-Lab services.

## Prerequisites

- Google Cloud SDK installed
- Access to the target GCP project
- Docker installed (if developing locally)

## Setup

1. Clone the repository:
```bash
git clone [repository-url]
cd token-exchange
```

2. Set up the signing key:
```bash
# Generate a new signing key (if you don't have one)
jose-util generate-key --use sig --alg EdDSA > private.json

# Add it as a secret to Google Cloud
gcloud secrets create token-exchange-private-key --data-file=private.json

# Clean up
rm private.json
```

## Deployment

The service can be built and deployed using the provided scripts:

```bash
# Build the container
./build.sh [PROJECT_ID] [REGION]

# Deploy to Cloud Run
./deploy.sh [PROJECT_ID] [REGION]
```

## API Endpoints

### 1. Token Exchange

Exchanges an API key for a signed JWT.

#### Request
```http
POST /token
Content-Type: application/json

{
    "api_key": "your-api-key"
}
```
#### Response
```json
{
    "token": "signed-jwt-token"
}
```

#### Example
```bash
curl -X POST https://[service-url]/token \
  -H "Content-Type: application/json" \
  -d '{"api_key": "your-api-key"}'
```

### 2. JWKS Endpoint

Returns the JSON Web Key Set (JWKS) containing the public key used to verify the JWTs.

#### Request
```http
GET /.well-known/jwks.json
```

#### Response
```json
{
    "keys": [
        {
            "kty": "RSA",
            "kid": "...",
            "n": "...",
            "e": "...",
            "alg": "...",
            "use": "sig"
        }
    ]
}
```

#### Example

```bash
curl https://[service-url]/.well-known/jwks.json
```


## Development

To run the service locally:

```bash
# Set up Go environment
go mod download

# Run the service
go run cmd/server/main.go
```

## Environment Variables

- `PROJECT_ID`: The Google Cloud project ID
- `PORT`: Port to run the service on (default: 8080)
- `PRIVATE_KEY_PATH`: Path to the private key file (default: /secrets/private.pem)

## Secret Configuration

The service expects the signing key to be mounted at `/secrets/private.pem` in JSON format. In Cloud Run, this is configured through the `--set-secrets` flag in the deployment script.
