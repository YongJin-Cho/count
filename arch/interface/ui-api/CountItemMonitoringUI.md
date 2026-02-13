# CountItemMonitoringUI UI-APIs

| ID | Method | Path | Trigger | Target | Swap | Description |
|---|---|---|---|---|---|---|
| GetCountValueAPI | GET | `/ui/counts/{id}/value` | `load, every 5s` | `#value-display` | `innerHTML` | Periodically fetches the current value of a count item and updates the display area. |
