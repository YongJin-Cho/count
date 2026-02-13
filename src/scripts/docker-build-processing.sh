#!/bin/bash
set -e

TAG=${1:-latest}
MODULE_DIR="src/count-processing-service"

echo "Building count-processing-service:${TAG}..."

# Ensure go.sum is up to date
cd $MODULE_DIR
go mod tidy
cd -

docker build -t count-processing-service:${TAG} $MODULE_DIR
