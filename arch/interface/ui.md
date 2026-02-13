# UI Interface Design

This document identifies the UI screens and user flows for the Count system, including count item management and value retrieval.

## UI List

| UI ID | Screen Name | Key Features | Related User Story |
|-------|-------------|--------------|-------------------|
| `CountItemManagementUI` | Count Item Management Dashboard | - List all registered count items (Name, Description, **Current Value**)<br>- **Display real-time/updated values for all items**<br>- Display empty state when no items exist<br>- Form to register a new count item (Unique Name, Optional Description)<br>- Inline deletion of count items via HTMX | FR-001-01, FR-001-02, FR-001-04, **FR-003-02** |
| `CountItemUpdateUI` | Count Item Edit Screen | - Form pre-filled with existing count item details<br>- Update name and description<br>- Validation and error message display | FR-001-03 |
| `CountItemMonitoringUI` | Specific Count Monitoring UI | - Display current value of a specific count item<br>- Real-time value updates via dedicated HTML fragment | FR-003-01 |

## User Flows

### 1. Count Item Registration Flow
Allows users to create a new counter to track.
- **Path**: `CountItemManagementUI` (Enter details) → Submit → `CountItemManagementUI` (Item appears in list)

### 2. Count Item List & Deletion Flow
Allows users to view all counters and remove those no longer needed.
- **Path**: `CountItemManagementUI` (View list) → Trigger Delete → `CountItemManagementUI` (Item removed from list)

### 3. Count Item Update Flow
Allows users to modify the metadata of an existing counter.
- **Path**: `CountItemManagementUI` (Select Edit) → `CountItemUpdateUI` (Modify details) → Submit → `CountItemManagementUI` (Updated details reflected in list)

### 4. Count Value Monitoring Flow
Allows users to monitor the current state of counters in real-time.
- **Path (Dashboard)**: `CountItemManagementUI` (View list with values) → (Auto-refresh) → `CountItemManagementUI` (Updated values reflected)
- **Path (Specific)**: `CountItemManagementUI` (Select Item) → `CountItemMonitoringUI` (View specific real-time value)
