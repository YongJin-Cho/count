# External API List

This document lists the external Sync APIs provided or required by the system.

## Provided APIs

| Interface ID | Method | Path | Name | Actor | Purpose | Auth | Related User Story |
|--------------|--------|------|------|-------|---------|------|---------------------|
| **CountCollectAPI** | POST | `/api/v1/collect` | Count Collection API | External System | Collect count data from external systems | Bearer Token | FR-01-01, FR-01-02 |
| **CountQueryAPI** | GET | `/api/v1/counts` | Integrated Count Query API | External System | Query integrated totals and per-source count data with pagination (limit, offset) and total_count in response | Bearer Token | FR-03-01, FR-03-02 |
