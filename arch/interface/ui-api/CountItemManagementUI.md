# CountItemManagementUI UI-APIs

| ID | Method | Path | Trigger | Target | Swap | Description |
|---|---|---|---|---|---|---|
| RegisterCountItemAPI | POST | `/ui/count-items` | `click` (on `btn-register`) | `#count-item-list` | `beforeend` | Registers a new count item and appends the new item fragment to the list. Resets the form on success. |
| ListCountItemsAPI | GET | `/ui/count-items` | `load` (on `#count-item-list`) | `this` | `innerHTML` | Fetches and displays the initial list of registered count items. |
| GetCountItemValueAPI | GET | `/ui/counts/{id}/value` | `every 30s`, `click` (on `item-value-{id}`) | `this` | `innerHTML` | Fetches the current value of a specific count item and updates its display. |
| DeleteCountItemAPI | DELETE | `/ui/count-items/{id}` | `click` (on `btn-delete-{id}`) | `closest .count-item-row` | `outerHTML` | Deletes a count item and removes its row fragment from the list. Requires user confirmation. |
