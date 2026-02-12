# Kubernetes Configuration Management Principles

## Naming Rules
- **Namespace**: kebab-case based on MSA.System name (e.g., `count-collection-system`).
- **Resources**: kebab-case starting with service name (e.g., `count-api-service`, `count-api-service-pvc`).
- **Labels**: 
  - `app: count-api-service`
  - `system: count-collection-system`

## Namespace Strategy
- Each MSA.System has its own dedicated namespace to ensure isolation.
- All resources for a system must be deployed in its specific namespace.

## Resource Management
- **Deployment**: Used for stateless services.
- **StatefulSet**: Used for services requiring persistent storage or stable identifiers.
- **Service**: ClusterIP is preferred for internal communication.
- **Gateway API**: Used for external traffic routing with HTTP method support.

## Security & Reliability
- **Resource Limits**: Always define CPU and memory requests/limits based on `system.json`.
- **Probes**: Include `livenessProbe` and `readinessProbe` for all services.
- **Secrets**: Use Kubernetes Secrets for sensitive information.
