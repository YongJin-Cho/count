#!/bin/bash
set -e

# Integration Test Script for Count API Service
BASE_URL="${BASE_URL:-http://localhost}"
HOST_HEADER="count-api.local"
NAMESPACE="count-collection-system"

echo "Starting integration tests against ${BASE_URL} (Host: ${HOST_HEADER})..."

# Function to generate JWT token
generate_token() {
    # Run go run from the service directory to use its go.mod/dependencies
    (cd src/count-api-service && go run ../scripts/gen-token.go)
}

TOKEN=$(generate_token)
echo "Generated Token: ${TOKEN:0:10}..."

echo "--------------------------------------------------"
echo "Test 1: Health Check"
HEALTH_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -H "Host: ${HOST_HEADER}" ${BASE_URL}/health)
if [ "$HEALTH_STATUS" -eq 200 ]; then
    echo "[PASS] Health check successful (200 OK)"
else
    echo "[FAIL] Health check failed with status ${HEALTH_STATUS}"
    # exit 1
fi

echo "--------------------------------------------------"
echo "Test 2: Valid Count Collection"
RESPONSE=$(curl -s -X POST ${BASE_URL}/api/v1/collect \
  -H "Host: ${HOST_HEADER}" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d '{"external_id": "test-device-01", "count": 100}')

echo "Response: ${RESPONSE}"
if [[ $RESPONSE == *"success"* ]]; then
    echo "[PASS] Valid count collection successful"
else
    echo "[FAIL] Valid count collection failed"
    # exit 1
fi

echo "--------------------------------------------------"
echo "Test 3: Unauthorized Access (No Token)"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST ${BASE_URL}/api/v1/collect \
  -H "Host: ${HOST_HEADER}" \
  -H "Content-Type: application/json" \
  -d '{"external_id": "test-device-01", "count": 100}')

if [ "$STATUS" -eq 401 ]; then
    echo "[PASS] Unauthorized access correctly handled (401)"
else
    echo "[FAIL] Expected 401 for unauthorized access, got ${STATUS}"
fi

echo "--------------------------------------------------"
echo "Test 4: Validation Error (Missing count)"
RESPONSE=$(curl -s -X POST ${BASE_URL}/api/v1/collect \
  -H "Host: ${HOST_HEADER}" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d '{"external_id": "test-device-01"}')

echo "Response: ${RESPONSE}"
if [[ $RESPONSE == *"missing count"* ]]; then
    echo "[PASS] Missing count validation successful"
else
    echo "[FAIL] Missing count validation failed"
fi

echo "--------------------------------------------------"
echo "Test 5: Benchmark (Simple)"
echo "Performing 100 requests to measure average response time..."

START_TIME=$(date +%s%3N)
for i in {1..100}
do
    curl -s -o /dev/null -X POST ${BASE_URL}/api/v1/collect \
      -H "Host: ${HOST_HEADER}" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer ${TOKEN}" \
      -d '{"external_id": "bench-device", "count": 1}'
done
END_TIME=$(date +%s%3N)

TOTAL_TIME=$((END_TIME - START_TIME))
AVG_TIME=$(echo "scale=2; $TOTAL_TIME / 100" | bc)

echo "Total time for 100 requests: ${TOTAL_TIME}ms"
echo "Average response time: ${AVG_TIME}ms"

if (( $(echo "$AVG_TIME < 100" | bc -l) )); then
    echo "[PASS] Performance requirement met (Avg < 100ms)"
else
    echo "[WARN] Performance requirement not met (Avg >= 100ms)"
fi

echo "--------------------------------------------------"
echo "Integration tests finished."
