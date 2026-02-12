# QR-03: API 및 웹 UI 인증/보안

## 1. Purpose & Background
- **Purpose**: 신뢰할 수 있는 외부 시스템 및 권한 있는 사용자만이 API 및 웹 UI에 접근할 수 있도록 보장하여, 데이터의 무결성을 유지하고 비인가된 접근으로부터 시스템을 보호함. (UR-03, UR-09 관련)
- **Scope**: 
    - 외부 count 수집 API (FR-01) 및 통합 조회 API (FR-03)의 모든 접근 경로.
    - 카운트 관리 웹 UI (FR-04) 전반 및 UI에서 발생하는 모든 비동기 요청.

## 2. Measurement Indicators (Measurable Criteria)
- **Indicator**: 인증 및 암호화 수준, 토큰 관리 방식.
- **Target Value**: 
    - **인증 방식**: API 및 웹 UI 모두 `Authorization: Bearer <token>` 헤더를 통한 Bearer Token(JWT 또는 API Key) 방식 필수.
    - **웹 UI (HTMX) 인증 구현 상세**:
        - **토큰 저장**: 브라우저의 `Session Storage`에 보안 토큰을 저장.
        - **요청 헤더 주입**: HTMX의 `hx-headers` 속성 또는 `htmx:configRequest` 이벤트를 사용하여 모든 비동기 요청 헤더에 `Authorization` 토큰을 포함.
        - **초기 로드 및 세션 처리**: 초기 페이지 로드 시 유효한 세션/토큰이 없는 경우 로그인 페이지로 리다이렉트 처리함. (서버 사이드 세션과 연동하여 토큰을 주입하는 래퍼 구성 가능)
    - **전송 보안**: 모든 통신(API 및 웹 UI)은 TLS 1.2 버전 이상을 사용해야 함.
- **Unit/Conditions**:
    - 인증 실패(토큰 누락, 유효하지 않은 토큰): HTTP 401 Unauthorized 응답 반환율 100%.
    - 권한 부족(유효한 토큰이나 접근 권한 없음): HTTP 403 Forbidden 응답 반환율 100%.

## 3. Verification Method
- **Measurement Tool/Method**: 
    - **API 보안 테스트**: Postman, Newman 등을 사용하여 유효/무효/권한 없음 시나리오별 401, 403 응답 확인.
    - **웹 UI 인증 테스트**: 
        - 브라우저 개발자 도구를 사용하여 HTMX 요청 시 `Authorization` 헤더에 Bearer 토큰이 올바르게 포함되는지 검증.
        - `Session Storage`에 저장된 토큰 삭제 후 요청 시 401 에러 및 로그인 페이지 리다이렉트 발생 여부 확인.
    - **취약점 스캐닝**: OWASP ZAP 등을 활용하여 웹 UI 및 API의 인증 우회 가능성 점검.
    - **TLS 검증**: SSLLabs 또는 `openssl` 도구를 이용한 TLS 1.2+ 설정 확인.
- **Pass Criteria**: 
    - API 및 웹 UI의 모든 보호된 자원에 대한 요청이 인증 절차를 거칠 것.
    - 인증 정보 누락/오류 시 HTTP 401, 권한 부족 시 HTTP 403을 엄격히 반환할 것.
    - 웹 UI에서 비동기 요청(HTMX) 시 브라우저 스토리지의 토큰이 헤더에 누락 없이 주입될 것.

## 4. References
- **Related FR**: 
    - [FR-01: 외부 count 수집 API](FR-01-collect-count-api.md)
    - [FR-03: 통합 count 조회 API](FR-03-integrated-count-api.md)
    - [FR-04: 카운트 관리 웹 UI](FR-04-count-management-ui.md)
- **Related URD**:
    - [UR-03: API 보안 및 인증](../URD.md)
    - [UR-09: 카운트 관리 웹 UI](../URD.md)
- **Related SC**:
    - [SC-03: API 표준 준수](SC-03-api-standard.md)
    - [SC-05: UI 기술 스택 (HTMX)](SC-05-ui-tech-stack-htmx.md)
