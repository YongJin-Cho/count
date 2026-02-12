# Module Specification: storage-module

## 1. Overview
- **Role**: Handles persistence and retrieval of count data. It subscribes to count collection events for storage and provides repository interfaces for retrieving data for queries.
- **Build Output**: Library (component logic).

## 2. Providing Interfaces
- **Event Consumer**:
    - Subscribes to `CountCollectedEvent` via `event-module`.
- **Repository (Internal API)**:
    - `FindAll(filter string, limit int, offset int) ([]CountItem, error)`: Retrieves paginated count records, optionally filtered by `external_id`.
    - `CountTotal(filter string) (int, error)`: Returns the total number of records matching the filter.

## 3. Functional Requirements
- **Core Logic**:
    ### Persistence (Event-driven)
    1. **Subscription**: On startup, register a handler for `CountCollectedEvent`.
    2. **Processing**: When an event is received, format it for storage (e.g., CSV or line-delimited JSON).
    3. **Writing**: Append the formatted data to the storage file (e.g., `/data/counts.log`).
    ### Retrieval (Repository)
    1. **FindAll**:
        - Open the storage file.
        - Scan records and filter by `external_id` if provided.
        - Apply `offset` and `limit` to the result set.
        - Return a slice of `CountItem` models.
    2. **CountTotal**:
        - Scan the storage file.
        - Count all records that match the `external_id` filter (or all records if no filter).
        - Return the total count.

## 4. Dependencies
- **Reference Modules**:
    - `common-layer/model-module` (src/common-layer/model-module)
    - `common-layer/event-module` (src/common-layer/event-module)
- **Technologies Used**: Go, OS File System.

## 5. Acceptance Tests
- [x] Successfully registers as a subscriber to `CountCollectedEvent`.
- [x] Correctly appends new count data to storage file.
- [x] `FindAll` returns correct subset of data based on `limit` and `offset`.
- [x] `FindAll` returns only records matching `external_id` when filter is applied.
- [x] `CountTotal` returns accurate total count matching the filter, ignoring pagination parameters.
- [x] Returns empty result and 0 count if storage file is empty or no matches found.
