# QA Report - Issue #1aeddeb

## Q-Gate Result: **PASS**

통합 count 조회 기능 및 페이지네이션/필터링 기능에 대한 통합 테스트를 수행한 결과, 모든 테스트 케이스를 통과하였으며 시스템이 설계대로 동작함을 확인하였습니다.

## 실제 배포·테스트 실행 결과 (Actual Deployment and Test Execution Results)

### 1. 이미지 빌드 (Docker Build)
- **명령어**: `bash src/scripts/docker-build.sh`
- **결과**: **성공**
- **내용**: `count-api-service:latest` 이미지 빌드 완료 및 로컬 도커 레지스트리 등록 확인.

### 2. Kubernetes 배포 (K8S Deployment)
- **명령어**: `bash src/scripts/deploy.sh`
- **결과**: **성공**
- **내용**: 
    - Namespace `count-collection-system` 생성/확인
    - StatefulSet `count-api-service` 배포 및 롤아웃 완료 (Ready: 1/1)
    - Service 및 Gateway API(HTTPRoute) 설정 완료.
- **특이사항**: 기존에 남아있던 구버전 Deployment(`count-api-service`)가 새 StatefulSet과 동일한 라벨을 사용하여 트래픽 간섭(404 발생)을 일으키는 것이 발견되어 수동으로 제거하였습니다.

### 3. 통합 테스트 실행 (Integration Test)
- **명령어**: `bash src/scripts/integration-test.sh`
- **결과**: **모든 테스트 통과 (PASS)**
- **주요 테스트 항목 및 결과**:
    - **Test 1: Health Check**: `GET /health` -> 200 OK (성공)
    - **Test 2: Valid Count Collection**: `POST /api/v1/collect` -> 성공
    - **Test 5: Integrated Query (신규)**: `GET /api/v1/counts` -> 전체 데이터 및 기본 조회 성공
    - **Test 6: Pagination & Filtering (신규)**: `GET /api/v1/counts?external_id=...&limit=2&offset=1` -> 필터링 및 페이지네이션 정확도 확인 (성공)
    - **Test 7: Performance**: 평균 응답 시간 8.97ms (기준 100ms 미만 충족)

---

## 발견된 결함 및 조치 사항 (Discovered Defects & Actions)

### 결함 #1: 테스트용 토큰 권한 부족
- **위치**: `src/scripts/gen-token.go`
- **현상**: 조회 API(`GET /api/v1/counts`) 호출 시 `403 Forbidden` 반환.
- **원인**: 테스트 스크립트에서 생성하는 JWT 토큰에 `query` 권한이 누락되어 있었음.
- **조치**: `src/scripts/gen-token.go`를 수정하여 `collect`와 `query` 권한을 모두 포함하도록 개선.
- **상태**: **조치 완료 (Fixed)**

### 결함 #2: 구버전 리소스(Deployment) 간섭
- **위치**: Kubernetes Cluster (`count-collection-system` namespace)
- **현상**: 간헐적으로 조회 API 호출 시 `404 page not found` 반환.
- **원인**: 이전 버전에 사용되던 `Deployment` 리소스가 삭제되지 않고 남아있었으며, 새 `StatefulSet`과 동일한 `app: count-api-service` 라벨을 사용하여 서비스 트래픽이 구버전 Pod로 유입됨.
- **조치**: `kubectl delete deployment count-api-service -n count-collection-system` 명령으로 구버전 리소스 제거.
- **권장 사항**: 향후 배포 스크립트(`deploy.sh`)에 리소스 정리 로직을 포함하거나, Helm 등을 사용하여 리소스 생명주기를 관리할 것을 권장함.
- **상태**: **조치 완료 (Fixed)**

---

## 사용자 제안 사항 (Proposals for user)

1. **배포 스크립트 고도화**: 현재 `deploy.sh`는 신규 리소스를 `apply`만 수행합니다. 기존 리소스와의 충돌을 방지하기 위해 사용하지 않는 리소스를 정리하는 로직 추가가 필요합니다.
2. **테스트 데이터 초기화**: 통합 테스트 시 기존 데이터의 영향을 받지 않도록 테스트 시작 전 데이터를 초기화하거나, 매번 고유한 ID를 사용하는 로직을 강화할 필요가 있습니다 (이번 테스트에서는 고유 ID 사용으로 보완함).
3. **API 에러 메시지 세분화**: 현재 403 Forbidden 발생 시 구체적인 부족 권한이 명시되지 않습니다. 보안 정책에 따라 다르겠지만, 개발 편의성을 위해 에러 메시지에 필요한 스코프를 명시하는 것을 고려해볼 수 있습니다.
