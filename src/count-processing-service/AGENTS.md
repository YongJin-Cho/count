# Module Specification: count-processing-service

## 1. Overview
- **Role**: Service for processing count value operations. It handles high-frequency updates from external systems, provides internal APIs for lifecycle management, and maintains an audit trail of all changes.
- **Build Output**: image (count-processing-service)

## 2. Module Structure
- **Architecture Pattern**: Hexagonal Architecture (Ports and Adapters)
- **Structure**:
  - **Domain Layer** (Core): Business logic for count value management, domain entities (CountValue, CountLog), use cases (Initialize, Get, Update, Delete, Increase, Decrease, Reset, GetHistory).
  - **Ports** (Interfaces): 
    - **Inbound Ports**: Interface for `InternalCountValueAPI`, `ExternalCountUpdateAPI`, `CountValueAPI`, and `CountHistoryAPI`.
    - **Outbound Ports**: Interface for `CountValueRepository` and `CountHistoryRepository` (PostgreSQL).
  - **Adapters** (Implementations):
    - **Inbound Adapters**: HTTP Gin/Fiber handlers for Internal and External APIs.
    - **Outbound Adapters**: GORM/sqlx repository for PostgreSQL.
- **Testing Strategy**: 
  - **Adapter Mocking**: Mock the PostgreSQL repository for unit tests.
  - **Port Testing**: Test domain logic through the Inbound API ports.
  - **Integration Testing**: Verify repository implementation against a test database, ensuring atomic SQL operations and correct transaction handling for logging.

## 3. Providing Interfaces
### External API - `ExternalCountUpdateAPI`
- **POST /api/v1/counts/{itemId}/increase**: Increase count value.
  - **Request Body**: `{ "amount": integer, "source": string }`
  - **Logic**: Atomically increment `current_value`. **Automatically log event** to `count_history` (FR-004-01-01).
  - **Response**: `200 OK` with `{ "itemId": string, "value": integer }`.
- **POST /api/v1/counts/{itemId}/decrease**: Decrease count value.
  - **Request Body**: `{ "amount": integer, "source": string }`
  - **Logic**: Atomically decrement `current_value`. **Automatically log event** to `count_history` (FR-004-01-02).
  - **Response**: `200 OK` with `{ "itemId": string, "value": integer }`.
- **POST /api/v1/counts/{itemId}/reset**: Reset count value to zero.
  - **Request Body**: `{ "source": string }`
  - **Logic**: Atomically set `current_value` to 0. **Automatically log event** to `count_history` recording the difference as change amount (FR-004-01-03).
  - **Response**: `200 OK` with `{ "itemId": string, "value": 0 }`.

### External API - `CountValueAPI` & `CountHistoryAPI`
- **GET /api/v1/counts/{id}/value**: Get the current value of a specific count.
- **GET /api/v1/counts/values**: Get current values for all count items.
- **GET /api/v1/counts/{id}/history**: Retrieve count change history (FR-004-02-01).
  - **Logic**: Fetch logs from `count_history` table for `id`, ordered by `timestamp` DESC.
  - **Response**: `200 OK` with JSON array of `{ timestamp, type, change, source }`.
  - **Error**: `404 Not Found` if item does not exist.

### Internal API - `InternalCountValueAPI`
- **POST /api/v1/internal/counts**: Initialize count value.
- **GET /api/v1/internal/counts**: Retrieve multiple count values.
- **DELETE /api/v1/internal/counts/{itemId}**: Delete count value record.

## 4. Functional Requirements
- **Atomic Operations & Logging (FR-004-01)**:
  - Updates and history logging MUST be part of the same transaction to ensure consistency.
  - Every `increase`, `decrease`, or `reset` must insert a row into `count_history`.
- **High-Frequency Performance**: Target 10,000 RPS. Logging must be optimized (e.g., efficient indexing on `item_id` and `timestamp`).
- **Persistence**:
  - **Database Schema**:
    - Table `count_values`: `item_id` (PK), `current_value`, `last_updated_at`.
    - Table `count_history` (FR-004-01): `id` (UUID PK), `item_id` (FK), `operation_type` (increase/decrease/reset), `change_amount` (Integer), `source` (String), `timestamp` (Timestamp).
- **Retrieval Logic**: `GET /api/v1/counts/{id}/history` must return logs in descending order of `timestamp`.

## 5. Dependencies
- **Reference Modules**: None
- **Technologies Used**: Go, PostgreSQL.

## 6. Acceptance Tests
- [ ] `POST /api/v1/counts/{itemId}/increase` returns 200 and inserts a record into `count_history` with type `increase` and correct amount.
- [ ] `POST /api/v1/counts/{itemId}/decrease` returns 200 and inserts a record into `count_history` with type `decrease` and negative amount.
- [ ] `POST /api/v1/counts/{itemId}/reset` returns 200 and inserts a record into `count_history` with type `reset` and the diff amount.
- [ ] `GET /api/v1/counts/{id}/history` returns 200 and a JSON array of logs ordered by newest first (FR-004-02-01).
- [ ] `GET /api/v1/counts/{invalid-id}/history` returns 404 (FR-004-02-03).
- [ ] High concurrency: 100 `increase` requests result in 100 log entries in `count_history`.
