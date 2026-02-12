#!/bin/bash
set -e

# Default tag
TAG=${1:-latest}
IMAGE_NAME="count-api-service"

echo "Building Docker image ${IMAGE_NAME}:${TAG}..."

# Move to the component directory
cd src/count-api-service

# Ensure go.sum is updated
echo "Running go mod tidy..."
go mod tidy

# Build the Docker image
# Using --network=host for local development if needed, but standard build is usually fine
docker build -t "${IMAGE_NAME}:${TAG}" .

echo "Build complete: ${IMAGE_NAME}:${TAG}"
