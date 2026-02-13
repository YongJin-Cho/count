# Module Specification: count-management-service

## 1. Overview
- **Role**: Service for managing count item metadata and providing the management UI. It acts as the orchestrator for metadata and ensures the processing service is synchronized when items are created or deleted.
- **Build Output**: image (count-management-service)

## 2. Module Structure
- **Architecture Pattern**: Hexagonal Architecture (Ports and Adapters)
- **Structure**:
  - **Domain Layer** (Core): Business logic for count item metadata, domain entities (CountItem), use cases (Register, List, Update, Delete).
  - **Ports** (Interfaces): 
    - **Inbound Ports**: Interfaces for CountItemAPI (External) and CountItemManagementUI (UI-API).
    - **Outbound Ports**: Interface for CountItemRepository (PostgreSQL) and ValueServiceClient (InternalCountValueAPI).
  - **Adapters** (Implementations):
    - **Inbound Adapters**: HTTP Gin/Fiber handlers for JSON and HTMX fragments.
    - **Outbound Adapters**: GORM/sqlx repository for PostgreSQL, HTTP client for count-processing-service.
- **Testing Strategy**: 
  - **Adapter Mocking**: Mock the PostgreSQL repository and the ValueServiceClient (which calls the processing service).
  - **Port Testing**: Test domain logic through ports using mocked adapters.
  - **Integration Testing**: Use mocked adapters to test the integration between domain logic and adapters.

## 3. Providing Interfaces
### UI-API (HTMX) - `CountItemManagementUI` & `CountItemMonitoringUI`
- **POST /ui/count-items**: Registers a new count item.
  - **Request**: `application/x-www-form-urlencoded` { `name`, `description` }.
  - **Response (201)**: HTML `<tr>` fragment for the new item.
    - Example: `<tr class="count-item-row" id="count-item-123"><td>Inventory</td><td>Stock count for warehouse A</td><td><span id="item-value-123">0</span></td><td><button id="btn-delete-123">Delete</button></td></tr>`
  - **Response (400/409)**: HTML `<div>` fragment with error message.
- **GET /ui/count-items**: Lists all registered count items.
  - **Logic**: Fetch all metadata from PostgreSQL, then call `GET /api/v1/internal/counts?itemIds=...` on processing service to get current values.
  - **Response (200)**: HTML `<tr>` list with current values, or empty state message.
- **GET /ui/counts/{id}/value**: Returns the current value fragment for a specific item.
  - **Logic**: Call `GET /api/v1/internal/counts/{id}` on processing service.
  - **Response (200)**: HTML fragment containing only the numeric value (e.g., `42`).
- **DELETE /ui/count-items/{id}**: Deletes a count item.
  - **Response (200)**: Empty string (removes row via HTMX `outerHTML` swap).
- **PUT /ui/counts/{count_id}**: Updates metadata.
  - **Response (200)**: Redirect (HX-Redirect) or Dashboard fragment.

### External API - `CountItemAPI`
- **GET /api/v1/count-items**: List all items (JSON). Returns array of `{id, name, description}`.
  - Note: This endpoint returns metadata only. Use `CountValueAPI` on processing service for values.
- **POST /api/v1/count-items**: Register item (JSON). Request: `{name, description}`.
- **PUT /api/v1/count-items/{id}**: Update item (JSON). Request: `{name, description}`.
- **DELETE /api/v1/count-items/{id}**: Delete item (JSON).

## 4. Functional Requirements
- **FR-001-01 (Register)**:
  - Validate `name` is not empty.
  - Save `id`, `name`, `description` to PostgreSQL `count_items` table.
  - **Call Processing Service**: `POST /api/v1/internal/counts` with `{ "itemId": id, "initialValue": 0 }`.
  - On success, return HTMX fragment or JSON object.
- **FR-001-04 (Delete)**:
  - Delete from `count_items` table.
  - **Call Processing Service**: `DELETE /api/v1/internal/counts/{id}`.
- **FR-003-01/02 (Value Retrieval)**:
  - **Single Value**: Handlers for `/ui/counts/{id}/value` call `InternalCountValueAPI` on `count-processing-service` to get the latest value and return it as an HTML fragment.
  - **Bulk Value**: The main list handler `/ui/count-items` performs a bulk call to `InternalCountValueAPI` using all registered item IDs to efficiently display the list with values.
- **Business Rule**: Metadata and initial value must be consistent. Ensure the processing service is notified of creation/deletion.

## 5. Dependencies
- **Reference Modules**: `src/count-processing-service/`
- **Technologies Used**: Go (Backend), HTMX (Frontend), PostgreSQL.

## 6. Acceptance Tests
- [ ] Returns 400 + error HTML when name is empty in registration form.
- [ ] Returns 409 + error HTML when duplicate name is submitted.
- [ ] On successful registration, a row is appended to `#count-item-list` and processing service receives initialization call.
- [ ] On deletion, the row is removed from UI and processing service receives deletion call.
- [ ] `GET /ui/count-items` returns a list where each item shows its current value retrieved from the processing service.
- [ ] `GET /ui/counts/{id}/value` returns status `200 OK` and an HTML fragment with the current value.
- [ ] `GET /ui/counts/{invalid-id}/value` returns status `404 Not Found` with an error fragment.
