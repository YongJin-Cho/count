# SC-01: 기술 스택 (Go 언어 사용)

## 1. Purpose & Background
- **Purpose**: 백엔드 API 시스템의 고성능 처리, 동시성 관리 용이성 및 운영 안정성을 확보하기 위해 특정 언어를 표준으로 지정함.
- **Application Scope**: 백엔드 API 서비스 전체 및 관련 라이브러리 개발.

## 2. Constraint Content (Concrete Conditions)
- **Technology/Platform/Regulation**: Go (Golang)
- **Conditions**:
    - Go 버전: 1.21 버전 이상의 Stable 버전을 사용해야 함.
    - 패키지 관리: `go mod`를 사용하여 의존성을 관리해야 함.
    - 코드 스타일: `gofmt` 또는 `goimports`를 사용하여 표준 코드 스타일을 유지해야 함.

## 3. Verification Method
- **Verification Method**: 
    - 빌드 단계에서 `go version` 확인.
    - `go.mod` 파일 존재 및 의존성 리스트 확인.
    - CI 단계에서 `go fmt` 또는 `golangci-lint`를 통한 코드 스타일 검사.
- **On Violation**:
    - 규격에 맞지 않는 언어로 개발된 모듈은 빌드/배포 대상에서 제외함.
    - 코드 스타일 미준수 시 CI/CD 파이프라인에서 실패 처리.

## 4. References
- **Related FR/QR**: 
    - [FR-01: 외부 count 수집 API](FR-01-collect-count-api.md)
    - [QR-01: 응답 시간 성능](QR-01-response-performance.md)
