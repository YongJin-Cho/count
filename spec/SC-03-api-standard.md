# SC-03: API 표준 준수

## 1. Purpose & Background
- **Purpose**: 외부 시스템과의 원활한 연동 및 인터페이스의 일관성을 유지하기 위해 표준화된 API 통신 방식을 정의함.
- **Application Scope**: 외부 노출되는 모든 수집 API 인터페이스.

## 2. Constraint Content (Concrete Conditions)
- **Technology/Platform/Regulation**: REST API 규격 준수
- **Conditions**:
    - 프로토콜: HTTP/1.1 이상 (TLS 1.2 이상 권장).
    - HTTP 메서드: 데이터 수집 API의 경우 반드시 **HTTP POST** 메서드를 사용해야 함.
    - 데이터 포맷: 요청 및 응답 본문(Body)은 **JSON** 형식을 사용해야 함.
    - Content-Type: `application/json`을 명시해야 함.
    - 상태 코드:
        - **성공**: `200 OK` (요청 처리 성공), `201 Created` (데이터 생성 완료).
        - **클라이언트 오류**: `400 Bad Request` (잘못된 데이터 형식/유효성 검증 실패), `401 Unauthorized` (인증 누락/실패), `403 Forbidden` (접근 권한 없음), `405 Method Not Allowed` (정의되지 않은 HTTP 메서드 사용).
        - **서버 오류**: `500 Internal Server Error` (서버 내부 처리 중 예외 발생).

## 3. Verification Method
- **Verification Method**: 
    - API 명세서(Swagger/OpenAPI) 검토.
    - Postman 또는 curl을 이용한 실제 호출 테스트를 통해 메서드 및 포맷 준수 여부 확인.
    - 유닛 테스트 및 통합 테스트 시 HTTP 메서드 검증 로직 포함.
- **On Violation**:
    - GET, PUT 등 허용되지 않은 메서드로 요청 시 `405 Method Not Allowed`를 반환해야 함.
    - JSON 형식이 아닌 요청에 대해 `400 Bad Request` 처리.

## 4. References
- **Related FR/QR**: 
    - [FR-01: 외부 count 수집 API](FR-01-collect-count-api.md)
    - [FR-02: 데이터 검증 및 응답 처리](FR-02-data-validation-response.md)
