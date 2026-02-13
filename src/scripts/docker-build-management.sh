#!/bin/bash
set -e

TAG=${1:-latest}
MODULE_DIR="src/count-management-service"

echo "Building count-management-service:${TAG}..."

# Ensure go.sum is up to date
cd $MODULE_DIR
go mod tidy
cd -

docker build -t count-management-service:${TAG} $MODULE_DIR
