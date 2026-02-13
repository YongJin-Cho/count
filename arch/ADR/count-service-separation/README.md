# ADR: Count Service Separation

## Review Background
The Count Management System needs to handle both metadata management (registering items) and high-frequency count updates (incrementing values). We need to decide whether to combine these into a single service or separate them into distinct microservices.

## Alternatives

### Alternative 1: Single `count-api-service`
All functionality (metadata CRUD and count processing) is handled by a single service.

**Pros:**
- Simpler development and deployment.
- No network overhead between management and processing logic.
- Easier data consistency (single database).
- Lower operational overhead.

**Cons:**
- Scaling management and processing independently is not possible.
- High traffic on count updates could impact the performance of management operations.
- Tighter coupling between different domain concerns.

### Alternative 2: Separated `count-management-service` and `count-processing-service`
Metadata management is handled by `count-management-service`, while count value operations are handled by `count-processing-service`.

**Pros:**
- **Scalability**: `count-processing-service` can be scaled independently to handle high-frequency updates.
- **Isolation**: Performance issues or failures in one service do not directly impact the other.
- **Focus**: Each service has a single, well-defined responsibility.
- **Tech Flexibility**: Processing service could use a different data store (e.g., Redis) optimized for counters in the future.

**Cons:**
- Increased complexity in deployment and operations.
- Network latency for internal calls (e.g., Management service calling Processing service to show count values in the list).
- Distributed data consistency challenges (e.g., creating a count item and initializing its value across two services).

## Comparison Table

| Criteria | Alternative 1 (Single) | Alternative 2 (Separated) |
|----------|-----------------------|---------------------------|
| **Scalability** | Lower | Higher |
| **Complexity** | Lower | Higher |
| **Independence** | Low | High |
| **Operational Effort** | Lower | Higher |

## Recommendation
Given the goal of "thinking about scalability" and the distinct nature of "Management" (metadata CRUD) vs. "Processing" (high-frequency updates), **Alternative 2 (Separated)** is recommended for long-term robustness, even if Alternative 1 is faster to implement initially.

## Decision
User confirmation pending. For this design task, we will proceed with **Alternative 2** to provide a scalable architecture.
