# Module Specification: ui-module

## 1. Overview
- **Role**: Handles UI-related requests and serves HTMX templates for Count Management (FR-04). It provides HTML fragments for dynamic UI updates, supporting listing, creating, and modifying count data.
- **Build Output**: Library (component logic).

## 2. Providing Interfaces
- **UI API (HTML Fragments)**:
    - `GET /ui/counts`: Returns `<tr>` fragments for the count list.
    - `POST /ui/counts`: Creates a new count and returns the new `<tr>` fragment. Trigger: `count-list-changed`.
    - `GET /ui/counts/new`: Returns the "Create New Count" form fragment.
    - `GET /ui/counts/{source_id}`: Returns the view-mode `<tr>` fragment for a specific source.
    - `PUT /ui/counts/{source_id}`: Updates the count value and returns the view-mode `<tr>` fragment.
    - `GET /ui/counts/{source_id}/edit`: Returns the edit-mode form fragment.
    - `POST /ui/counts/{source_id}/increment`: Increments value by 1 and returns updated `<tr>`.
    - `POST /ui/counts/{source_id}/decrement`: Decrements value by 1 (min 0) and returns updated `<tr>`.
- **HTML Fragment Structure**:
    - **Row Fragment (`<tr>`)**:
      ```html
      <tr id="count-row-${source_id}">
        <td>${source_id}</td>
        <td><span class="badge">${current_count}</span></td>
        <td>${last_updated}</td>
        <td>
          <button hx-post="/ui/counts/${source_id}/increment" hx-target="closest tr" hx-swap="outerHTML">+1</button>
          <button hx-post="/ui/counts/${source_id}/decrement" hx-target="closest tr" hx-swap="outerHTML">-1</button>
          <button hx-get="/ui/counts/${source_id}/edit" hx-target="#main-content" hx-push-url="true">수정</button>
        </td>
      </tr>
      ```
    - **Form Fragment**: Use `<form hx-post="/ui/counts" hx-target="#count-list-body" hx-swap="beforeend">` for creation.

## 3. Functional Requirements
- **User Story**: FR-04 (카운트 관리 웹 UI)
    - FR-04-01: Create new count source with Source ID (`[a-z0-9-]+`) and initial value.
    - FR-04-02: Display all registered counts in a table.
    - FR-04-03: Manual increment/decrement and updates.
    - FR-04-04: Error handling with HTMX (e.g., duplicate ID).
- **Core Logic**:
    - **List Rendering**: Call `storage-module` to get all counts and render `<tr>` list.
    - **Creation**: Validate `source_id` regex. Call `collector-module` to create. Return 409 if exists.
    - **Modification**: Call `collector-module` for increment/decrement/update. Return updated row fragment.
    - **Error Handling**: On failure, return HTML fragment with error message (e.g., `<div class="text-red-500">...</div>`).

## 4. Dependencies
- **Reference Modules**:
    - `component-layer/collector-module` (src/component-layer/collector-module)
    - `component-layer/storage-module` (src/component-layer/storage-module)
    - `common-layer/model-module` (src/common-layer/model-module)
    - `common-layer/auth-module` (src/common-layer/auth-module)
- **Technologies Used**: Go, HTMX, HTML Templates.

## 5. Acceptance Tests
- [ ] `GET /ui/counts` returns 200 with `<tr>` elements for all stored counts.
- [ ] `POST /ui/counts` returns 409 fragment when `source_id` is duplicated.
- [ ] `POST /ui/counts` returns 200 with new `<tr>` on success and triggers `count-list-changed`.
- [ ] Clicking "+1" button updates the specific row's count value without page refresh.
- [ ] Edit form submission updates the value and returns to view-mode row.
- [ ] Invalid `source_id` (e.g., containing Uppercase) returns 400 with error fragment.
