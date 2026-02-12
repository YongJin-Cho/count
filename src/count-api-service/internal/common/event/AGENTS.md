# Module Specification: event-module

## 1. Overview
- **Role**: Defines internal event structures and provides an interface for the internal event bus.
- **Build Output**: Library (shared package).

## 2. Providing Interfaces
- **Events**:
    - `CountCollectedEvent`: Internal event emitted when a count is successfully collected.
        - Payload: `{ "external_id": string, "count": integer, "timestamp": string }`
        - Ref: `src/interface/event/CountCollectedEvent.json`
- **Bus Interface**:
    - `Publisher`: Interface for publishing events.
    - `Subscriber`: Interface for subscribing to events.

## 3. Functional Requirements
- **Core Logic**:
    - Define the message format for `CountCollected` according to AsyncAPI spec.
    - Provide a mechanism for the `collector` to publish events without knowing the `storage` implementation.
    - Provide a mechanism for the `storage` to subscribe to events.

## 4. Dependencies
- **Reference Modules**: None.
- **Technologies Used**: Go, Internal Event Bus (e.g., Channels or lightweight Pub/Sub).

## 5. Acceptance Tests
- [ ] `CountCollected` event structure matches `CountCollectedEvent.json`.
- [ ] Events published by one component can be received by another component via the bus.
