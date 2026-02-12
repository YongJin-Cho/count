# SC-02: 인프라 (Kubernetes 배포 환경)

## 1. Purpose & Background
- **Purpose**: 서비스의 확장성(Scalability), 가용성(Availability) 및 운영 자동화를 보장하기 위해 컨테이너 오케스트레이션 플랫폼을 표준화함.
- **Application Scope**: 서비스 배포, 운영 환경 및 인프라 구성 전체.

## 2. Constraint Content (Concrete Conditions)
- **Technology/Platform/Regulation**: Kubernetes (k8s)
- **Conditions**:
    - 배포 형태: 모든 서비스는 Docker 컨테이너 이미지로 패키징되어야 함.
    - 리소스 정의: Deployment, Service, HPA(Horizontal Pod Autoscaler), ConfigMap, Secret 등 Kubernetes 표준 리소스를 사용하여 배포 정의가 작성되어야 함.
    - 헬스 체크: Liveness Probe 및 Readiness Probe가 반드시 설정되어야 함.
    - 리소스 제한: 모든 Pod에는 CPU 및 Memory에 대한 Request와 Limit이 설정되어야 함.

## 3. Verification Method
- **Verification Method**: 
    - Kubernetes Manifest 파일 (YAML) 또는 Helm Chart 검사.
    - CI/CD 파이프라인에서 `kubectl apply --dry-run` 또는 `kubeval`, `pluto` 등을 사용한 검증.
    - 배포 후 Kubernetes 클러스터 내에서 Pod 상태 및 리소스 설정 확인.
- **On Violation**:
    - 표준 매니페스트 형식을 준수하지 않을 경우 배포를 승인하지 않음.
    - 리소스 제한이 설정되지 않은 서비스는 운영 클러스터에 배포 불가.

## 4. References
- **Related FR/QR**: 
    - [QR-02: 시스템 가용성](QR-02-system-availability.md)
