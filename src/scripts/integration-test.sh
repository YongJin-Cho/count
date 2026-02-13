#!/bin/bash
set -e

GATEWAY_URL=${1:-"http://localhost"}

echo "Running integration tests against $GATEWAY_URL..."

# 0. Check connectivity
echo "Checking connectivity to $GATEWAY_URL..."
if ! curl -sSf "$GATEWAY_URL/api/v1/count-items" > /dev/null 2>&1; then
  echo "ERROR: Cannot reach $GATEWAY_URL/api/v1/count-items"
  echo "If you are running in Kubernetes, make sure you have a Gateway/Ingress controller or use port-forwarding:"
  echo "  kubectl port-forward svc/count-management-service 8888:8080 -n count-system"
  echo "Then run: ./src/scripts/integration-test.sh http://localhost:8888"
  exit 1
fi
echo "Connectivity OK."

# 1. Register a count item via API
echo "Testing Register Item API..."
REGISTER_RESPONSE=$(curl -sS -X POST "$GATEWAY_URL/api/v1/count-items" \
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
LIST_RESPONSE=$(curl -s -X GET "$GATEWAY_URL/api/v1/count-items")
if [[ "$LIST_RESPONSE" == *"$ITEM_ID"* ]]; then
  echo "SUCCESS: Item found in list."
else
  echo "FAILURE: Item not found in list."
  exit 1
fi

# 3. Update count item via API
echo "Testing Update Item API..."
UPDATE_RESPONSE=$(curl -s -X PUT "$GATEWAY_URL/api/v1/count-items/$ITEM_ID" \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Integration Item", "description": "Updated by integration test"}')
echo "Update response: $UPDATE_RESPONSE"

# Verify update in list
LIST_RESPONSE=$(curl -s -X GET "$GATEWAY_URL/api/v1/count-items")
if [[ "$LIST_RESPONSE" == *"Updated Integration Item"* ]]; then
  echo "SUCCESS: Item updated in list."
else
  echo "FAILURE: Item not updated in list."
  exit 1
fi

# 4. Test UI Registration (HTMX)
echo "Testing UI Register Item (HTMX)..."
UI_REGISTER_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/ui/count-items" \
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

UI_DELETE_RESPONSE=$(curl -s -X DELETE "$GATEWAY_URL/ui/count-items/$UI_ITEM_ID")
echo "UI Delete response: $UI_DELETE_RESPONSE"

# Verify deletion in list
LIST_RESPONSE=$(curl -s -X GET "$GATEWAY_URL/api/v1/count-items")
if [[ "$LIST_RESPONSE" == *"$UI_ITEM_ID"* ]]; then
  echo "FAILURE: UI Item still found in list after delete."
  exit 1
else
  echo "SUCCESS: UI Item removed from list."
fi

# 6. Delete count item via API
echo "Testing API Delete Item..."
API_DELETE_RESPONSE=$(curl -s -X DELETE "$GATEWAY_URL/api/v1/count-items/$ITEM_ID")
echo "API Delete response: $API_DELETE_RESPONSE"

# Verify deletion in list
LIST_RESPONSE=$(curl -s -X GET "$GATEWAY_URL/api/v1/count-items")
if [[ "$LIST_RESPONSE" == *"$ITEM_ID"* ]]; then
  echo "FAILURE: API Item still found in list after delete."
  exit 1
else
  echo "SUCCESS: API Item removed from list."
fi

echo "Integration tests passed!"
