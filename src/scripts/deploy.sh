#!/bin/bash
set -e

NAMESPACE="count-collection-system"

echo "Starting deployment to namespace: ${NAMESPACE}"

# 1. Create Namespace
echo "Creating namespace..."
kubectl apply -f src/k8s/namespace.yaml

# 2. Deploy StatefulSet
echo "Deploying StatefulSet..."
kubectl apply -f src/k8s/count-api-service/statefulset.yaml

# 3. Deploy Service
echo "Deploying Service..."
kubectl apply -f src/k8s/count-api-service/service.yaml

# 4. Deploy Gateway API resources
echo "Deploying Gateway API..."
kubectl apply -f src/k8s/gateway.yaml

# 5. Wait for deployment to be ready
echo "Waiting for count-api-service to be ready..."
kubectl rollout status statefulset/count-api-service -n ${NAMESPACE} --timeout=120s

echo "Deployment completed successfully."
