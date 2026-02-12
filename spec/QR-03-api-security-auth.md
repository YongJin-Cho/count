# QR-03: API 인증 및 보안

## 1. Purpose & Background
- **Purpose**: 신뢰할 수 있는 외부 시스템만이 API에 접근할 수 있도록 보장하여, 데이터의 무결성을 유지하고 비인가된 접근으로부터 시스템을 보호함. (UR-03 관련)
- **Scope**: 외부 count 수집 API (FR-01)의 모든 접근 경로.

## 2. Measurement Indicators (Measurable Criteria)
- **Indicator**: 인증 및 암호화 수준.
- **Target Value**: 
  - 인증 방식: Authorization 헤더를 통한 Bearer Token(JWT 또는 API Key) 방식 필수.
  - 전송 보안: 모든 통신은 TLS 1.2 버전 이상을 사용해야 함.
- **Unit/Conditions**:
  - 인증 실패(토큰 누락, 유효하지 않은 토큰): HTTP 401 Unauthorized 응답 반환율 100%.
  - 권한 부족(유효한 토큰이나 접근 권한 없음): HTTP 403 Forbidden 응답 반환율 100%.

## 3. Verification Method
- **Measurement Tool/Method**: 
  - 보안 취약점 스캐닝 도구(예: OWASP ZAP)를 사용한 인증 우회 테스트.
  - API 자동화 테스트 도구(Postman, Newman 등)를 사용하여 유효한/유효하지 않은/권한 없는 토큰 테스트 시나리오 수행.
  - SSLLabs 등을 이용한 TLS 설정 검증.
- **Pass Criteria**: 
  - 모든 API 요청이 사전에 정의된 인증 절차를 통과해야만 처리될 것.
  - 인증 정보 누락/오류 시 HTTP 401 Unauthorized를 반환할 것.
  - 권한 부족 시 HTTP 403 Forbidden을 반환할 것.

## 4. References
- **Related FR**: 
  - [FR-01: 외부 count 수집 API](FR-01-collect-count-api.md)
- **Related URD**:
  - [UR-03: 외부 API 요청은 신뢰할 수 있는 소스로부터의 접근임을 보장하기 위해 보안 및 인증 절차를 거쳐야 한다.](../URD.md)
- **Related SC**:
  - [SC-03: API 표준 준수](SC-03-api-standard.md)
