# CountEditUI UI-API Identification

## UI-API List

| ID | Method | URL | Trigger | Target | Swap | Description |
|---|---|---|---|---|---|---|
| DecrementCountAPI | POST | `/ui/counts/{source_id}/decrement` | click | `#count-row-{source_id}` | outerHTML | 카운트 값을 1 감소시키고 결과 행(HTML fragment)을 반환 |
| IncrementCountAPI | POST | `/ui/counts/{source_id}/increment` | click | `#count-row-{source_id}` | outerHTML | 카운트 값을 1 증가시키고 결과 행(HTML fragment)을 반환 |
| UpdateCountAPI | PUT | `/ui/counts/{source_id}` | submit | `#count-row-{source_id}` | outerHTML | 카운트 값을 직접 입력한 값으로 수정하고 결과 행(HTML fragment)을 반환 |
| CancelEditCountAPI | GET | `/ui/counts/{source_id}` | click | `#count-row-{source_id}` | outerHTML | 수정을 취소하고 해당 카운트의 조회 모드 행(HTML fragment)을 반환 |

## API Details

### DecrementCountAPI
- **Description**: 현재 카운트 값을 1 줄입니다.
- **Request**:
    - URL: `/ui/counts/{source_id}/decrement`
    - Method: `POST`
- **Response**:
    - Type: HTML Fragment
    - Content: 업데이트된 카운트 정보를 담은 `CountEditUI`의 테이블 행(`<tr>`). 
    - Note: 값이 성공적으로 변경되면 수정 모드 행을 다시 반환하여 현재 값을 갱신하거나, 특정 조건에 따라 `CountViewUI`로 전환될 수 있습니다. (이 UI 스펙에서는 `outerHTML`로 자기 자신을 대체함)
    - Success: `200 OK`. 

### IncrementCountAPI
- **Description**: 현재 카운트 값을 1 늘립니다.
- **Request**:
    - URL: `/ui/counts/{source_id}/increment`
    - Method: `POST`
- **Response**:
    - Type: HTML Fragment
    - Content: 업데이트된 카운트 정보를 담은 `CountEditUI`의 테이블 행(`<tr>`).
    - Note: `outerHTML` 스왑을 통해 행 전체를 갱신합니다.
    - Success: `200 OK`.

### UpdateCountAPI
- **Description**: 사용자가 입력한 값으로 카운트를 설정합니다.
- **Request**:
    - URL: `/ui/counts/{source_id}`
    - Method: `PUT`
    - Body: `value` (number, required)
- **Response**:
    - Type: HTML Fragment
    - Content: 업데이트된 카운트 정보를 담은 `CountViewUI`의 테이블 행(`<tr>`).
    - Note: 업데이트 성공 후에는 보통 조회 모드(`CountViewUI`)로 복귀하는 것이 일반적이나, 스펙상 `#count-row-{source_id}`를 `outerHTML`로 대체하므로 서버는 적절한 상태의 행 fragment를 반환해야 합니다.
    - Success: `200 OK`.

### CancelEditCountAPI
- **Description**: 편집을 취소하고 조회 모드로 돌아갑니다.
- **Request**:
    - URL: `/ui/counts/{source_id}`
    - Method: `GET`
- **Response**:
    - Type: HTML Fragment
    - Content: 해당 카운트의 정보를 담은 `CountViewUI`의 테이블 행(`<tr>`).
    - Success: `200 OK`.
