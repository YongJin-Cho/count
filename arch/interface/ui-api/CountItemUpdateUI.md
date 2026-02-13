# UI-API List: CountItemUpdateUI

| ID | Method | Path | Trigger | Target | Swap | HTML Fragment Description |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| CountItemUpdateUI_GetDashboardAPI | GET | `/ui/counts` | `click` (link-back-dashboard, btn-cancel) | `body` | - | Returns the HTML fragment for the Count Item Management Dashboard (CountItemManagementUI). |
| CountItemUpdateUI_UpdateItemAPI | PUT | `/ui/counts/{count_id}` | `submit` | `body` | `outerHTML` | Updates the count item metadata (name, description). On success, returns the Dashboard UI fragment or an `HX-Redirect` to the dashboard. On validation failure, returns the `CountItemUpdateUI` fragment with error messages rendered in the `validation-error-area`. |
