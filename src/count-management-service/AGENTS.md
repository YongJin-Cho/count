# Module Specification: count-management-service

## 1. Overview
- **Role**: Service for managing count item metadata and providing the management UI. It acts as the orchestrator for metadata and ensures the processing service is synchronized when items are created or deleted.
- **Build Output**: image (count-management-service)

## 2. Module Structure
- **Architecture Pattern**: Hexagonal Architecture (Ports and Adapters)
- **Structure**:
  - **Domain Layer** (Core): Business logic for count item metadata, domain entities (CountItem), use cases (Register, List, Update, Delete, ViewHistory).
  - **Ports** (Interfaces): 
    - **Inbound Ports**: Interfaces for CountItemAPI (External), CountItemManagementUI (UI-API), and CountItemHistoryUI (UI-API).
    - **Outbound Ports**: Interface for CountItemRepository (PostgreSQL), ValueServiceClient (InternalCountValueAPI), and HistoryServiceClient (CountHistoryAPI).
  - **Adapters** (Implementations):
    - **Inbound Adapters**: HTTP Gin/Fiber handlers for JSON and HTMX fragments.
    - **Outbound Adapters**: GORM/sqlx repository for PostgreSQL, HTTP client for count-processing-service.
- **Testing Strategy**: 
  - **Adapter Mocking**: Mock the PostgreSQL repository and the clients (ValueServiceClient, HistoryServiceClient) calling the processing service.
  - **Port Testing**: Test domain logic through ports using mocked adapters.
  - **Integration Testing**: Use mocked adapters to test the integration between domain logic and adapters.

## 3. Providing Interfaces
### UI-API (HTMX) - `CountItemManagementUI`, `CountItemMonitoringUI` & `CountItemHistoryUI`
- **POST /ui/count-items**: Registers a new count item.
  - **Request**: `application/x-www-form-urlencoded` { `name`, `description` }.
  - **Response (201)**: HTML `<tr>` fragment for the new item.
  - **Response (400/409)**: HTML `<div>` fragment with error message.
- **GET /ui/count-items**: Lists all registered count items.
  - **Logic**: Fetch all metadata from PostgreSQL, then call `GET /api/v1/internal/counts?itemIds=...` on processing service to get current values.
  - **Response (200)**: HTML `<tr>` list with current values, or empty state message.
- **GET /ui/counts/{id}/value**: Returns the current value fragment for a specific item.
  - **Logic**: Call `GET /api/v1/internal/counts/{id}` on processing service.
  - **Response (200)**: HTML fragment containing only the numeric value (e.g., `42`).
- **GET /ui/counts/{id}/history**: Returns the change history HTML fragment (FR-004-02-02).
  - **Logic**: Call `GET /api/v1/counts/{id}/history` on processing service. Render the received JSON logs into an HTML table (or empty state).
  - **Response (200)**: HTML fragment (`<table id="history-table">` or `<div id="empty-history">`).
  - **HTMX Spec**: Trigger: `load`, Target: `#history-table-container`, Swap: `innerHTML`.
- **DELETE /ui/count-items/{id}**: Deletes a count item.
  - **Response (200)**: Empty string (removes row via HTMX `outerHTML` swap).
- **PUT /ui/counts/{count_id}**: Updates metadata.
  - **Response (200)**: Redirect (HX-Redirect) or Dashboard fragment.

### External API - `CountItemAPI`
- **GET /api/v1/count-items**: List all items (JSON). Returns array of `{id, name, description}`.
- **POST /api/v1/count-items**: Register item (JSON). Request: `{name, description}`.
- **PUT /api/v1/count-items/{id}**: Update item (JSON). Request: `{name, description}`.
- **DELETE /api/v1/count-items/{id}**: Delete item (JSON).

## 4. Functional Requirements
- **FR-001-01 (Register)**: Validate `name`, save to `count_items`, call `POST /api/v1/internal/counts` on processing service.
- **FR-001-04 (Delete)**: Delete from `count_items`, call `DELETE /api/v1/internal/counts/{id}` on processing service.
- **FR-003-01/02 (Value Retrieval)**: Handlers call `InternalCountValueAPI` on `count-processing-service` to get latest value(s).
- **FR-004-02 (History Inquiry UI)**:
  - Handler for `/ui/counts/{id}/history` must call `CountHistoryAPI` on `count-processing-service`.
  - Format the response into the structure defined in `src/interface/ui/CountItemHistoryUI.json`.
  - If processing service returns 404, return error fragment.
- **Business Rule**: Metadata and initial value must be consistent. Ensure the processing service is notified of creation/deletion.

## 5. Dependencies
- **Reference Modules**: `src/count-processing-service/`
- **Technologies Used**: Go (Backend), HTMX (Frontend), PostgreSQL.

## 6. Acceptance Tests
- [ ] Returns 400 + error HTML when name is empty in registration form.
- [ ] On successful registration, processing service receives initialization call.
- [ ] `GET /ui/counts/{id}/value` returns HTML fragment with current value.
- [ ] `GET /ui/counts/{id}/history` returns an HTML table containing chronological change logs (timestamp, operation, amount, source) (FR-004-02-02).
- [ ] `GET /ui/counts/{id}/history` returns empty state fragment when no history records exist.
- [ ] `GET /ui/counts/{invalid-id}/history` returns 404 + error HTML fragment (FR-004-02-03).
