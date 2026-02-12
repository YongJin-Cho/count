package collector

import (
	"bytes"
	"count-api-service/internal/common/auth"
	"count-api-service/internal/common/event"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type mockPublisher struct {
	events []event.CountCollectedEvent
}

func (m *mockPublisher) Publish(ev event.CountCollectedEvent) {
	m.events = append(m.events, ev)
}

func generateToken(permissions []string) string {
	claims := &auth.Claims{
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(auth.SecretKey)
	return tokenString
}

func TestCollectorHandler_CollectCount(t *testing.T) {
	ap := auth.NewAuthProvider()
	pub := &mockPublisher{}
	handler := NewCollectorHandler(ap, pub)

	validToken := generateToken([]string{"collect"})
	forbiddenToken := generateToken([]string{"other"})

	tests := []struct {
		name           string
		token          string
		body           interface{}
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "Missing token",
			token:          "",
			body:           map[string]interface{}{"external_id": "test", "count": 10},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Forbidden token",
			token:          forbiddenToken,
			body:           map[string]interface{}{"external_id": "test", "count": 10},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Missing external_id",
			token:          validToken,
			body:           map[string]interface{}{"count": 10},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "missing external_id",
		},
		{
			name:           "Missing count",
			token:          validToken,
			body:           map[string]interface{}{"external_id": "test"},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "missing count",
		},
		{
			name:           "Invalid count value",
			token:          validToken,
			body:           map[string]interface{}{"external_id": "test", "count": -1},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "invalid count value",
		},
		{
			name:           "Success",
			token:          validToken,
			body:           map[string]interface{}{"external_id": "test", "count": 10},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pub.events = nil // Reset events
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/collect", bytes.NewBuffer(bodyBytes))
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			w := httptest.NewRecorder()

			handler.CollectCount(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedMsg != "" {
				var resp map[string]string
				json.NewDecoder(w.Body).Decode(&resp)
				if resp["message"] != tt.expectedMsg {
					t.Errorf("Expected message %q, got %q", tt.expectedMsg, resp["message"])
				}
			}

			if tt.expectedStatus == http.StatusOK {
				if len(pub.events) != 1 {
					t.Errorf("Expected 1 event to be published, got %d", len(pub.events))
				}
			}
		})
	}
}
