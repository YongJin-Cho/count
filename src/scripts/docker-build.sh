#!/bin/bash
set -e

TAG=${1:-latest}

echo "Building all services..."

./src/scripts/docker-build-management.sh ${TAG}
./src/scripts/docker-build-processing.sh ${TAG}
