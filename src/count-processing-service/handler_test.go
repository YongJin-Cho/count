package domain_test

import (
	"bytes"
	"count-processing-service/adapters/inbound"
	"count-processing-service/domain"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"context"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCountValueHandler_Initialize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		expected := &domain.CountValue{ItemID: "item-1", CurrentValue: 10}
		mockUC := &MockCountValueUseCase{
			InitializeFunc: func(ctx context.Context, itemID string, initialValue int) (*domain.CountValue, error) {
				assert.Equal(t, "item-1", itemID)
				return expected, nil
			},
		}
		handler := inbound.NewCountValueHandler(mockUC)
		r := gin.New()
		handler.RegisterRoutes(r)

		reqBody, _ := json.Marshal(inbound.InitializeRequest{
			ItemID:       "item-1",
			InitialValue: 10,
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/internal/counts", bytes.NewBuffer(reqBody))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		var actual domain.CountValue
		json.Unmarshal(resp.Body.Bytes(), &actual)
		assert.Equal(t, expected.ItemID, actual.ItemID)
		assert.Equal(t, expected.CurrentValue, actual.CurrentValue)
	})
}

func TestCountValueHandler_GetMultiple(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		counts := []domain.CountValue{
			{ItemID: "item-1", CurrentValue: 10},
		}
		mockUC := &MockCountValueUseCase{
			GetMultipleFunc: func(ctx context.Context, itemIDs []string) ([]domain.CountValue, error) {
				assert.Equal(t, []string{"item-1"}, itemIDs)
				return counts, nil
			},
		}
		handler := inbound.NewCountValueHandler(mockUC)
		r := gin.New()
		handler.RegisterRoutes(r)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/internal/counts?itemIds=item-1", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var actual map[string][]domain.CountValue
		json.Unmarshal(resp.Body.Bytes(), &actual)
		assert.Equal(t, counts, actual["counts"])
	})
}

func TestCountValueHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockUC := &MockCountValueUseCase{
			DeleteFunc: func(ctx context.Context, itemID string) error {
				assert.Equal(t, "item-1", itemID)
				return nil
			},
		}
		handler := inbound.NewCountValueHandler(mockUC)
		r := gin.New()
		handler.RegisterRoutes(r)

		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/internal/counts/item-1", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var actual map[string]string
		json.Unmarshal(resp.Body.Bytes(), &actual)
		assert.Equal(t, "successfully deleted", actual["message"])
	})

	t.Run("not found", func(t *testing.T) {
		mockUC := &MockCountValueUseCase{
			DeleteFunc: func(ctx context.Context, itemID string) error {
				return domain.ErrNotFound
			},
		}
		handler := inbound.NewCountValueHandler(mockUC)
		r := gin.New()
		handler.RegisterRoutes(r)

		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/internal/counts/item-1", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
}

type MockCountValueUseCase struct {
	InitializeFunc  func(ctx context.Context, itemID string, initialValue int) (*domain.CountValue, error)
	GetFunc         func(ctx context.Context, itemID string) (*domain.CountValue, error)
	GetMultipleFunc func(ctx context.Context, itemIDs []string) ([]domain.CountValue, error)
	DeleteFunc      func(ctx context.Context, itemID string) error
}

func (m *MockCountValueUseCase) Initialize(ctx context.Context, itemID string, initialValue int) (*domain.CountValue, error) {
	return m.InitializeFunc(ctx, itemID, initialValue)
}

func (m *MockCountValueUseCase) Get(ctx context.Context, itemID string) (*domain.CountValue, error) {
	return m.GetFunc(ctx, itemID)
}

func (m *MockCountValueUseCase) GetMultiple(ctx context.Context, itemIDs []string) ([]domain.CountValue, error) {
	return m.GetMultipleFunc(ctx, itemIDs)
}

func (m *MockCountValueUseCase) Delete(ctx context.Context, itemID string) error {
	return m.DeleteFunc(ctx, itemID)
}
