# Internal Event List

This document lists the internal events used for communication between components or services.

## Events

| Interface ID | Name | Producer | Consumer | Purpose | Related User Story |
|--------------|------|----------|----------|---------|---------------------|
| **CountCollectedEvent** | Count Collected Event | CountCollector | CountStorage | Notifies that a count has been collected and needs persistence | FR-01-01 |
