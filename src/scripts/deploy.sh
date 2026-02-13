#!/bin/bash
set -e

NAMESPACE="count-system"

echo "Deploying Count Management System..."

# 1. Namespace
kubectl apply -f src/k8s/namespace.yaml

# 2. Secret
kubectl apply -f src/k8s/secret.yaml

# 3. Databases
kubectl apply -f src/k8s/database-management.yaml
kubectl apply -f src/k8s/database-processing.yaml

echo "Waiting for databases to be ready..."
kubectl wait --for=condition=available --timeout=60s deployment/management-db -n $NAMESPACE
kubectl wait --for=condition=available --timeout=60s deployment/processing-db -n $NAMESPACE

# 4. Services
kubectl apply -f src/k8s/count-processing-service.yaml
kubectl apply -f src/k8s/count-management-service.yaml

echo "Waiting for services to be ready..."
kubectl wait --for=condition=available --timeout=60s deployment/count-processing-service -n $NAMESPACE
kubectl wait --for=condition=available --timeout=60s deployment/count-management-service -n $NAMESPACE

# 5. Gateway API CRDs (Option A)
echo "Ensuring Gateway API CRDs are installed..."
if ! kubectl get crd gatewayclasses.gateway.networking.k8s.io > /dev/null 2>&1; then
    echo "Gateway API CRDs not found. Attempting to install..."
    # Using standard v1.0.0 release of Gateway API CRDs
    kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.0.0/standard-install.yaml || echo "Warning: Failed to install Gateway API CRDs. This might be due to insufficient permissions."
fi

# 6. Gateway / Ingress (Option B as fallback)
if kubectl get crd gatewayclasses.gateway.networking.k8s.io > /dev/null 2>&1; then
    echo "Applying Gateway API resources..."
    kubectl apply -f src/k8s/gateway.yaml
else
    echo "Gateway API CRDs missing. Falling back to standard Ingress..."
    kubectl apply -f src/k8s/ingress.yaml
fi

echo "Deployment complete!"
