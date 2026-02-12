# Module Specification: collector-module

## 1. Overview
- **Role**: Receives count data from external systems via HTTP API, validates it, and emits internal events.
- **Build Output**: Library (component logic).

## 2. Providing Interfaces
- **API**:
    - `POST /api/v1/collect` (operationId: `collectCount`)
    - Request Body: JSON with `external_id` (string) and `count` (integer).
    - Responses:
        - 200 OK: `{ "status": "success" }`
        - 400 Bad Request: `{ "error": "Validation failed", "message": "..." }`
        - 401 Unauthorized: Authentication failed.
    - Ref: `src/interface/api/CountCollectAPI.json`

## 3. Functional Requirements
- **User Story**: FR-01 (API Collection), FR-02 (Data Validation).
- **Core Logic**:
    1. **Authentication**: Call `auth-module` to validate Bearer Token. Return 401 if invalid.
    2. **Parsing**: Parse JSON body. Return 400 if JSON is malformed (FR-02-03-01).
    3. **Validation**:
        - Check if `external_id` is present. Return 400 with "missing external_id" if not (FR-02-01-01).
        - Check if `count` is present. Return 400 with "missing count" if not (FR-02-01-02).
        - Check if `count` >= 0. Return 400 with "invalid count value" if negative (FR-02-01-03).
    4. **Event Emission**:
        - Create a `CountCollected` event using `event-module`.
        - Publish the event to the internal bus.
    5. **Response**: Return 200 OK with success status.

## 4. Dependencies
- **Reference Modules**:
    - `common-layer/model-module` (src/count-api-service/internal/common/model)
    - `common-layer/auth-module` (src/count-api-service/internal/common/auth)
    - `common-layer/event-module` (src/count-api-service/internal/common/event)
- **Technologies Used**: Go, Net/HTTP.

## 5. Acceptance Tests
- [ ] Returns 401 when Authorization header is missing.
- [ ] Returns 400 + "missing external_id" when `external_id` is empty.
- [ ] Returns 400 + "missing count" when `count` is missing.
- [ ] Returns 400 + "invalid count value" when `count` is -1.
- [ ] Returns 200 + "status": "success" and emits event when input is valid.
