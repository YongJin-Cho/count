# Consolidated UI-API List

| ID | Method | Path | Trigger | Target | Swap | Using UI (UI ID) | Description |
|---|---|---|---|---|---|---|---|
| RegisterCountItemAPI | POST | `/ui/count-items` | `click` (on `btn-register`) | `#count-item-list` | `beforeend` | CountItemManagementUI | Registers a new count item and appends the new item fragment to the list. Resets the form on success. |
| ListCountItemsAPI | GET | `/ui/count-items` | `load` (on `#count-item-list`) | `this` | `innerHTML` | CountItemManagementUI | Fetches and displays the initial list of registered count items. |
| GetCountItemValueAPI | GET | `/ui/counts/{id}/value` | `every 30s`, `click` (on `item-value-{id}`) | `this` | `innerHTML` | CountItemManagementUI | Fetches the current value of a specific count item and updates its display. |
| DeleteCountItemAPI | DELETE | `/ui/count-items/{id}` | `click` (on `btn-delete-{id}`) | `closest .count-item-row` | `outerHTML` | CountItemManagementUI | Deletes a count item and removes its row fragment from the list. Requires user confirmation. |
| GetCountValueAPI | GET | `/ui/counts/{id}/value` | `load, every 5s` | `#value-display` | `innerHTML` | CountItemMonitoringUI | Periodically fetches the current value of a count item and updates the display area. |
| CountItemUpdateUI_GetDashboardAPI | GET | `/ui/counts` | `click` (link-back-dashboard, btn-cancel) | `body` | - | CountItemUpdateUI | Returns the Dashboard UI HTML fragment. `hx-push-url="true"` is used for the back link. |
| CountItemUpdateUI_UpdateItemAPI | PUT | `/ui/counts/{count_id}` | `submit` | `body` | `outerHTML` | CountItemUpdateUI | Updates the count item (name, description). Uses `#update-progress` indicator. On success, returns Dashboard fragment. On validation failure, returns `CountItemUpdateUI` fragment with errors in `validation-error-area`. |
| GetCountItemHistoryAPI | GET | `/ui/counts/{id}/history` | `load` | `#history-table-container` | `innerHTML` | CountItemHistoryUI | Returns an HTML fragment containing either a table (`history-table`) with chronological change logs (timestamp, operation, amount, source) or an empty state message (`empty-history`) if no records exist. |
