#!/bin/bash
set -e

GATEWAY_URL=${1:-"http://localhost"}
MGMT_URL=${MGMT_URL:-$GATEWAY_URL}
PROC_URL=${PROC_URL:-$GATEWAY_URL}

echo "Running integration tests..."
echo "Management URL: $MGMT_URL"
echo "Processing URL: $PROC_URL"

# 0. Check connectivity
echo "Checking connectivity..."
if ! curl -sSf "$MGMT_URL/api/v1/count-items" > /dev/null 2>&1; then
  echo "ERROR: Cannot reach $MGMT_URL/api/v1/count-items"
  exit 1
fi
echo "Connectivity OK."

# 1. Register a count item via API
echo "Testing Register Item API..."
REGISTER_RESPONSE=$(curl -sS -X POST "$MGMT_URL/api/v1/count-items" \
  -H "Content-Type: application/json" \
  -d '{"name": "Integration Test Item", "description": "Created by integration test"}')

echo "Register response: $REGISTER_RESPONSE"
ITEM_ID=$(echo "$REGISTER_RESPONSE" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')

if [ -z "$ITEM_ID" ]; then
  echo "FAILURE: Could not get Item ID from register response. Response was: $REGISTER_RESPONSE"
  exit 1
fi
echo "Item ID: $ITEM_ID"

# 2. List count items via API
echo "Testing List Items API..."
LIST_RESPONSE=$(curl -s -X GET "$MGMT_URL/api/v1/count-items")
if [[ "$LIST_RESPONSE" == *"$ITEM_ID"* ]]; then
  echo "SUCCESS: Item found in list."
else
  echo "FAILURE: Item not found in list."
  exit 1
fi

# 3. Update count item via API
echo "Testing Update Item API..."
UPDATE_RESPONSE=$(curl -s -X PUT "$MGMT_URL/api/v1/count-items/$ITEM_ID" \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Integration Item", "description": "Updated by integration test"}')
echo "Update response: $UPDATE_RESPONSE"

# Verify update in list
LIST_RESPONSE=$(curl -s -X GET "$MGMT_URL/api/v1/count-items")
if [[ "$LIST_RESPONSE" == *"Updated Integration Item"* ]]; then
  echo "SUCCESS: Item updated in list."
else
  echo "FAILURE: Item not updated in list."
  exit 1
fi

# 4. Test UI Registration (HTMX)
echo "Testing UI Register Item (HTMX)..."
UI_REGISTER_RESPONSE=$(curl -s -X POST "$MGMT_URL/ui/count-items" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "name=UI Integration Item&description=Created by UI test")

if [[ "$UI_REGISTER_RESPONSE" == *"UI Integration Item"* ]] && [[ "$UI_REGISTER_RESPONSE" == *"<tr"* ]]; then
  echo "SUCCESS: UI returned HTML fragment for new item."
else
  echo "FAILURE: UI did not return correct HTML fragment."
  echo "Response: $UI_REGISTER_RESPONSE"
  exit 1
fi

# 5. Delete count item via UI (HTMX)
echo "Testing UI Delete Item (HTMX)..."
UI_ITEM_ID=$(echo "$UI_REGISTER_RESPONSE" | sed -n 's/.*id="count-item-\([^"]*\)".*/\1/p')
if [ -z "$UI_ITEM_ID" ]; then
  echo "FAILURE: Could not extract UI Item ID from response: $UI_REGISTER_RESPONSE"
  exit 1
fi
echo "UI Item ID: $UI_ITEM_ID"

UI_DELETE_RESPONSE=$(curl -s -X DELETE "$MGMT_URL/ui/count-items/$UI_ITEM_ID")
echo "UI Delete response: $UI_DELETE_RESPONSE"

# Verify deletion in list
LIST_RESPONSE=$(curl -s -X GET "$MGMT_URL/api/v1/count-items")
if [[ "$LIST_RESPONSE" == *"$UI_ITEM_ID"* ]]; then
  echo "FAILURE: UI Item still found in list after delete."
  exit 1
else
  echo "SUCCESS: UI Item removed from list."
fi

# 6. Delete count item via API
echo "Testing API Delete Item..."
API_DELETE_RESPONSE=$(curl -s -X DELETE "$MGMT_URL/api/v1/count-items/$ITEM_ID")
echo "API Delete response: $API_DELETE_RESPONSE"

# Verify deletion in list
LIST_RESPONSE=$(curl -s -X GET "$MGMT_URL/api/v1/count-items")
if [[ "$LIST_RESPONSE" == *"$ITEM_ID"* ]]; then
  echo "FAILURE: API Item still found in list after delete."
  exit 1
else
  echo "SUCCESS: API Item removed from list."
fi

# 7. Test External Count Update API
echo "Testing External Count Update API..."

# Create a new item for testing updates
UNIQUE_NAME="Update Test Item $(date +%s)"
REGISTER_RESPONSE=$(curl -sS -X POST "$MGMT_URL/api/v1/count-items" \
  -H "Content-Type: application/json" \
  -d "{\"name\": \"$UNIQUE_NAME\", \"description\": \"For testing External API\"}")
echo "Register response: $REGISTER_RESPONSE"
ITEM_ID=$(echo "$REGISTER_RESPONSE" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
echo "Test Item ID: $ITEM_ID"

# 7.1 Increase
echo "Testing Increase..."
INC_RESPONSE=$(curl -sS -X POST "$PROC_URL/api/v1/counts/$ITEM_ID/increase" \
  -H "Content-Type: application/json" \
  -d '{"amount": 5}')
echo "Increase response: $INC_RESPONSE"
if [[ "$INC_RESPONSE" == *"\"value\":5"* ]]; then
  echo "SUCCESS: Value increased to 5."
else
  echo "FAILURE: Value not increased correctly. Expected 5 in response."
  exit 1
fi

# 7.2 Decrease
echo "Testing Decrease..."
DEC_RESPONSE=$(curl -sS -X POST "$PROC_URL/api/v1/counts/$ITEM_ID/decrease" \
  -H "Content-Type: application/json" \
  -d '{"amount": 2}')
echo "Decrease response: $DEC_RESPONSE"
if [[ "$DEC_RESPONSE" == *"\"value\":3"* ]]; then
  echo "SUCCESS: Value decreased to 3."
else
  echo "FAILURE: Value not decreased correctly. Expected 3 in response."
  exit 1
fi

# 7.3 Reset
echo "Testing Reset..."
RESET_RESPONSE=$(curl -sS -X POST "$PROC_URL/api/v1/counts/$ITEM_ID/reset")
echo "Reset response: $RESET_RESPONSE"
if [[ "$RESET_RESPONSE" == *"\"value\":0"* ]]; then
  echo "SUCCESS: Value reset to 0."
else
  echo "FAILURE: Value not reset correctly. Expected 0 in response."
  exit 1
fi

# 7.4 404 for non-existent item
echo "Testing 404 for non-existent item..."
STATUS_404=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$PROC_URL/api/v1/counts/non-existent-id/increase")
if [ "$STATUS_404" == "404" ]; then
  echo "SUCCESS: Got 404 for non-existent item."
else
  echo "FAILURE: Did not get 404 for non-existent item (got $STATUS_404)."
  exit 1
fi

# 7.5 Concurrency test (10 rapid calls)
echo "Testing atomicity with 10 rapid increase calls..."
# Ensure it starts at 0
curl -s -X POST "$PROC_URL/api/v1/counts/$ITEM_ID/reset" > /dev/null

for i in {1..10}; do
  curl -s -X POST "$PROC_URL/api/v1/counts/$ITEM_ID/increase" \
    -H "Content-Type: application/json" \
    -d '{"amount": 1}' > /dev/null &
done
wait

echo "Rapid calls complete. Checking final value..."
FINAL_INC_RESPONSE=$(curl -s -X POST "$PROC_URL/api/v1/counts/$ITEM_ID/increase" \
  -H "Content-Type: application/json" \
  -d '{"amount": 1}')
echo "Final increase response (after 10 concurrent + 1 sync): $FINAL_INC_RESPONSE"

if [[ "$FINAL_INC_RESPONSE" == *"\"value\":11"* ]]; then
  echo "SUCCESS: Atomicity verified. Final value is 11."
else
  echo "FAILURE: Atomicity failed or inconsistency detected. Expected 11, got response: $FINAL_INC_RESPONSE"
  exit 1
fi

# 8. Test API Retrieval
echo "Testing API Retrieval..."
RETRIEVAL_RESPONSE=$(curl -s "$PROC_URL/api/v1/counts/$ITEM_ID/value")
echo "Retrieval response: $RETRIEVAL_RESPONSE"
if [[ "$RETRIEVAL_RESPONSE" == *"\"currentValue\":11"* ]]; then
  echo "SUCCESS: Retrieved correct value via API."
else
  echo "FAILURE: Incorrect value retrieved via API."
  exit 1
fi

# 9. Test Bulk API Retrieval
echo "Testing Bulk API Retrieval..."
BULK_RETRIEVAL_RESPONSE=$(curl -s "$PROC_URL/api/v1/counts/values")
if [[ "$BULK_RETRIEVAL_RESPONSE" == *"$ITEM_ID"* ]] && [[ "$BULK_RETRIEVAL_RESPONSE" == *"11"* ]]; then
  echo "SUCCESS: Item found in bulk retrieval."
else
  echo "FAILURE: Item or correct value not found in bulk retrieval."
  exit 1
fi

# 10. Test UI Retrieval (HTMX)
echo "Testing UI Retrieval (HTMX)..."
UI_RETRIEVAL_RESPONSE=$(curl -s "$MGMT_URL/ui/counts/$ITEM_ID/value")
echo "UI Retrieval response: $UI_RETRIEVAL_RESPONSE"
if [[ "$UI_RETRIEVAL_RESPONSE" == *"11"* ]]; then
  echo "SUCCESS: Retrieved correct value via UI fragment."
else
  echo "FAILURE: Incorrect value retrieved via UI fragment."
  exit 1
fi

# 11. Test History Logging and Inquiry
echo "Testing History Logging and Inquiry..."

# The previous operations (increase to 11, etc.) should have already generated history logs.
# Let's generate one more specific log for clear verification.
echo "Generating an additional log entry..."
curl -s -X POST "$PROC_URL/api/v1/counts/$ITEM_ID/increase" \
  -H "Content-Type: application/json" \
  -d '{"amount": 5}' > /dev/null

# 11.1 Verify History via API (Processing Service)
echo "Checking history via API..."
HISTORY_API_RESPONSE=$(curl -s "$PROC_URL/api/v1/counts/$ITEM_ID/history")
echo "History API response: $HISTORY_API_RESPONSE"
# We expect to find 'increase' and amount '5' in the history.
if [[ "$HISTORY_API_RESPONSE" == *"increase"* ]] && [[ "$HISTORY_API_RESPONSE" == *"5"* ]]; then
  echo "SUCCESS: History log found in API response."
else
  echo "FAILURE: History log not found or incorrect in API response."
  exit 1
fi

# 11.2 Verify History via UI (Management Service)
echo "Checking history via UI..."
HISTORY_UI_RESPONSE=$(curl -s "$MGMT_URL/ui/counts/$ITEM_ID/history")
# The UI should return HTML fragment with history details.
if [[ "$HISTORY_UI_RESPONSE" == *"increase"* ]] && [[ "$HISTORY_UI_RESPONSE" == *"5"* ]]; then
  echo "SUCCESS: History log found in UI response."
else
  echo "FAILURE: History log not found or incorrect in UI response."
  echo "Response: $HISTORY_UI_RESPONSE"
  exit 1
fi

# Clean up
curl -s -X DELETE "$MGMT_URL/api/v1/count-items/$ITEM_ID" > /dev/null

echo "Integration tests passed!"
