#!/bin/bash

PROJECT_ID=${1:-mlab-sandbox}
REGION=${2:-us-central1}

echo "Building for:"
echo "Project ID: $PROJECT_ID"
echo "Region: $REGION"

# Build the image
gcloud --project=$PROJECT_ID \
    builds submit --tag $REGION-docker.pkg.dev/$PROJECT_ID/m-lab/token-exchange
