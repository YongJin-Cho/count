# CountItemHistoryUI UI-API List

| ID | Endpoint | Method | Trigger | Target | Swap | Response Description |
|:---|:---|:---|:---|:---|:---|:---|
| GetCountItemHistoryAPI | `/ui/counts/{id}/history` | GET | `load` | `#history-table-container` | `innerHTML` | Returns an HTML fragment containing either a table (`history-table`) with chronological change logs (timestamp, operation, amount, source) or an empty state message (`empty-history`) if no records exist. |
