# Consolidated UI-Event List

| ID | Trigger Source | Action | Using UI | Description |
|---|---|---|---|---|
| CountListChangedEvent | `count-list-changed` (from: `body`) | `GetCountListAPI` 호출 | CountManagementUI | 카운트 목록 데이터 최신화 |
| HtmxResponseErrorEvent | `htmx:responseError` (from: `body`) | error-toast 표시 | CountManagementUI | 서버 통신 에러 발생 시 사용자 알림 |
