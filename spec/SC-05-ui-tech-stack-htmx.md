# SC-05: UI 기술 스택 (HTMX)

## 1. Purpose & Background
- **Purpose**: 복잡한 클라이언트 사이드 상태 관리 및 무거운 프론트엔드 프레임워크(React, Vue 등) 도입 없이, 서버 사이드 렌더링(SSR) 기반에서 동적이고 인터랙티브한 사용자 경험(UX)을 효율적으로 제공하기 위함.
- **Application Scope**: [FR-04] 카운트 관리 웹 UI 기능을 포함한 시스템의 모든 웹 프론트엔드 영역.

## 2. Constraint Content (Concrete Conditions)
- **HTMX 활용**: 동적 데이터 갱신, 폼 제출, 부분 화면 업데이트 등 모든 인터랙티브한 기능은 [HTMX](https://htmx.org/) 라이브러리를 사용하여 구현해야 함.
- **SSR 및 HTML Fragment**: 
    - 서버(Go)는 전체 페이지 요청뿐만 아니라, HTMX 요청에 대해 필요한 최소 단위의 HTML 조각(Fragment)을 반환해야 함.
    - Go의 `html/template` 패키지 또는 호환 가능한 템플릿 엔진을 사용하여 서버 사이드에서 HTML을 생성함.
- **프레임워크 배제**: React, Vue, Angular, Svelte 등 SPA(Single Page Application) 지향의 중량급 프레임워크 사용을 금지함.
- **최소한의 JavaScript**: 
    - 비즈니스 로직은 최대한 서버에서 처리하며, 클라이언트 사이드 스크립트는 HTMX로 처리하기 어려운 복잡한 UI 인터랙션이나 라이브러리 연동에 한해 최소한으로 작성함.
    - 가급적 바닐라 자바스크립트(Vanilla JS) 또는 HTMX와 궁합이 좋은 경량 스크립트(예: Alpine.js) 사용을 권장함.

## 3. Verification Method
- **Verification Method**:
    - **코드 리뷰**: HTML 템플릿 내에 `hx-get`, `hx-post`, `hx-target`, `hx-swap` 등 HTMX 특성이 적절히 사용되었는지 확인.
    - **네트워크 검사**: 브라우저 개발자 도구의 Network 탭에서 페이지 전환 없이 부분적인 HTML 조각이 송수신되는지 확인.
    - **번들 크기 확인**: 대규모 프론트엔드 프레임워크 라이브러리가 포함되어 있지 않은지, 초기 로딩 시 JS 번들 크기가 비정상적으로 크지 않은지 확인.
- **On Violation**:
    - 프론트엔드 프레임워크(React 등)가 사용된 경우, 해당 기능을 HTMX와 SSR 기반으로 재구현해야 함.
    - 불필요하게 큰 클라이언트 사이드 스크립트가 발견될 경우 서버 사이드 로직으로 이전하거나 경량화함.

## 4. References
- **Related UR**: [UR-10] 웹 UI HTMX 활용
- **Related FR**: [FR-04] 카운트 관리 웹 UI
- **Related QR**: [QR-05] UI 비동기 업데이트 성능
