# FR-003: Count Value Retrieval

## 1. Feature Description
- **Purpose**: To allow users and external systems to retrieve the current value of count items, providing visibility into the tracked metrics. This aligns with [UR-06].
- **Scope**: 
    - API endpoints for retrieving the value of a specific count item and all count items in JSON format.
    - UI support for displaying the current value of count items using HTML fragments (compatible with HTMX).
    - Support for both specific counter retrieval and bulk retrieval.
- **References**: [UR-06], [QR-001], [SC-001]

## 2. User Stories
- **FR-003-01**: As a user or external system, I want to retrieve the current value of a specific count item so that I can monitor its status. (UR-06)
- **FR-003-02**: As a user or external system, I want to retrieve the current values of all count items so that I can get a comprehensive overview of all tracked metrics. (UR-06)

## 3. Acceptance Criteria (Gherkin)

### FR-003-01: Specific Count Item Retrieval
- **FR-003-01-01 (API Success)**:
    - **Given** an existing count item with ID `count-123` and current value `42`
    - **When** an external system sends a GET request to `/api/v1/counts/count-123/value`
    - **Then** the system returns status code `200 OK`
    - **And** the response body is a JSON object containing the `itemId` and `currentValue`
    - **And** the `currentValue` in the response is `42`

- **FR-003-01-02 (UI Success)**:
    - **Given** an existing count item with ID `count-123` and current value `42`
    - **When** a user's browser requests the value UI fragment from `/ui/counts/count-123/value`
    - **Then** the system returns status code `200 OK`
    - **And** the response is an HTML fragment containing the current value `42`

- **FR-003-01-03 (Non-existent Item)**:
    - **Given** no count item exists with ID `invalid-id`
    - **When** a request is sent to retrieve its value (via `/api/v1/counts/invalid-id/value` or `/ui/counts/invalid-id/value`)
    - **Then** the system returns status code `404 Not Found`
    - **And** the response contains an error message in the appropriate format (JSON for API, HTML for UI)

### FR-003-02: All Count Items Retrieval
- **FR-003-02-01 (API Success)**:
    - **Given** multiple count items exist (e.g., `count-1` with value `10`, `count-2` with value `20`)
    - **When** an external system sends a GET request to `/api/v1/counts/values`
    - **Then** the system returns status code `200 OK`
    - **And** the response body is a JSON array of objects, each containing `itemId` and `currentValue`
    - **And** the array contains the values for both `count-1` and `count-2`

- **FR-003-02-02 (UI Success)**:
    - **Given** multiple count items exist
    - **When** a user's browser requests the all-values UI fragment from `/ui/counts/values`
    - **Then** the system returns status code `200 OK`
    - **And** the response is an HTML fragment representing a list or table of all count items with their current values

- **FR-003-02-03 (Empty List Success)**:
    - **Given** no count items have been registered
    - **When** a request is sent to retrieve all count values
    - **Then** the system returns status code `200 OK`
    - **And** the response returns an empty state representation (JSON empty array `[]` or HTML empty state fragment)
