# Module Specification: collector-module

## 1. Overview
- **Role**: Handles count data collection and query requests. It validates incoming collection data, emits events, and provides an integrated query interface for collected counts.
- **Build Output**: Library (component logic).

## 2. Providing Interfaces
- **API**:
    - `POST /api/v1/collect` (operationId: `collectCount`)
        - Request Body: JSON with `external_id` (string) and `count` (integer).
        - Responses: 200 OK, 400 Bad Request, 401 Unauthorized.
        - Ref: `src/interface/api/CountCollectAPI.json`
    - `GET /api/v1/counts` (operationId: `getCounts`)
        - Query Parameters:
            - `external_id` (string, optional): Filter by source identifier.
            - `limit` (integer, default: 10, min: 1): Max items to return.
            - `offset` (integer, default: 0, min: 0): Items to skip.
        - Responses:
            - 200 OK: `CountListResponse` { `total_count` (int), `counts` (array of `CountItem`) }
            - 401 Unauthorized: Invalid or missing token.
            - 500 Internal Server Error: Unexpected error.
        - Ref: `src/interface/api/CountQueryAPI.json`

## 3. Functional Requirements
- **User Story**:
    - FR-01 (API Collection): Collect count data from external systems.
    - FR-02 (Data Validation): Validate collected data formats and values.
    - FR-03 (Integrated Query): Provide paginated and filtered view of all collected data.
- **Core Logic**:
    ### Collection (POST /api/v1/collect)
    1. **Authentication**: Call `auth-module` to validate Bearer Token.
    2. **Validation**: Check `external_id` (required) and `count` (required, non-negative).
    3. **Event Emission**: Create and publish `CountCollected` event via `event-module`.
    ### Query (GET /api/v1/counts)
    1. **Authentication**: Call `auth-module` to validate Bearer Token. Return 401 if invalid.
    2. **Parameter Parsing**:
        - Parse `limit` and `offset` from query string. Apply defaults if missing.
        - Validate `limit` >= 1 and `offset` >= 0.
    3. **Data Retrieval**:
        - Use `storage-module` to fetch `total_count` and paginated `counts`.
        - If `external_id` is provided, filter results by that ID.
    4. **Response Construction**:
        - Return 200 OK with the retrieved data.
        - If no records found for `external_id`, return `total_count: 0` and empty `counts` list.

## 4. Dependencies
- **Reference Modules**:
    - `common-layer/model-module` (src/count-api-service/internal/common/model)
    - `common-layer/auth-module` (src/count-api-service/internal/common/auth)
    - `common-layer/event-module` (src/count-api-service/internal/common/event)
    - `component-layer/storage-module` (src/count-api-service/internal/component/storage)
- **Technologies Used**: Go, Net/HTTP.

## 5. Acceptance Tests
- [x] Returns 401 when Authorization header is missing or invalid.
- [x] Returns 400 when collection data is invalid (missing fields, negative count).
- [x] Successfully emits `CountCollected` event on valid collection.
- [x] GET /api/v1/counts returns 200 with `total_count` and `counts` array.
- [x] GET /api/v1/counts correctly applies `limit` and `offset` parameters.
- [x] GET /api/v1/counts filters by `external_id` when provided.
- [x] GET /api/v1/counts returns empty list and 0 total_count for non-existent `external_id`.
