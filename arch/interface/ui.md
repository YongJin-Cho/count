# UI Identification

## UI List

| UI ID | Screen Name | Key Features | Related User Story |
|-------|-------------|--------------|-------------------|
| `CountManagementUI` | 카운트 목록 조회 | 전체 카운트 목록 및 현재 값 테이블 표시, 생성/수정 액션 버튼 제공 | FR-04-02 |
| `CountCreateUI` | 카운트 생성 폼 | 소스 식별자(Source ID) 입력 및 유효성 검사, 초기값 설정, 중복 ID 에러 메시지 표시 | FR-04-01, FR-04-04 |
| `CountEditUI` | 카운트 수정 및 수동 변경 | 특정 카운트 값의 수동 증감(+1, -1) 및 수동 값 수정/상태 업데이트 | FR-04-03, FR-04-04 |

## User Flows

### 1. 카운트 생성 흐름 (Count Creation Flow)
- **Path**: `CountManagementUI` (목록) -> `CountCreateUI` (폼 입력) -> 생성 완료 -> `CountManagementUI` (목록 자동 갱신)
- **Description**: 관리자가 새로운 카운트 소스를 등록하기 위해 생성 버튼을 클릭하고, 정보를 입력한 후 HTMX를 통해 페이지 새로고침 없이 목록에 즉시 반영되는 흐름.

### 2. 카운트 수동 변경 및 수정 흐름 (Count Manual Adjustment Flow)
- **Path**: `CountManagementUI` (목록) -> `CountEditUI` (증감 버튼 또는 수정 폼) -> 변경 완료 -> `CountManagementUI` (해당 항목 값 업데이트)
- **Description**: 관리자가 특정 카운트의 값을 '+' 또는 '-' 버튼을 눌러 즉시 변경하거나, 수정 폼을 통해 값을 직접 조정하는 흐름. HTMX를 사용하여 해당 행(row)만 부분 업데이트됨.

### 3. 카운트 목록 조회 흐름 (Count List View Flow)
- **Path**: `CountManagementUI`
- **Description**: 관리자가 카운트 관리 메인 페이지에 접속하여 시스템의 모든 카운트 상태를 한눈에 확인하는 기본 흐름.
