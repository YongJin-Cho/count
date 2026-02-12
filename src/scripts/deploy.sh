#!/bin/bash
set -e

NAMESPACE="count-collection-system"

echo "Starting deployment to namespace: ${NAMESPACE}"

# 1. Create Namespace
echo "Creating namespace..."
kubectl apply -f src/k8s/namespace.yaml

# 2. Deploy PVC
echo "Deploying PVC..."
kubectl apply -f src/k8s/count-api-service/pvc.yaml

# 3. Deploy Service
echo "Deploying Service..."
kubectl apply -f src/k8s/count-api-service/service.yaml

# 4. Deploy Deployment
echo "Deploying Deployment..."
kubectl apply -f src/k8s/count-api-service/deployment.yaml

# 5. Deploy Gateway API resources
echo "Deploying Gateway API..."
kubectl apply -f src/k8s/gateway.yaml

# 6. Wait for deployment to be ready
echo "Waiting for count-api-service to be ready..."
kubectl rollout status deployment/count-api-service -n ${NAMESPACE} --timeout=120s

echo "Deployment completed successfully."
