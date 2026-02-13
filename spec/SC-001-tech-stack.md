# SC-001: Tech Stack

## 1. Purpose & Background
- **Purpose**: To ensure consistency across the development and operation of the Count system by standardizing the technical stack. This facilitates easier maintenance, scalability, and integration with established infrastructure.
- **Application Scope**: All backend services, frontend user interfaces, and deployment infrastructure of the Count system.

## 2. Constraint Content (Concrete Conditions)
- **Backend**: Use **Go (Golang)** (stable latest or 1.21+) for all server-side logic and API development. Go's performance and concurrency models are critical for the system's scalability.
- **Frontend**: Use **HTMX** (latest stable) for dynamic content updates and user interactions. This approach minimizes complex JavaScript frameworks and aligns with a simpler, server-driven UI architecture.
- **Infrastructure/Orchestration**: Deploy and manage services using **Kubernetes (K8S)**. All components must be containerized (Docker/OCI compliant) and orchestrated via K8S for reliability and automated scaling.

## 3. Verification Method
- **Verification Method**:
    - **Code Review**: Check if the language used is Go and HTMX attributes are utilized in templates.
    - **Build Pipeline**: Verify that the CI/CD pipelines use Go build tools and produce K8S-ready deployment manifests (e.g., Helm charts or Kustomize).
    - **Deployment Check**: Inspect the running environment to ensure services are operating within a Kubernetes cluster.
- **On Violation**: Any component not adhering to this stack must be refactored to comply with these standards unless an explicit architectural exception is granted.

## 4. References
- **Related URD**: System technical constraints
- **Related FR/QR**: 
    - [QR-001: Performance and Reliability](QR-001-performance-and-reliability.md)
    - [FR-001: Count Item Management](FR-001-count-item-management.md)
