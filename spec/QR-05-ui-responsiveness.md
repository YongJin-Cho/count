# QR-05: UI 비동기 업데이트 성능

## 1. Purpose & Background
- **Purpose**: [UR-10]에 따라 HTMX를 사용하여 동적인 UI를 구현할 때, 서버로부터의 부분 HTML 응답 및 클라이언트 DOM 반영 과정이 사용자에게 지연 없이 즉각적으로 느껴지도록 하여 최적의 사용자 경험(UX)을 제공하기 위함이다.
- **Scope**: 카운트 관리 웹 UI([FR-04]) 내에서 발생하는 모든 HTMX 기반 비동기 요청 및 부분 화면 갱신 작업.

## 2. Measurement Indicators (Measurable Criteria)
- **Indicator**: UI 상호작용 지연 시간 (User Interaction to UI Update Completion)
  - 사용자의 입력(클릭, 키 입력 등) 시점부터 HTMX 요청 결과가 DOM에 반영되어 화면에 최종 렌더링이 완료된 시점까지의 시간.
- **Target Value**: 
  - **P95 (95th Percentile)**: 300ms 이하
  - **P99 (99th Percentile)**: 500ms 이하
- **Unit/Conditions**: 
  - **측정 단위**: 밀리초(ms)
  - **환경**: 표준 데스크톱 브라우저 환경 (Chrome, Edge 등 Chromium 기반 브라우저 최신 버전)
  - **제외 조건**: 클라이언트의 네트워크 상태 불량(RTT 100ms 초과)으로 인한 지연은 목표치 계산에서 제외할 수 있음.

## 3. Verification Method
- **Measurement Tool/Method**:
  - **HTMX Events Logging**: 브라우저 콘솔에서 `htmx:beforeRequest`와 `htmx:afterOnLoad` 이벤트를 가로채어 시간 차이를 측정하는 스크립트를 활용.
  - **Browser DevTools**: Network 탭의 'Time' 항목(TTFB + Content Download) 및 Performance 탭의 상호작용 분석.
  - **Load Testing Tool (Optional)**: 서버 측의 응답 시간(Latency)이 UI 성능의 병목이 아닌지 검증하기 위해 k6 또는 JMeter 등을 활용한 동시 요청 성능 테스트 병행.
- **Pass Criteria**:
  - 정의된 사용자 시나리오(카운트 생성, 목록 필터링, 값 변경)별로 각 50회 이상 테스트를 수행하여, 전체 결과의 95%가 300ms 이내에 도달해야 함.

## 4. References
- **Related FR**: [FR-04] 카운트 관리 웹 UI
- **Related SC**: [SC-05] UI 기술 스택 (HTMX)
- **Related URD**: [UR-10] HTMX 활용 동적 UI 구현
