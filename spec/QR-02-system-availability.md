# QR-02: 시스템 가용성

## 1. Purpose & Background
- **Purpose**: 외부 소스로부터 발생하는 데이터를 누락 없이 수집하기 위해 시스템의 지속적인 운영 가능 상태를 보장함.
- **Scope**: 중앙 API 수집 시스템 전체 서비스.

## 2. Measurement Indicators (Measurable Criteria)
- **Indicator**: 월간 가용성 (Monthly Uptime Percentage).
- **Target Value**: 99.9% 이상.
- **Unit/Conditions**:
  - 가용성 계산: `(총 시간 - 다운타임) / 총 시간 * 100`.
  - 다운타임 정의: 시스템이 외부 요청에 대해 정상적으로 응답하지 못하는 상태(HTTP 5xx 오류 또는 응답 없음)가 1분 이상 지속되는 경우.
  - 계획된 점검 시간은 다운타임에서 제외하되, 최소 7일 전 사전 공지되어야 함.

## 3. Verification Method
- **Measurement Tool/Method**: 
  - 외부 업타임 모니터링 도구(예: UptimeRobot, Pingdom) 또는 내부 상태 확인(Health Check) 메트릭.
  - Kubernetes의 Liveness/Readiness Probe 로그 분석.
- **Pass Criteria**: 월간 리포트 상의 가용성이 99.9%를 상회할 것.

## 4. References
- **Related FR**: 
  - [FR-01: 외부 count 수집 API](FR-01-collect-count-api.md)
- **Related SC**: 
  - [SC-02: 인프라 (Kubernetes)](SC-02-infra-k8s.md) - 고가용성 배포 및 자동 복구 환경.
