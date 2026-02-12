# CountCreateUI UI-API Identification

## UI-API List

| ID | Method | URL | Trigger | Target | Swap | Description |
|---|---|---|---|---|---|---|
| CreateCountAPI | POST | `/ui/counts` | click | `#count-list` | beforeend | 새로운 카운트 소스를 생성하고 목록(HTML fragment)을 추가합니다. |

## API Details

### CreateCountAPI
- **Description**: 입력된 정보를 바탕으로 새로운 카운트를 생성합니다. 성공 시 카운트 목록에 새로운 항목을 추가합니다.
- **Request**:
    - URL: `/ui/counts`
    - Method: `POST`
    - Body: 
        - `source_id` (string, required): 소스 식별자 (영문 소문자, 숫자, 하이픈)
        - `initial_value` (number, required): 초기 카운트 값
- **Response**:
    - Type: HTML Fragment
    - Content: 생성된 카운트의 상세 정보를 포함한 테이블 행(`<tr>`) 또는 리스트 아이템.
    - Success (200 OK): HTML Fragment 반환 및 클라이언트에서 폼 초기화 (`hx-on::after-request`)
    - Error (4xx/5xx): 에러 메시지 표시 (`hx-on::error`)
- **HTMX Specifics**:
    - **Indicator**: `#create-loading-indicator` (요청 중 표시)
    - **Before Request**: 기존 에러 메시지 숨김 (`document.querySelector('#count-create-error-msg').style.display = 'none';`)
    - **After Request**: 성공 시 폼 리셋 (`document.querySelector('#count-create-form').reset();`)
    - **Error**: 에러 메시지 표시 (`#count-create-error-msg` 영역)
