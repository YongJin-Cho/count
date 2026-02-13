# Module Specification: count-processing-service

## 1. Overview
- **Role**: Service for processing count value operations. It handles high-frequency updates and provides internal APIs for count value lifecycle management.
- **Build Output**: image (count-processing-service)

## 2. Module Structure
- **Architecture Pattern**: Hexagonal Architecture (Ports and Adapters)
- **Structure**:
  - **Domain Layer** (Core): Business logic for count value management, domain entities (CountValue), use cases (Initialize, Get, Update, Delete).
  - **Ports** (Interfaces): 
    - **Inbound Ports**: Interface for InternalCountValueAPI.
    - **Outbound Ports**: Interface for CountValueRepository (PostgreSQL).
  - **Adapters** (Implementations):
    - **Inbound Adapters**: HTTP Gin/Fiber handlers for Internal API.
    - **Outbound Adapters**: GORM/sqlx repository for PostgreSQL.
- **Testing Strategy**: 
  - **Adapter Mocking**: Mock the PostgreSQL repository for unit tests.
  - **Port Testing**: Test domain logic through the Internal API port.
  - **Integration Testing**: Verify repository implementation against a test database.

## 3. Providing Interfaces
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
- [ ] `GET /api/v1/internal/counts?itemIds=A&itemIds=B` returns current values for both.
- [ ] `DELETE /api/v1/internal/counts/{itemId}` removes the record from DB.
