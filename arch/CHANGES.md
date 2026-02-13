# Changes

## 2026-02-13
- Initial system architecture design for Count Management System.
- Split system into `count-management-service` and `count-processing-service` for scalability.
- Defined `InternalCountValueAPI` for service-to-service communication.
- Added ADR for service separation.
- Added `CountValueAPI` and `CountItemMonitoringUI` to `CountManagementSystem`.
- Updated `count-management-service` and `ManagementBackend` to provide `CountItemMonitoringUI`.
- Updated `count-processing-service` and `ProcessingBackend` to provide `CountValueAPI` for external value retrieval.
- Updated `CountItemManagementUI` type to `UI` in `system.json` for CDL compliance.
- Added `CountHistoryAPI` and `CountItemHistoryUI` to support audit trail functionality.
- Updated `count-processing-service` and `ProcessingBackend` to handle history logging and provide `CountHistoryAPI`.
- Updated `count-management-service` and `ManagementBackend` to provide `CountItemHistoryUI`.
- Increased resource allocations for `ProcessingBackend` and `ProcessingPostgreSQL` to handle history log processing and storage.
- Added connector for `CountHistoryAPI` between management and processing services in `CountManagementSystem`.
