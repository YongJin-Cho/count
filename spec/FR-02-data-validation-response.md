# FR-02: 데이터 검증 및 성공/실패 응답 처리

## 1. Feature Description
- **Purpose**: 수집된 데이터의 정밀한 유효성 검사를 수행하고, 결과에 따라 적절한 HTTP 상태 코드와 메시지를 반환함 (UR-01, UR-02 대응).
- **Scope**: 필수 필드 누락 검사, 데이터 형식 검사, 성공 및 실패 응답 메시지 정의.
- **References**: UR-01, UR-02, SC-03 (Issue #71e94b3)

## 2. User Stories
- **FR-02-01**: 시스템으로서, 나는 전송된 `external_id`와 `count` 값이 유효한지 검증하고 싶다. 그래야 잘못된 데이터가 저장되는 것을 방지할 수 있다.
- **FR-02-02**: 외부 시스템으로서, 나는 데이터가 성공적으로 처리되었을 때 성공 응답을 받고 싶다.
- **FR-02-03**: 외부 시스템으로서, 나는 데이터가 유효하지 않을 때 실패 이유와 함께 실패 응답을 받고 싶다.

## 3. Acceptance Criteria (Gherkin)
### FR-02-01
- **FR-02-01-01**: `external_id` 누락 시 처리
    - **Given**: 외부 시스템이 `external_id` 필드를 제외하고 POST 요청을 보냄
    - **When**: 서버가 요청을 검증함
    - **Then**: 서버는 HTTP 400 Bad Request를 반환하고 JSON 응답 본문에 `"message": "missing external_id"`를 포함함

- **FR-02-01-02**: `count` 필드 누락 시 처리
    - **Given**: 외부 시스템이 `count` 필드를 제외하고 POST 요청을 보냄
    - **When**: 서버가 요청을 검증함
    - **Then**: 서버는 HTTP 400 Bad Request를 반환하고 JSON 응답 본문에 `"message": "missing count"`를 포함함

- **FR-02-01-03**: `count` 값이 유효하지 않을 때 처리 (음수)
    - **Given**: 외부 시스템이 `count` 값을 음수로 설정하여 POST 요청을 보냄
    - **When**: 서버가 요청을 검증함
    - **Then**: 서버는 HTTP 400 Bad Request를 반환하고 JSON 응답 본문에 `"message": "invalid count value"`를 포함함

### FR-02-02
- **FR-02-02-01**: 유효한 데이터 전송 시 성공 응답
    - **Given**: 외부 시스템이 유효한 `external_id` (string)와 0 이상의 `count` (integer)를 전송함
    - **When**: 서버가 요청을 성공적으로 처리함
    - **Then**: 서버는 HTTP 200 OK 또는 201 Created를 반환하고 JSON 응답 본문에 `"status": "success"` 메시지를 포함함

### FR-02-03
- **FR-02-03-01**: 잘못된 JSON 형식 요청 시 처리
    - **Given**: 외부 시스템이 JSON 문법에 맞지 않는 데이터를 전송함
    - **When**: 서버가 요청을 파싱하려고 시도함
    - **Then**: 서버는 HTTP 400 Bad Request를 반환하고 JSON 응답 본문에 `"message": "invalid JSON format"`을 포함함
