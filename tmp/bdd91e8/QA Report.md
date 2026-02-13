# QA Report - Issue #bdd91e8 (K8S Environment Verification - Retry)

## Q-Gate Result: PASS

## 실제 배포·테스트 실행 결과 (Actual Deployment and Test Execution Results)

### 1. 이미지 빌드 (Image Build)
- **명령어**: `bash src/scripts/docker-build.sh latest`
- **결과**: **성공**
- **내용**: `count-management-service:latest` 및 `count-processing-service:latest` 이미지 재빌드 완료.

### 2. 배포 실행 (Deployment Execution)
- **명령어**: `bash src/scripts/deploy.sh`
- **결과**: **성공**
- **내용**: 
    - **Gateway API CRD**: `deploy.sh`가 클러스터 내 CRD 존재 여부를 확인하고, 필요한 경우 설치하도록 업데이트됨을 확인.
    - **리소스 생성**: `GatewayClass`, `Gateway`, `HTTPRoute` 리소스가 오류 없이 정상적으로 생성됨. (이전의 "no matches for kind" 오류 해결됨)
    - **컴포넌트 상태**: `count-system` 네임스페이스 내 모든 Pod (`management-db`, `processing-db`, `count-management-service`, `count-processing-service`)가 `Running` 상태임을 확인.

### 3. 통합 테스트 실행 (Integration Test Execution)
- **환경**: Kubernetes (via Port-Forwarding)
- **명령어**: 
    - `kubectl port-forward -n count-system svc/count-management-service 8888:8080 &`
    - `bash src/scripts/integration-test.sh http://localhost:8888`
- **결과**: **PASS**
- **상세 결과**:
    - **Testing Register Item API**: 성공 (ID: 298ef379-e0e3-448a-a9d6-daf24950fdc9)
    - **Testing List Items API**: 성공 (아이템 발견됨)
    - **Testing Update Item API**: 성공 (이름 업데이트 확인)
    - **Testing UI Register Item (HTMX)**: 성공 (HTML fragment 반환 확인)
    - **Testing UI Delete Item (HTMX)**: 성공 (목록에서 제거 확인)
    - **Testing API Delete Item**: 성공 (목록에서 제거 확인)
- **검증 내용**: MSA 간 통신(Management -> Processing) 및 DB 연동이 K8S 내부 네트워크 환경에서 정상 작동함을 확인.

---

## 발견된 결함 (Discovered Defects)

### [Defect #3] Gateway/Ingress Controller 부재 (Infrastructure Issue)
- **상황**: Gateway API 리소스는 정상적으로 생성되었으나, 이를 처리할 Controller(예: Kong)가 클러스터에 존재하지 않아 Gateway가 `Programmed: Unknown` 상태로 유지됨.
- **영향**: `http://localhost`를 통한 외부 접속이 불가능함. (테스트는 Port-Forwarding으로 우회 수행)
- **수정 제안**: 
    - `deploy.sh`에서 Gateway가 Ready 상태가 될 때까지 기다리는 로직 추가.
    - 클러스터 환경에 Gateway Controller 또는 Ingress Controller가 설치되어 있는지 확인 필요.

---

## Proposals for user

1. **데이터 영속성 확보**: 현재 `management-db` 및 `processing-db`가 `Deployment`로 구성되어 있으며 Volume Mount가 없습니다. 운영 환경을 위해 `StatefulSet`과 `PersistentVolumeClaim` 사용을 권장합니다.
2. **Health Check 엔드포인트 구현**: `count-management-service`에 `/health` 엔드포인트를 추가하여 K8S의 `livenessProbe` 및 `readinessProbe`가 애플리케이션 상태를 더 정확히 판단할 수 있도록 개선을 권장합니다.
3. **Gateway Ready 대기 로직**: `deploy.sh`에 `kubectl wait --for=condition=programmed gateway/count-gateway -n count-system`와 같은 대기 로직을 추가하여 배포 완결성을 높일 것을 권장합니다.
