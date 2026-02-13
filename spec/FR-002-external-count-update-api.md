# FR-002: External Count Update API

## 1. Feature Description
- **Purpose**: To provide a high-performance, atomic API for external sources to modify count values. This aligns with [UR-05] and ensures that counting events from various external systems are accurately recorded.
- **Scope**: 
    - API endpoints for increasing, decreasing, and resetting count values.
    - Ensuring atomicity of operations to prevent data corruption during concurrent access.
    - Supporting high-frequency requests as specified in performance requirements.
- **References**: 
    - [UR-05]: External API for count updates.
    - [QR-001]: Performance and Reliability targets (P99 < 200ms, 100% Consistency).
    - [SC-001]: Tech Stack (Go-based implementation for performance).

## 2. User Stories
- **FR-002-01**: As an external system, I want to increase a specific count item's value so that I can record incrementing events in real-time.
- **FR-002-02**: As an external system, I want to decrease a specific count item's value so that I can record decrementing events.
- **FR-002-03**: As an external system, I want to reset a specific count item's value so that I can re-initialize the counter when necessary.

## 3. Acceptance Criteria (Gherkin)

### FR-002-01: Increase Count
- **FR-002-01-01 (Successful Increase)**:
    - **Given** an existing count item with ID `count-123` and current value `10`.
    - **When** I send a POST request to `/api/v1/counts/count-123/increase`.
    - **Then** the system returns status code `200 OK`.
    - **And** the response body contains the updated count value `11`.
    - **And** the stored value of `count-123` is `11`.

- **FR-002-01-02 (Atomic Concurrent Increase)**:
    - **Given** an existing count item with ID `count-123` and current value `0`.
    - **When** 100 concurrent POST requests are sent to `/api/v1/counts/count-123/increase`.
    - **Then** all requests return status code `200 OK`.
    - **And** the final stored value of `count-123` is exactly `100`.

- **FR-002-01-03 (Non-existent Item)**:
    - **Given** no count item exists with ID `invalid-id`.
    - **When** I send a POST request to `/api/v1/counts/invalid-id/increase`.
    - **Then** the system returns status code `404 Not Found`.

### FR-002-02: Decrease Count
- **FR-002-02-01 (Successful Decrease)**:
    - **Given** an existing count item with ID `count-123` and current value `10`.
    - **When** I send a POST request to `/api/v1/counts/count-123/decrease`.
    - **Then** the system returns status code `200 OK`.
    - **And** the response body contains the updated count value `9`.
    - **And** the stored value of `count-123` is `9`.

- **FR-002-02-02 (Decrease below zero - Optional/Policy)**:
    - **Given** an existing count item with ID `count-123` and current value `0`.
    - **When** I send a POST request to `/api/v1/counts/count-123/decrease`.
    - **Then** the system returns status code `200 OK`.
    - **And** the updated count value is `-1` (unless a non-negative constraint is explicitly added later).

- **FR-002-02-03 (Non-existent Item)**:
    - **Given** no count item exists with ID `invalid-id`.
    - **When** I send a POST request to `/api/v1/counts/invalid-id/decrease`.
    - **Then** the system returns status code `404 Not Found`.

### FR-002-03: Reset Count
- **FR-002-03-01 (Successful Reset)**:
    - **Given** an existing count item with ID `count-123` and current value `50`.
    - **When** I send a POST request to `/api/v1/counts/count-123/reset`.
    - **Then** the system returns status code `200 OK`.
    - **And** the response body contains the updated count value `0`.
    - **And** the stored value of `count-123` is `0`.

- **FR-002-03-02 (Non-existent Item)**:
    - **Given** no count item exists with ID `invalid-id`.
    - **When** I send a POST request to `/api/v1/counts/invalid-id/reset`.
    - **Then** the system returns status code `404 Not Found`.
