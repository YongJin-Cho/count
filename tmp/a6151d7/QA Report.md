# QA Report

## 1. Q-Gate Result: PASS

## 2. 실제 배포·테스트 실행 결과 (Actual Deployment and Test Execution Results)

### 이미지 빌드
- **명령어**: `bash src/scripts/docker-build.sh latest`
- **결과**: 성공
- **주요 내용**: `count-management-service` 및 `count-processing-service` 이미지 빌드 완료.

### 배포 실행
- **명령어**: `bash src/scripts/deploy.sh`
- **결과**: 성공
- **주요 내용**: 
  - `count-system` 네임스페이스에 데이터베이스(PostgreSQL) 및 백엔드 서비스 배포 완료.
  - Ingress 및 Gateway API 리소스 생성 완료.
  - **추가 조치**: 빌드된 최신 이미지를 반영하기 위해 `kubectl rollout restart`를 통해 서비스를 강제 재시작함.

### 통합 테스트 실행
- **명령어**: `bash src/scripts/integration-test.sh` (MGMT_URL 및 PROC_URL 환경변수 사용)
- **환경**: 
  - `kubectl port-forward`를 통해 서비스 노출 (관리 서비스: 8888, 처리 서비스: 8889)
- **테스트 항목 및 결과**:

| 테스트 그룹 | 테스트 항목 | 결과 | 상세 내용 |
| :--- | :--- | :--- | :--- |
| **Connectivity** | 초기 연결 확인 | **Pass** | `MGMT_URL/api/v1/count-items` 접속 확인 |
| **Management API** | 아이템 등록 (POST) | **Pass** | 아이템 ID 반환 확인 |
| | 아이템 목록 조회 (GET) | **Pass** | 등록된 아이템 ID 포함 확인 |
| | 아이템 수정 (PUT) | **Pass** | 수정된 명칭 반영 확인 |
| | 아이템 삭제 (API) | **Pass** | 목록에서 삭제 확인 |
| **HTMX UI** | 아이템 등록 (HTMX) | **Pass** | `<tr>` HTML 조각 반환 확인 |
| | 아이템 삭제 (HTMX) | **Pass** | UI 삭제 동작 확인 |
| **External API** | Count 증가 (Increase) | **Pass** | `amount: 5` 요청 시 `value: 5` 확인 |
| | Count 감소 (Decrease) | **Pass** | `amount: 2` 요청 시 `value: 3` 확인 |
| | Count 초기화 (Reset) | **Pass** | `reset` 요청 시 `value: 0` 확인 |
| | 404 Error (Invalid ID) | **Pass** | 존재하지 않는 ID 요청 시 404 확인 |
| **Concurrency** | 원자성 테스트 (10회) | **Pass** | 10회 동시 증가 요청 후 최종 값 `11` 확인 (10 + 1) |

- **최종 판정**: **Pass**. 모든 Acceptance Criteria 및 시나리오 통과 확인.

## 3. 발견된 결함 (Discovered Defects)

| ID | Location | Problem | Severity | Status | Fix Instructions |
| :--- | :--- | :--- | :--- | :--- | :--- |
| #1 | src/count-processing-service | 초기 테스트 시 External API 경로 404 발생 | Medium | **Fixed** | `docker-build.sh` 후 기존 Pod가 `imagePullPolicy: Never` 설정으로 인해 새 이미지를 사용하지 않음. `kubectl rollout restart`를 통해 해결함. |
| #2 | src/scripts/integration-test.sh | 동시성 테스트 시 고정된 아이템 이름 사용 | Low | **Fixed** | 이전 테스트 잔재로 인해 중복 이름 에러 발생 가능성. `date +%s`를 이용한 유니크 네임으로 스크립트 보완. |

## 4. Proposals for user

1. **Health Check Endpoint 추가**: 현재 Kubernetes Liveness/Readiness Probe가 TCP 체크 또는 기본 404 응답에 의존하고 있습니다. `/health` 또는 `/ready` 전용 엔드포인트를 구현하여 서비스 상태를 더 정확하게 모니터링할 것을 권장합니다.
2. **Negative Count 제한**: 현재 사양상 Count 값이 음수가 될 수 있습니다. 비즈니스 로직상 0 이하로 내려가지 않아야 한다면, `Decrease` 처리 시 유효성 검사 추가를 고려해 보시기 바랍니다.
3. **Gateway API 안정성**: 현재 환경에서 Gateway API Controller(Kong)가 "Unknown" 상태로 조회됩니다. 실제 운영 환경에서는 Ingress Controller 또는 Gateway API Implementation이 정상 작동하는지 확인이 필요합니다.

**Q-Gate Status: PASS**
