# Module Specification: count-processing-service

## 1. Overview
- **Role**: Service for processing count value operations. It handles high-frequency updates from external systems and provides internal APIs for count value lifecycle management.
- **Build Output**: image (count-processing-service)

## 2. Module Structure
- **Architecture Pattern**: Hexagonal Architecture (Ports and Adapters)
- **Structure**:
  - **Domain Layer** (Core): Business logic for count value management, domain entities (CountValue), use cases (Initialize, Get, Update, Delete, Increase, Decrease, Reset).
  - **Ports** (Interfaces): 
    - **Inbound Ports**: Interface for `InternalCountValueAPI` and `ExternalCountUpdateAPI`.
    - **Outbound Ports**: Interface for `CountValueRepository` (PostgreSQL).
  - **Adapters** (Implementations):
    - **Inbound Adapters**: HTTP Gin/Fiber handlers for Internal and External APIs.
    - **Outbound Adapters**: GORM/sqlx repository for PostgreSQL.
- **Testing Strategy**: 
  - **Adapter Mocking**: Mock the PostgreSQL repository for unit tests.
  - **Port Testing**: Test domain logic through the Inbound API ports.
  - **Integration Testing**: Verify repository implementation against a test database, ensuring atomic SQL operations.

## 3. Providing Interfaces
### External API - `ExternalCountUpdateAPI`
- **POST /api/v1/counts/{itemId}/increase**: Increase count value.
  - **Request Body**: `{ "amount": integer }` (Default: 1)
  - **Logic**: Atomically increment `current_value` by `amount`.
  - **Response**: `200 OK` with `{ "itemId": string, "value": integer }`.
  - **Error**: `404 Not Found` if `itemId` does not exist.
- **POST /api/v1/counts/{itemId}/decrease**: Decrease count value.
  - **Request Body**: `{ "amount": integer }` (Default: 1)
  - **Logic**: Atomically decrement `current_value` by `amount`.
  - **Response**: `200 OK` with `{ "itemId": string, "value": integer }`.
  - **Error**: `404 Not Found` if `itemId` does not exist.
- **POST /api/v1/counts/{itemId}/reset**: Reset count value to zero.
  - **Logic**: Atomically set `current_value` to 0.
  - **Response**: `200 OK` with `{ "itemId": string, "value": 0 }`.
  - **Error**: `404 Not Found` if `itemId` does not exist.

### Internal API - `InternalCountValueAPI`
- **POST /api/v1/internal/counts**: Initialize count value.
  - **Request Body**: `{ "itemId": string, "initialValue": integer }`
  - **Logic**: Create record in `count_values` table. Return 409 if already exists.
- **GET /api/v1/internal/counts**: Retrieve multiple count values.
  - **Query Param**: `itemIds[]` (e.g., `?itemIds=id1&itemIds=id2`)
  - **Response**: List of `{ "itemId", "currentValue", "lastUpdatedAt" }`.
- **GET /api/v1/internal/counts/{itemId}**: Retrieve single value.
- **DELETE /api/v1/internal/counts/{itemId}**: Delete count value record.

## 4. Functional Requirements
- **High-Frequency Performance**:
  - Target: **10,000 RPS** (Requests Per Second).
  - P99 Latency: < 200ms.
- **Atomic Operations**:
  - Updates MUST be atomic at the database level to handle high concurrency.
  - Use SQL: `UPDATE count_values SET current_value = current_value + :amount, last_updated_at = NOW() WHERE item_id = :itemId RETURNING current_value`.
  - Do NOT use "Select-then-Update" in application code.
- **Initialization**: When called (typically by management service), create a new entry in `count_values` table with `current_value = initialValue` (default 0).
- **Integrity**: Ensure `itemId` is unique in `count_values` table.
- **Persistence**: Store all values in PostgreSQL `count_values` table.
  - **Database Schema Requirements**:
    - Table: `count_values`
    - `item_id`: UUID or String, Primary Key
    - `current_value`: Integer, Default 0
    - `last_updated_at`: Timestamp, auto-updated on change.

## 5. Dependencies
- **Reference Modules**: None
- **Technologies Used**: Go, PostgreSQL.

## 6. Acceptance Tests
- [ ] `POST /api/v1/internal/counts` returns 201 and creates record in DB.
- [ ] `POST /api/v1/internal/counts` returns 409 if `itemId` already has a value.
- [ ] `POST /api/v1/counts/{itemId}/increase` returns 200 and correctly incremented value.
- [ ] `POST /api/v1/counts/{itemId}/increase` returns 404 for non-existent `itemId`.
- [ ] 100 concurrent `increase` requests to same `itemId` result in exactly +100 in DB.
- [ ] `POST /api/v1/counts/{itemId}/decrease` returns 200 and correctly decremented value.
- [ ] `POST /api/v1/counts/{itemId}/reset` returns 200 and value becomes 0.
- [ ] `GET /api/v1/internal/counts?itemIds=A&itemIds=B` returns current values for both.
- [ ] `DELETE /api/v1/internal/counts/{itemId}` removes the record from DB.
