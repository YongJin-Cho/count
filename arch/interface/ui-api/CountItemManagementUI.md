# CountItemManagementUI UI-APIs

| ID | Method | Path | Trigger | Target | Swap | Description |
|---|---|---|---|---|---|---|
| RegisterCountItemAPI | POST | `/ui/count-items` | `click` (on `btn-register`) | `#count-item-list` | `beforeend` | Registers a new count item and appends the new item fragment to the list. |
| ListCountItemsAPI | GET | `/ui/count-items` | `load` (on `#count-item-list`) | `this` | `innerHTML` | Fetches and displays the initial list of registered count items. |
| DeleteCountItemAPI | DELETE | `/ui/count-items/{id}` | `click` (on `btn-delete-{id}`) | `closest .count-item-row` | `outerHTML` | Deletes a count item and removes its row fragment from the list. |
