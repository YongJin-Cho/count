# CountManagementUI APIs

| ID | Method | Path | Trigger | Target | Swap | Description |
|---|---|---|---|---|---|---|
| GetCreateNewUIAPI | GET | /ui/counts/new | click | #main-content | innerHTML | 신규 카운트 생성 화면 요청 |
| GetCountListAPI | GET | /ui/counts | load, count-list-changed (body) | #count-list-body | innerHTML | 카운트 목록 데이터(<tr>) 요청 |
| PostIncrementCountAPI | POST | /ui/counts/${source_id}/increment | click | closest tr | outerHTML | 특정 카운트 1 증가 및 행 갱신 |
| PostDecrementCountAPI | POST | /ui/counts/${source_id}/decrement | click | closest tr | outerHTML | 특정 카운트 1 감소 및 행 갱신 |
| GetEditCountUIAPI | GET | /ui/counts/${source_id}/edit | click | #main-content | innerHTML | 특정 카운트 상세 수정 화면 요청 |
