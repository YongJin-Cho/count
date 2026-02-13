# Consolidated UI-API List

| ID | Method | Path | Trigger | Target | Swap | Using UI (UI ID) | Description |
|---|---|---|---|---|---|---|---|
| RegisterCountItemAPI | POST | `/ui/count-items` | `click` (on `btn-register`) | `#count-item-list` | `beforeend` | CountItemManagementUI | Registers a new count item and appends the new item fragment to the list. |
| ListCountItemsAPI | GET | `/ui/count-items` | `load` (on `#count-item-list`) | `this` | `innerHTML` | CountItemManagementUI | Fetches and displays the initial list of registered count items. |
| DeleteCountItemAPI | DELETE | `/ui/count-items/{id}` | `click` (on `btn-delete-{id}`) | `closest .count-item-row` | `outerHTML` | CountItemManagementUI | Deletes a count item and removes its row fragment from the list. |
| CountItemUpdateUI_GetDashboardAPI | GET | `/ui/counts` | `click` (link-back-dashboard, btn-cancel) | `body` | - | CountItemUpdateUI | Returns the HTML fragment for the Count Item Management Dashboard (CountItemManagementUI). |
| CountItemUpdateUI_UpdateItemAPI | PUT | `/ui/counts/{count_id}` | `submit` | `body` | `outerHTML` | CountItemUpdateUI | Updates the count item metadata (name, description). On success, returns the Dashboard UI fragment or an `HX-Redirect` to the dashboard. On validation failure, returns the `CountItemUpdateUI` fragment with error messages rendered in the `validation-error-area`. |
