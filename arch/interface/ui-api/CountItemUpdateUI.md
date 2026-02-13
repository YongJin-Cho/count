# UI-API List: CountItemUpdateUI

| ID | Method | Path | Trigger | Target | Swap | HTML Fragment Description |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| CountItemUpdateUI_GetDashboardAPI | GET | `/ui/counts` | `click` (link-back-dashboard, btn-cancel) | `body` | - | Returns the Dashboard UI HTML fragment. `hx-push-url="true"` is used for the back link. |
| CountItemUpdateUI_UpdateItemAPI | PUT | `/ui/counts/{count_id}` | `submit` | `body` | `outerHTML` | Updates the count item (name, description). Uses `#update-progress` indicator. On success, returns Dashboard fragment. On validation failure, returns `CountItemUpdateUI` fragment with errors in `validation-error-area`. |
