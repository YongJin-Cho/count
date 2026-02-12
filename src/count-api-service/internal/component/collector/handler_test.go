package collector

import (
	"bytes"
	"count-api-service/internal/common/auth"
	"count-api-service/internal/common/event"
	"count-api-service/internal/common/model"
	"encoding/json"
	"fmt"
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

type mockRepository struct {
	items []model.CountItem
}

func (m *mockRepository) FindAll(filter string, limit int, offset int) ([]model.CountItem, error) {
	var filtered []model.CountItem
	for _, item := range m.items {
		if filter == "" || item.ExternalID == filter {
			filtered = append(filtered, item)
		}
	}

	if offset < 0 {
		offset = 0
	}
	if offset >= len(filtered) {
		return []model.CountItem{}, nil
	}

	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[offset:end], nil
}

func (m *mockRepository) FindById(id string) (*model.CountItem, error) {
	for _, item := range m.items {
		if item.ExternalID == id {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("not found: %s", id)
}

func (m *mockRepository) CountTotal(filter string) (int, error) {
	count := 0
	for _, item := range m.items {
		if filter == "" || item.ExternalID == filter {
			count++
		}
	}
	return count, nil
}

func (m *mockRepository) Create(item model.CountItem) error {
	m.items = append(m.items, item)
	return nil
}

func (m *mockRepository) UpdateValue(id string, value int) error {
	found := false
	for i, item := range m.items {
		if item.ExternalID == id {
			m.items[i].Count = value
			m.items[i].UpdatedAt = time.Now().Format(time.RFC3339)
			found = true
		}
	}
	if !found {
		return fmt.Errorf("not found: %s", id)
	}
	return nil
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
	repo := &mockRepository{}
	handler := NewCollectorHandler(ap, pub, repo)

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

func TestCollectorHandler_GetCounts(t *testing.T) {
	ap := auth.NewAuthProvider()
	pub := &mockPublisher{}
	repo := &mockRepository{
		items: []model.CountItem{
			{ExternalID: "A", Count: 1, UpdatedAt: "2023-01-01T00:00:00Z"},
			{ExternalID: "A", Count: 2, UpdatedAt: "2023-01-02T00:00:00Z"},
			{ExternalID: "B", Count: 3, UpdatedAt: "2023-01-03T00:00:00Z"},
		},
	}
	handler := NewCollectorHandler(ap, pub, repo)

	validToken := generateToken([]string{"query"})
	forbiddenToken := generateToken([]string{"collect"})

	tests := []struct {
		name           string
		token          string
		query          string
		expectedStatus int
		expectedTotal  int
		expectedCount  int
	}{
		{
			name:           "Unauthorized",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Forbidden",
			token:          forbiddenToken,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "All counts",
			token:          validToken,
			query:          "",
			expectedStatus: http.StatusOK,
			expectedTotal:  3,
			expectedCount:  3,
		},
		{
			name:           "Filter by A",
			token:          validToken,
			query:          "external_id=A",
			expectedStatus: http.StatusOK,
			expectedTotal:  2,
			expectedCount:  2,
		},
		{
			name:           "Limit",
			token:          validToken,
			query:          "limit=1",
			expectedStatus: http.StatusOK,
			expectedTotal:  3,
			expectedCount:  1,
		},
		{
			name:           "Offset",
			token:          validToken,
			query:          "offset=2",
			expectedStatus: http.StatusOK,
			expectedTotal:  3,
			expectedCount:  1,
		},
		{
			name:           "Non-existent ID",
			token:          validToken,
			query:          "external_id=C",
			expectedStatus: http.StatusOK,
			expectedTotal:  0,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/counts?"+tt.query, nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			w := httptest.NewRecorder()

			handler.GetCounts(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var resp model.CountListResponse
				json.NewDecoder(w.Body).Decode(&resp)
				if resp.TotalCount != tt.expectedTotal {
					t.Errorf("Expected total count %d, got %d", tt.expectedTotal, resp.TotalCount)
				}
				if len(resp.Counts) != tt.expectedCount {
					t.Errorf("Expected %d items, got %d", tt.expectedCount, len(resp.Counts))
				}
			}
		})
	}
}
