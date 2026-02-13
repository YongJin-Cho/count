# Changes

## 2026-02-13
- Initial system architecture design for Count Management System.
- Split system into `count-management-service` and `count-processing-service` for scalability.
- Defined `InternalCountValueAPI` for service-to-service communication.
- Added ADR for service separation.
- Added `ExternalCountUpdateAPI` to `count-processing-service` for external high-frequency updates.
