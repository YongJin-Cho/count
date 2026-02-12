# SC-04: 조회 API 메서드 준수

## 1. Purpose & Background
- **Purpose**: 수집된 통합 count 데이터를 조회하는 API의 일관성을 유지하고, 표준 HTTP RESTful 규약을 준수하여 예측 가능하고 효율적인 캐싱 환경을 제공하기 위함이다.
- **Application Scope**: 통합 count 조회 기능을 제공하는 모든 API (FR-03 관련 API 전체)

## 2. Constraint Content (Concrete Conditions)
- **HTTP Method**: 통합 count 조회 API는 반드시 **HTTP GET** 메서드를 사용해야 한다.
- **URI Naming**: 조회 대상 리소스를 명확히 식별할 수 있는 명사 형태의 경로를 사용한다. (예: `/api/v1/counts`)
- **Query Parameters**: 필터링(소스별) 및 정렬 등의 조건은 Request Body가 아닌 Query String을 통해 전달해야 한다.
- **Idempotency**: 조회 API는 여러 번 호출해도 서버의 상태를 변경하지 않는 멱등성(Idempotency)과 안전성(Safe)을 보장해야 한다.

## 3. Verification Method
- **Verification Method**: 
    - API 설계 문서 및 Swagger/OAS 명세 검토를 통해 GET 메서드 사용 여부 확인
    - 통합 테스트 단계에서 다른 메서드(POST, PUT 등) 요청 시 실패하는지 검증
- **On Violation**: 
    - 설계 가이드라인 미준수 시 설계 승인 반려 및 수정 요청
    - 구현 시 다른 메서드를 사용한 경우 코드 리뷰에서 차단 및 GET 메서드로 변경

## 4. References
- **Related UR**: [UR-08] 수집된 소스별/전체 통합 count 조회 기능 제공
- **Related FR**: FR-03 통합 count 조회 API
