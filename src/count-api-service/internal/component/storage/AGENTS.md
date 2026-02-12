# Module Specification: storage-module

## 1. Overview
- **Role**: Subscribes to internal count events and persists the data to permanent storage (file-based).
- **Build Output**: Library (component logic).

## 2. Providing Interfaces
- **Event Consumer**:
    - Subscribes to `CountCollectedEvent` on channel `internal.count.collected`.

## 3. Functional Requirements
- **Core Logic**:
    1. **Subscription**: On startup, register a handler for `CountCollectedEvent` via `event-module`.
    2. **Processing**: When an event is received:
        - Extract `external_id`, `count`, and `timestamp`.
        - Format the data for storage (e.g., CSV or line-delimited JSON).
    3. **Persistence**: Write the data to the configured storage path (e.g., `/data/counts.log`).
    4. **Error Handling**: Log errors if file writing fails.

## 4. Dependencies
- **Reference Modules**:
    - `common-layer/model-module` (src/count-api-service/internal/common/model)
    - `common-layer/event-module` (src/count-api-service/internal/common/event)
- **Technologies Used**: Go, OS File System.

## 5. Acceptance Tests
- [ ] Successfully registers as a subscriber to `CountCollectedEvent`.
- [ ] Correctly parses event data into storage format.
- [ ] New line is appended to the storage file for every received event.
- [ ] Data integrity is maintained (written values match event values).
