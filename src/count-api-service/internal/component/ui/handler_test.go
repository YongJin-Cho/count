package ui

import (
	"count-api-service/internal/common/auth"
	"count-api-service/internal/common/model"
	"count-api-service/internal/component/collector"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateTestToken(permissions []string) string {
	claims := &auth.Claims{
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(auth.SecretKey)
	return "Bearer " + tokenString
}

type mockUIStorage struct {
	items []model.CountItem
}

func (m *mockUIStorage) FindAll(filter string, limit int, offset int) ([]model.CountItem, error) {
	return m.items, nil
}

func (m *mockUIStorage) FindById(id string) (*model.CountItem, error) {
	for _, item := range m.items {
		if item.ExternalID == id {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (m *mockUIStorage) CountTotal(filter string) (int, error) {
	return len(m.items), nil
}

func (m *mockUIStorage) Create(item model.CountItem) error {
	m.items = append(m.items, item)
	return nil
}

func (m *mockUIStorage) UpdateValue(id string, value int) error {
	for i, item := range m.items {
		if item.ExternalID == id {
			m.items[i].Count = value
			return nil
		}
	}
	return fmt.Errorf("not found")
}

func TestUIHandler_GetCountList(t *testing.T) {
	// Setup templates
	tmpDir, _ := os.MkdirTemp("", "templates")
	defer os.RemoveAll(tmpDir)

	rowTmpl := `<tr id="count-row-{{.ExternalID}}"><td>{{.ExternalID}}</td><td>{{.Count}}</td></tr>`
	os.WriteFile(filepath.Join(tmpDir, "count_row.html"), []byte(rowTmpl), 0644)
	// dummy others
	os.WriteFile(filepath.Join(tmpDir, "create_form.html"), []byte("form"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "edit_row.html"), []byte("edit"), 0644)

	storage := &mockUIStorage{
		items: []model.CountItem{
			{ExternalID: "test-1", Count: 10},
			{ExternalID: "test-2", Count: 20},
		},
	}
	authProvider := auth.NewAuthProvider()
	coll := collector.NewCollectorHandler(authProvider, nil, storage)
	handler, _ := NewUIHandler(coll, storage, tmpDir, authProvider)

	t.Run("Authorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ui/counts", nil)
		req.Header.Set("Authorization", generateTestToken([]string{"query"}))
		w := httptest.NewRecorder()

		handler.GetCountList(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		if !strings.Contains(body, "test-1") || !strings.Contains(body, "test-2") {
			t.Errorf("Response body missing count items: %s", body)
		}
	})

	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ui/counts", nil)
		w := httptest.NewRecorder()

		handler.GetCountList(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})
}

func TestUIHandler_CreateCount(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "templates")
	defer os.RemoveAll(tmpDir)
	os.WriteFile(filepath.Join(tmpDir, "count_row.html"), []byte(`<tr>{{.ExternalID}}</tr>`), 0644)
	os.WriteFile(filepath.Join(tmpDir, "create_form.html"), []byte("form"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "edit_row.html"), []byte("edit"), 0644)

	storage := &mockUIStorage{
		items: []model.CountItem{{ExternalID: "existing"}},
	}
	authProvider := auth.NewAuthProvider()
	coll := collector.NewCollectorHandler(authProvider, nil, storage)
	handler, _ := NewUIHandler(coll, storage, tmpDir, authProvider)

	token := generateTestToken([]string{"collect"})

	t.Run("Duplicate Error", func(t *testing.T) {
		form := "source_id=existing&initial_value=5"
		req := httptest.NewRequest(http.MethodPost, "/ui/counts", strings.NewReader(form))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", token)
		w := httptest.NewRecorder()

		handler.CreateCount(w, req)

		if w.Code != http.StatusConflict {
			t.Errorf("Expected status 409, got %d", w.Code)
		}
		if w.Header().Get("HX-Retarget") != "#count-create-error-msg" {
			t.Errorf("Expected HX-Retarget header")
		}
		if !strings.Contains(w.Body.String(), "already exists") {
			t.Errorf("Expected error message in response")
		}
	})

	t.Run("Format Error", func(t *testing.T) {
		form := "source_id=INVALID_ID&initial_value=5"
		req := httptest.NewRequest(http.MethodPost, "/ui/counts", strings.NewReader(form))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", token)
		w := httptest.NewRecorder()

		handler.CreateCount(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
		if w.Header().Get("HX-Retarget") != "#count-create-error-msg" {
			t.Errorf("Expected HX-Retarget header")
		}
		if !strings.Contains(w.Body.String(), "invalid source_id format") {
			t.Errorf("Expected format error message")
		}
	})

	t.Run("Unauthorized", func(t *testing.T) {
		form := "source_id=new-id&initial_value=5"
		req := httptest.NewRequest(http.MethodPost, "/ui/counts", strings.NewReader(form))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		handler.CreateCount(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})
}

func TestUIHandler_Increment(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "templates")
	defer os.RemoveAll(tmpDir)
	os.WriteFile(filepath.Join(tmpDir, "count_row.html"), []byte(`<tr>{{.ExternalID}}:{{.Count}}</tr>`), 0644)
	os.WriteFile(filepath.Join(tmpDir, "create_form.html"), []byte("form"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "edit_row.html"), []byte("edit"), 0644)

	storage := &mockUIStorage{
		items: []model.CountItem{{ExternalID: "test", Count: 10}},
	}
	authProvider := auth.NewAuthProvider()
	coll := collector.NewCollectorHandler(authProvider, nil, storage)
	handler, _ := NewUIHandler(coll, storage, tmpDir, authProvider)

	req := httptest.NewRequest(http.MethodPost, "/ui/counts/test/increment", nil)
	req.Header.Set("Authorization", generateTestToken([]string{"collect"}))
	w := httptest.NewRecorder()

	handler.IncrementCount(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "test:11") {
		t.Errorf("Expected incremented value, got %s", w.Body.String())
	}
}
