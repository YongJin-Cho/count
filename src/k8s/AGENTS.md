# Kubernetes Configuration Management Basic Principles

## 1. Naming Rules
- **Namespace**: `count-system` (kebab-case based on CountManagementSystem).
- **Resources**: All resources should use kebab-case and include the service name prefix where appropriate (e.g., `count-management-service`, `management-db`).
- **Labels**: Use standard labels:
  - `app.kubernetes.io/name`: Name of the application
  - `app.kubernetes.io/instance`: Instance of the application
  - `app.kubernetes.io/component`: Component type (backend, database)
  - `app.kubernetes.io/part-of`: System name (`count-management-system`)

## 2. Namespace Strategy
- All resources are deployed in the `count-system` namespace to provide isolation and resource management boundaries.

## 3. Resource Management
- **Deployments**: Used for stateless backend services (`count-management-service`, `count-processing-service`).
- **StatefulSets**: Used for databases (`management-db`, `processing-db`) to ensure stable network identifiers and persistent storage.
- **Services**:
  - `ClusterIP` for internal communication.
- **Gateway API**:
  - `GatewayClass`: `kong`
  - `Gateway`: `count-gateway` handles external traffic.
  - `HTTPRoute`: Maps external paths to backend services with method-based routing.
- **Ingress Fallback**: If Gateway API CRDs are not available in the cluster, a standard `Ingress` resource (`src/k8s/ingress.yaml`) is used as a fallback to ensure accessibility.

## 4. Configuration and Secrets
- **Environment Variables**: Used for database connection strings and service endpoints.
- **Secrets**: Used for sensitive information like database passwords.

## 5. Persistence
- **PersistentVolumeClaims**: Used by StatefulSets for database data persistence.
- **StorageClass**: Uses default storage class unless specified otherwise.

## 6. Health Checks
- Services include `livenessProbe` and `readinessProbe` pointing to `/health` endpoint if available, or basic TCP checks for databases.
