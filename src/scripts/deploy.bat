@echo off
set NAMESPACE=count-collection-system

echo Starting deployment to namespace: %NAMESPACE%

echo Creating namespace...
kubectl apply -f src/k8s/namespace.yaml

echo Deploying StatefulSet...
kubectl apply -f src/k8s/count-api-service/statefulset.yaml

echo Deploying Service...
kubectl apply -f src/k8s/count-api-service/service.yaml

echo Deploying Gateway API...
kubectl apply -f src/k8s/gateway.yaml

echo Waiting for count-api-service to be ready...
kubectl rollout status statefulset/count-api-service -n %NAMESPACE% --timeout=120s

echo Deployment completed successfully.
