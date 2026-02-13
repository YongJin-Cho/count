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
