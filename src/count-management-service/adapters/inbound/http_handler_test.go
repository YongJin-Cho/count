package inbound_test

import (
	"context"
	"count-management-service/adapters/inbound"
	"count-management-service/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type MockService struct {
	RegisterItemFunc func(ctx context.Context, name, description string) (*domain.CountItem, error)
	ListItemFunc     func(ctx context.Context) ([]domain.CountItem, error)
	DeleteItemFunc   func(ctx context.Context, id string) error
	UpdateItemFunc   func(ctx context.Context, id, name, description string) (*domain.CountItem, error)
}

func (m *MockService) RegisterItem(ctx context.Context, name, description string) (*domain.CountItem, error) {
	return m.RegisterItemFunc(ctx, name, description)
}
func (m *MockService) ListItem(ctx context.Context) ([]domain.CountItem, error) {
	return m.ListItemFunc(ctx)
}
func (m *MockService) DeleteItem(ctx context.Context, id string) error {
	return m.DeleteItemFunc(ctx, id)
}
func (m *MockService) UpdateItem(ctx context.Context, id, name, description string) (*domain.CountItem, error) {
	return m.UpdateItemFunc(ctx, id, name, description)
}

func TestRegisterItemUI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockService{}
	handler := inbound.NewHTTPHandler(mockSvc)

	r := gin.New()
	r.LoadHTMLGlob("../../templates/*.html")
	r.POST("/ui/count-items", handler.RegisterItemUI)

	t.Run("success", func(t *testing.T) {
		mockSvc.RegisterItemFunc = func(ctx context.Context, name, description string) (*domain.CountItem, error) {
			return &domain.CountItem{ID: "123", Name: name, Description: description}, nil
		}

		w := httptest.NewRecorder()
		data := url.Values{}
		data.Set("name", "test")
		data.Set("description", "desc")
		req, _ := http.NewRequest("POST", "/ui/count-items", strings.NewReader(data.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", w.Code)
		}
		if !strings.Contains(w.Body.String(), "id=\"count-item-123\"") {
			t.Error("response should contain item row")
		}
	})

	t.Run("empty name", func(t *testing.T) {
		mockSvc.RegisterItemFunc = func(ctx context.Context, name, description string) (*domain.CountItem, error) {
			return nil, domain.ErrEmptyName
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/ui/count-items", strings.NewReader("name="))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
		if !strings.Contains(w.Body.String(), "name is required") {
			t.Error("response should contain error message")
		}
	})

	t.Run("duplicate name", func(t *testing.T) {
		mockSvc.RegisterItemFunc = func(ctx context.Context, name, description string) (*domain.CountItem, error) {
			return nil, domain.ErrDuplicateName
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/ui/count-items", strings.NewReader("name=Existing"))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusConflict {
			t.Errorf("expected 409, got %d", w.Code)
		}
		if !strings.Contains(w.Body.String(), "item name already exists") {
			t.Error("response should contain error message")
		}
	})
}

func TestListItemsUI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockService{}
	handler := inbound.NewHTTPHandler(mockSvc)

	r := gin.New()
	r.LoadHTMLGlob("../../templates/*.html")
	r.GET("/ui/count-items", handler.ListItemsUI)

	t.Run("success", func(t *testing.T) {
		mockSvc.ListItemFunc = func(ctx context.Context) ([]domain.CountItem, error) {
			return []domain.CountItem{
				{ID: "1", Name: "Item 1", Description: "Desc 1"},
				{ID: "2", Name: "Item 2", Description: "Desc 2"},
			}, nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/ui/count-items", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		if !strings.Contains(w.Body.String(), "Item 1") || !strings.Contains(w.Body.String(), "Item 2") {
			t.Error("response should contain item names")
		}
	})

	t.Run("empty list", func(t *testing.T) {
		mockSvc.ListItemFunc = func(ctx context.Context) ([]domain.CountItem, error) {
			return []domain.CountItem{}, nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/ui/count-items", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		if !strings.Contains(w.Body.String(), "No items found.") {
			t.Error("response should contain empty state message")
		}
	})
}

