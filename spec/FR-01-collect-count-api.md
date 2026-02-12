# FR-01: 외부 count 수집 API

## 1. Feature Description
- **Purpose**: 외부 시스템으로부터 HTTP POST 방식을 통해 count 데이터를 수집하여 중앙 시스템에 기록함 (UR-01, UR-02 대응).
- **Scope**: 외부 시스템 식별자(external_id)와 count 값을 포함한 HTTP POST 요청 수신 기능.
- **References**: UR-01, UR-02, QR-03, SC-03 (Issue #71e94b3)

## 2. User Stories
- **FR-01-01**: 외부 시스템으로서, 나는 HTTP POST 방식을 통해 count 데이터를 전송하고 싶다. 그래야 중앙 시스템이 이를 기록할 수 있다.
- **FR-01-02**: 외부 시스템으로서, 나는 전송 시 나의 식별자와 count 값을 포함하고 싶다. 그래야 데이터가 올바르게 분류될 수 있다.

## 3. Acceptance Criteria (Gherkin)
### FR-01-01
- **FR-01-01-01**: HTTP POST 요청 수신 및 성공 응답 확인
    - **Given**: 외부 시스템이 유효한 Bearer Token 인증 정보를 가지고 있음
    - **When**: `/api/v1/collect` 경로로 유효한 데이터를 포함한 HTTP POST 요청을 보냄
    - **Then**: 서버는 HTTP 200 OK 또는 201 Created를 반환함

- **FR-01-01-02**: 유효하지 않은 토큰으로 요청 시 인증 실패 처리
    - **Given**: 외부 시스템이 유효하지 않거나 만료된 Bearer Token을 가지고 있음
    - **When**: `/api/v1/collect` 경로로 HTTP POST 요청을 보냄
    - **Then**: 서버는 HTTP 401 Unauthorized를 반환함

- **FR-01-01-03**: 권한이 없는 토큰으로 요청 시 접근 거부 처리
    - **Given**: 외부 시스템이 인증은 성공했으나 해당 API 호출 권한이 없는 Bearer Token을 가지고 있음
    - **When**: `/api/v1/collect` 경로로 HTTP POST 요청을 보냄
    - **Then**: 서버는 HTTP 403 Forbidden을 반환함

### FR-01-02
- **FR-01-02-01**: 요청 본문에 필수 데이터 포함 확인
    - **Given**: 외부 시스템이 `external_id`와 `count` 필드를 포함한 JSON 본문을 준비함
    - **When**: 유효한 인증 정보와 함께 해당 데이터를 POST 요청으로 전송함
    - **Then**: 서버는 데이터를 정상적으로 파싱하고 HTTP 200 OK 또는 201 Created를 반환함
