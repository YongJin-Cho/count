# Module Specification: storage-module

## 1. Overview
- **Role**: Handles persistence and retrieval of count data. It subscribes to count collection events for storage and provides repository interfaces for retrieving data for queries.
- **Build Output**: Library (component logic).

## 2. Providing Interfaces
- **Event Consumer**:
    - Subscribes to `CountCollectedEvent` via `event-module`.
- **Repository (Internal API)**:
    - `FindAll(filter string, limit int, offset int) ([]CountItem, error)`: Retrieves paginated count records, optionally filtered by `external_id`.
    - `FindById(id string) (*CountItem, error)`: Retrieves a specific count record by its Source ID.
    - `CountTotal(filter string) (int, error)`: Returns the total number of records matching the filter.
    - `Create(item CountItem) error`: Persists a new count record.
    - `UpdateValue(id string, value int) error`: Updates the count value for a specific Source ID.

## 3. Functional Requirements
- **Core Logic**:
    ### Persistence (Event-driven)
    1. **Subscription**: On startup, register a handler for `CountCollectedEvent`.
    2. **Processing**: When an event is received, format it for storage (e.g., CSV or line-delimited JSON).
    3. **Writing**: Append the formatted data to the storage file (e.g., `/data/counts.log`).
    ### Retrieval & Management (Repository)
    1. **FindAll**:
        - Open the storage file.
        - Scan records and filter by `external_id` if provided.
        - Apply `offset` and `limit` to the result set.
        - Return a slice of `CountItem` models.
    2. **FindById**:
        - Scan the storage file for a record with matching `source_id`.
        - Return the first match found or error if not found.
    3. **CountTotal**:
        - Scan the storage file.
        - Count all records that match the `external_id` filter (or all records if no filter).
        - Return the total count.
    4. **Create**:
        - Append a new record to the storage file with the given `source_id` and initial value.
    5. **UpdateValue**:
        - (Implementation specific: e.g., append a new record with the updated value, or rewrite the file).
        - Ensure the latest value is the one retrieved by `FindById` and `FindAll`.

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
