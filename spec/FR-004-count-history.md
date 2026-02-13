# FR-004: Count Change History Logging and Inquiry

## 1. Feature Description
- **Purpose**: This feature ensures that every modification to a count value is logged for auditability and analysis, and provides mechanisms to retrieve this history. This aligns with [UR-07].
- **Scope**:
  - Automatic logging of every count update (increase, decrease, reset).
  - Recording the source of the update, the timestamp, and the amount of change.
  - Querying the history of a specific count item.
  - Support for both JSON (API) and HTML (UI) output formats.
- **References**: [UR-07], [QR-001], [SC-001], [FR-002]

## 2. User Stories
- **FR-004-01**: As a system administrator, I want every count update event to be automatically logged with its source, timestamp, and change amount so that I have a reliable audit trail. (UR-07)
- **FR-004-02**: As a user or analyst, I want to query the change history for a specific count item so that I can perform time-series analysis or verify past updates. (UR-07)

## 3. Acceptance Criteria (Gherkin)

### FR-004-01: Automatic Logging
- **FR-004-01-01 (Log on Increase)**:
  - **Given** an existing count item with ID `count-123`
  - **When** an external source `source-A` sends an increase request to `/api/v1/counts/count-123/increase`
  - **Then** the system records a log entry with `itemId: "count-123"`, `type: "increase"`, `change: 1`, `source: "source-A"`, and a current `timestamp`
  - **And** the log entry is persistent in the database

- **FR-004-01-02 (Log on Decrease)**:
  - **Given** an existing count item with ID `count-123`
  - **When** an external source `source-B` sends a decrease request to `/api/v1/counts/count-123/decrease`
  - **Then** the system records a log entry with `itemId: "count-123"`, `type: "decrease"`, `change: -1`, `source: "source-B"`, and a current `timestamp`

- **FR-004-01-03 (Log on Reset)**:
  - **Given** an existing count item with ID `count-123` and current value `10`
  - **When** an external source `source-C` sends a reset request to `/api/v1/counts/count-123/reset`
  - **Then** the system records a log entry with `itemId: "count-123"`, `type: "reset"`, `change: -10`, `source: "source-C"`, and a current `timestamp`

### FR-004-02: History Inquiry
- **FR-004-02-01 (API Query Success - JSON)**:
  - **Given** a count item `count-123` has multiple log entries
  - **When** an external system sends a GET request to `/api/v1/counts/count-123/history` with header `Accept: application/json`
  - **Then** the system returns status code `200 OK`
  - **And** the response body is a JSON array of log objects, each containing `timestamp`, `type`, `change`, and `source`
  - **And** the logs are ordered by timestamp (descending)

- **FR-004-02-02 (UI Query Success - HTML)**:
  - **Given** a count item `count-123` has multiple log entries
  - **When** a user requests the history UI fragment from `/ui/counts/count-123/history`
  - **Then** the system returns status code `200 OK`
  - **And** the response is an HTML fragment containing a table or list of the change history
  - **And** each history record displays the timestamp, operation type, amount, and source

- **FR-004-02-03 (Query Non-existent Item)**:
  - **Given** no count item exists with ID `invalid-id`
  - **When** a request is sent to retrieve history for `invalid-id` (via `/api/v1/counts/invalid-id/history` or `/ui/counts/invalid-id/history`)
  - **Then** the system returns status code `404 Not Found`
  - **And** the response contains an error message in the appropriate format (JSON for API, HTML for UI)
