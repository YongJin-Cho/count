# Module Specification: count-api-service-module

## 1. Overview
- **Role**: The main entry point for the Count API Service. It orchestrates all components (collector, storage) and starts the server.
- **Build Output**: image (count-api-service).

## 2. Providing Interfaces
- **Entry Point**: Main function (`main.go`).
- **HTTP Server**: Listens on port (default 8080) for the `collector-module` endpoints.

## 3. Functional Requirements
- **Core Logic**:
    1. **Initialization**:
        - Initialize `event-module` (internal bus).
        - Initialize `storage-module` and register its subscription to the event bus.
        - Initialize `auth-module` (auth middleware).
        - Initialize `collector-module` with its handlers.
    2. **Orchestration**:
        - Setup the HTTP server routes using handlers from `collector-module`.
        - Inject `auth-module` middleware into the collector routes.
    3. **Lifecycle**:
        - Handle graceful shutdown (wait for pending events/requests to finish).
        - Log service startup and ready status.

## 4. Dependencies
- **Reference Modules**:
    - `component-layer/collector-module` (src/count-api-service/internal/component/collector)
    - `component-layer/storage-module` (src/count-api-service/internal/component/storage)
- **Technologies Used**: Go, Docker (for image build).

## 5. Acceptance Tests
- [ ] Application starts without errors.
- [ ] HTTP server is reachable on port 8080.
- [ ] Integrating collector and storage via event bus works as expected (End-to-end flow from API to File).
