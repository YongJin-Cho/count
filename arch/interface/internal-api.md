# Internal API List

This document lists the internal Sync APIs provided or required by the system components.

## Provided APIs

| Interface ID | Method | Path | Name | From | To | Purpose | Parameters | Response |
|--------------|--------|------|------|------|----|---------|------------|----------|
| **CountReadAPI** | GET | (Internal) | Count Read API | CountCollector | CountStorage | Retrieve count data with pagination support (QR-04) | limit, offset | data[], total_count |