func TestUpdateItemUI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockService{}
	handler := inbound.NewHTTPHandler(mockSvc)

	r := gin.New()
	r.LoadHTMLGlob("../../templates/*.html")
	r.POST("/ui/count-items/:id", handler.UpdateItemUI)

	t.Run("success", func(t *testing.T) {
		mockSvc.UpdateItemFunc = func(ctx context.Context, id, name, description string) (*domain.CountItem, error) {
			return &domain.CountItem{ID: id, Name: name, Description: description}, nil
		}

		w := httptest.NewRecorder()
		data := url.Values{}
		data.Set("name", "Updated Name")
		req, _ := http.NewRequest("POST", "/ui/count-items/123", strings.NewReader(data.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		if w.Header().Get("HX-Redirect") != "/ui/count-items" {
			t.Error("expected HX-Redirect header")
		}
	})

	t.Run("duplicate name", func(t *testing.T) {
		mockSvc.UpdateItemFunc = func(ctx context.Context, id, name, description string) (*domain.CountItem, error) {
			return nil, domain.ErrDuplicateName
		}

		w := httptest.NewRecorder()
		data := url.Values{}
		data.Set("name", "Existing")
		req, _ := http.NewRequest("POST", "/ui/count-items/123", strings.NewReader(data.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusConflict {
			t.Errorf("expected 409, got %d", w.Code)
		}
	})

	t.Run("empty name", func(t *testing.T) {
		mockSvc.UpdateItemFunc = func(ctx context.Context, id, name, description string) (*domain.CountItem, error) {
			return nil, domain.ErrEmptyName
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/ui/count-items/123", strings.NewReader("name="))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestDeleteItemUI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockService{}
	handler := inbound.NewHTTPHandler(mockSvc)

	r := gin.New()
	r.LoadHTMLGlob("../../templates/*.html")
	r.DELETE("/ui/count-items/:id", handler.DeleteItemUI)

	t.Run("success", func(t *testing.T) {
		mockSvc.DeleteItemFunc = func(ctx context.Context, id string) error {
			return nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/ui/count-items/123", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockSvc.DeleteItemFunc = func(ctx context.Context, id string) error {
			return domain.ErrItemNotFound
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/ui/count-items/non-existent", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestRegisterItemAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockService{}
	handler := inbound.NewHTTPHandler(mockSvc)

	r := gin.New()
	r.POST("/api/v1/count-items", handler.RegisterItemAPI)

	t.Run("success", func(t *testing.T) {
		mockSvc.RegisterItemFunc = func(ctx context.Context, name, description string) (*domain.CountItem, error) {
			return &domain.CountItem{ID: "123", Name: name, Description: description}, nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/count-items", strings.NewReader(`{"name":"test", "description":"desc"}`))
		req.Header.Add("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", w.Code)
		}
	})

	t.Run("duplicate name", func(t *testing.T) {
		mockSvc.RegisterItemFunc = func(ctx context.Context, name, description string) (*domain.CountItem, error) {
			return nil, domain.ErrDuplicateName
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/count-items", strings.NewReader(`{"name":"Inventory"}`))
		req.Header.Add("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusConflict {
			t.Errorf("expected 409, got %d", w.Code)
		}
	})
}

func TestListItemsAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockService{}
	handler := inbound.NewHTTPHandler(mockSvc)

	r := gin.New()
	r.GET("/api/v1/count-items", handler.ListItemsAPI)

	t.Run("success", func(t *testing.T) {
		mockSvc.ListItemFunc = func(ctx context.Context) ([]domain.CountItem, error) {
			return []domain.CountItem{{ID: "1", Name: "Item 1"}}, nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/count-items", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})
}

func TestUpdateItemAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockService{}
	handler := inbound.NewHTTPHandler(mockSvc)

	r := gin.New()
	r.PUT("/api/v1/count-items/:id", handler.UpdateItemAPI)

	t.Run("success", func(t *testing.T) {
		mockSvc.UpdateItemFunc = func(ctx context.Context, id, name, description string) (*domain.CountItem, error) {
			return &domain.CountItem{ID: id, Name: name}, nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/count-items/123", strings.NewReader(`{"name":"New Name"}`))
		req.Header.Add("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockSvc.UpdateItemFunc = func(ctx context.Context, id, name, description string) (*domain.CountItem, error) {
			return nil, domain.ErrItemNotFound
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/count-items/non-existent", strings.NewReader(`{"name":"Test"}`))
		req.Header.Add("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestDeleteItemAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockService{}
	handler := inbound.NewHTTPHandler(mockSvc)

	r := gin.New()
	r.DELETE("/api/v1/count-items/:id", handler.DeleteItemAPI)

	t.Run("success", func(t *testing.T) {
		mockSvc.DeleteItemFunc = func(ctx context.Context, id string) error {
			return nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/count-items/123", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})
}
