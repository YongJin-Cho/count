# FR-03: 통합 count 조회 API

## 1. Feature Description
- **Purpose**: 외부 소스별로 수집된 count 데이터와 전체 통합된 count 합계를 조회할 수 있는 인터페이스를 제공함 (UR-08 대응).
- **Scope**: HTTP GET 메서드를 통한 데이터 조회 기능, 특정 소스 식별자(external_id)를 이용한 필터링 기능.
- **References**: [UR-08](../URD.md), [QR-03](QR-03-api-security-auth.md), [QR-04](QR-04-query-performance.md), SC-04

## 2. User Stories
- **FR-03-01**: 사용자로서, 나는 모든 외부 소스의 count 합계와 소스별 상세 count 목록을 한 번에 확인하고 싶다. 그래야 전체 현황을 파악할 수 있다.
- **FR-03-02**: 사용자로서, 나는 특정 외부 소스의 식별자를 지정하여 해당 소스의 count 값만 조회하고 싶다. 그래야 특정 소스에 대한 집중적인 분석이 가능하다.

## 3. Acceptance Criteria (Gherkin)
### FR-03-01
- **FR-03-01-01**: 전체 통합 count 및 소스별 목록 조회 성공
    - **Given**: 시스템에 여러 외부 소스(예: 'source-A', 'source-B')의 데이터가 수집되어 있음
    - **When**: `/api/v1/counts?limit=10&offset=0` 경로로 HTTP GET 요청을 보냄 (유효한 인증 정보 포함)
    - **Then**: 서버는 HTTP 200 OK와 함께 전체 합계(total_count) 및 각 소스별 ID와 count가 포함된 목록(sources)을 JSON 형태로 반환함 (페이지네이션 적용)

- **FR-03-01-02**: 인증 실패 시 조회 거부
    - **Given**: 유효하지 않거나 누락된 Bearer Token을 사용함
    - **When**: `/api/v1/counts` 경로로 HTTP GET 요청을 보냄
    - **Then**: 서버는 HTTP 401 Unauthorized를 반환함

### FR-03-02
- **FR-03-02-01**: 특정 소스 식별자를 통한 필터링 조회 성공
    - **Given**: 시스템에 'source-A'와 'source-B'의 데이터가 수집되어 있음
    - **When**: `/api/v1/counts?external_id=source-A&limit=10&offset=0` 경로로 HTTP GET 요청을 보냄
    - **Then**: 서버는 HTTP 200 OK와 함께 'source-A'의 데이터만 포함된 결과(total_count가 source-A의 값과 동일)를 반환함

- **FR-03-02-02**: 존재하지 않는 소스 식별자로 조회 시 결과 처리
    - **Given**: 시스템에 'source-Z'에 대한 데이터가 존재하지 않음
    - **When**: `/api/v1/counts?external_id=source-Z` 경로로 HTTP GET 요청을 보냄
    - **Then**: 서버는 HTTP 200 OK를 반환하며, `total_count`는 0이고 `sources` 목록은 비어있어야 함
