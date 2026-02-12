# Consolidated UI-API List

| ID | Method | URL | Trigger | Target | Swap | Using UI | Description |
|---|---|---|---|---|---|---|---|
| GetCountListAPI | GET | `/ui/counts` | load, count-list-changed (body) | `#count-list-body` | innerHTML | CountManagementUI | 카운트 목록 데이터(`<tr>`) 요청 |
| CreateCountAPI | POST | `/ui/counts` | click | `#count-list` | beforeend | CountCreateUI | 새로운 카운트 소스 생성 및 목록 추가 |
| GetCreateNewUIAPI | GET | `/ui/counts/new` | click | `#main-content` | innerHTML | CountManagementUI | 신규 카운트 생성 화면 요청 |
| CancelEditCountAPI | GET | `/ui/counts/{source_id}` | click | `#count-row-{source_id}` | outerHTML | CountEditUI | 수정 취소 및 조회 모드 행 반환 |
| UpdateCountAPI | PUT | `/ui/counts/{source_id}` | submit | `#count-row-{source_id}` | outerHTML | CountEditUI | 카운트 값 수정 및 행 갱신 |
| GetEditCountUIAPI | GET | `/ui/counts/{source_id}/edit` | click | `#main-content` | innerHTML | CountManagementUI | 특정 카운트 상세 수정 화면 요청 |
| IncrementCountAPI | POST | `/ui/counts/{source_id}/increment` | click | `#count-row-{source_id}` (or `closest tr`) | outerHTML | CountEditUI, CountManagementUI | 카운트 1 증가 및 행 갱신 |
| DecrementCountAPI | POST | `/ui/counts/{source_id}/decrement` | click | `#count-row-{source_id}` (or `closest tr`) | outerHTML | CountEditUI, CountManagementUI | 카운트 1 감소 및 행 갱신 |
