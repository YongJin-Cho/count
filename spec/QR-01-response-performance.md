# QR-01: 응답 시간 성능

## 1. Purpose & Background
- **Purpose**: 외부 시스템으로부터 count 데이터를 수집하는 API의 지연 시간을 최소화하여, 데이터 수집의 효율성을 높이고 호출 시스템의 대기 시간을 줄이기 위함.
- **Scope**: 외부 count 수집 API (FR-01) 및 데이터 검증/응답 처리 (FR-02).

## 2. Measurement Indicators (Measurable Criteria)
- **Indicator**: API 응답 시간 (Latency) 95th Percentile (P95).
- **Target Value**: 200ms 이내.
- **Unit/Conditions**: 
  - 측정 환경: 운영(Production) 환경 또는 운영과 동일한 사양의 스테이징 환경.
  - 샘플 기간: 24시간 동안 발생한 모든 성공적인(HTTP 200 OK) 요청 대상.
  - 부하 조건: 정상 운영 범위(평균 TPS) 내의 요청량.

## 3. Verification Method
- **Measurement Tool/Method**: 
  - Application Monitoring Tool (예: Prometheus + Grafana, Datadog 등)을 통한 메트릭 수집.
  - 부하 테스트 도구(예: k6, JMeter)를 사용하여 정의된 TPS 내에서 P95 응답 시간 측정.
- **Pass Criteria**: 24시간 관찰 기간 또는 부하 테스트 결과, P95 응답 시간이 200ms 이하로 유지될 것.

## 4. References
- **Related FR**: 
  - [FR-01: 외부 count 수집 API](FR-01-collect-count-api.md)
  - [FR-02: 데이터 검증 및 응답 처리](FR-02-data-validation-response.md)
- **Related SC**: 
  - [SC-01: 기술 스택 (Go)](SC-01-tech-stack-go.md) - 고성능 처리를 위한 언어 선택.
