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

func TestCountValueHandler_ExternalAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("increase success", func(t *testing.T) {
		mockUC := &MockCountValueUseCase{
			IncreaseFunc: func(ctx context.Context, itemID string, amount int, source string) (*domain.CountValue, error) {
				assert.Equal(t, "item-1", itemID)
				assert.Equal(t, 5, amount)
				assert.Equal(t, "source-A", source)
				return &domain.CountValue{ItemID: itemID, CurrentValue: 15}, nil
			},
		}
		handler := inbound.NewCountValueHandler(mockUC)
		r := gin.New()
		handler.RegisterRoutes(r)

		reqBody, _ := json.Marshal(map[string]interface{}{"amount": 5, "source": "source-A"})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/counts/item-1/increase", bytes.NewBuffer(reqBody))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var actual map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &actual)
		assert.Equal(t, "item-1", actual["itemId"])
		assert.Equal(t, float64(15), actual["value"])
	})

	t.Run("increase default amount", func(t *testing.T) {
		mockUC := &MockCountValueUseCase{
			IncreaseFunc: func(ctx context.Context, itemID string, amount int, source string) (*domain.CountValue, error) {
				assert.Equal(t, 1, amount)
				return &domain.CountValue{ItemID: itemID, CurrentValue: 11}, nil
			},
		}
		handler := inbound.NewCountValueHandler(mockUC)
		r := gin.New()
		handler.RegisterRoutes(r)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/counts/item-1/increase", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("decrease success", func(t *testing.T) {
		mockUC := &MockCountValueUseCase{
			DecreaseFunc: func(ctx context.Context, itemID string, amount int, source string) (*domain.CountValue, error) {
				assert.Equal(t, "source-B", source)
				return &domain.CountValue{ItemID: itemID, CurrentValue: 5}, nil
			},
		}
		handler := inbound.NewCountValueHandler(mockUC)
		r := gin.New()
		handler.RegisterRoutes(r)

		reqBody, _ := json.Marshal(map[string]interface{}{"amount": 5, "source": "source-B"})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/counts/item-1/decrease", bytes.NewBuffer(reqBody))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var actual map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &actual)
		assert.Equal(t, float64(5), actual["value"])
	})

	t.Run("reset success", func(t *testing.T) {
		mockUC := &MockCountValueUseCase{
			ResetFunc: func(ctx context.Context, itemID string, source string) (*domain.CountValue, error) {
				assert.Equal(t, "source-C", source)
				return &domain.CountValue{ItemID: itemID, CurrentValue: 0}, nil
			},
		}
		handler := inbound.NewCountValueHandler(mockUC)
		r := gin.New()
		handler.RegisterRoutes(r)

		reqBody, _ := json.Marshal(map[string]string{"source": "source-C"})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/counts/item-1/reset", bytes.NewBuffer(reqBody))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var actual map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &actual)
		assert.Equal(t, float64(0), actual["value"])
	})

	t.Run("not found", func(t *testing.T) {
		mockUC := &MockCountValueUseCase{
			IncreaseFunc: func(ctx context.Context, itemID string, amount int, source string) (*domain.CountValue, error) {
				return nil, domain.ErrNotFound
			},
		}
		handler := inbound.NewCountValueHandler(mockUC)
		r := gin.New()
		handler.RegisterRoutes(r)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/counts/item-1/increase", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("get history success", func(t *testing.T) {
		logs := []domain.CountLog{
			{ItemID: "item-1", OperationType: "increase", ChangeAmount: 5, Source: "source-A"},
		}
		mockUC := &MockCountValueUseCase{
			GetHistoryFunc: func(ctx context.Context, itemID string) ([]domain.CountLog, error) {
				assert.Equal(t, "item-1", itemID)
				return logs, nil
			},
		}
		handler := inbound.NewCountValueHandler(mockUC)
		r := gin.New()
		handler.RegisterRoutes(r)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/counts/item-1/history", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var actual []domain.CountLog
		json.Unmarshal(resp.Body.Bytes(), &actual)
		assert.Len(t, actual, 1)
		assert.Equal(t, "increase", actual[0].OperationType)
	})

	t.Run("get single external success", func(t *testing.T) {
		mockUC := &MockCountValueUseCase{
			GetFunc: func(ctx context.Context, itemID string) (*domain.CountValue, error) {
				return &domain.CountValue{ItemID: "item-1", CurrentValue: 42}, nil
			},
		}
		handler := inbound.NewCountValueHandler(mockUC)
		r := gin.New()
		handler.RegisterRoutes(r)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/counts/item-1/value", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var actual map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &actual)
		assert.Equal(t, "item-1", actual["itemId"])
		assert.Equal(t, float64(42), actual["currentValue"])
	})

	t.Run("get all external success", func(t *testing.T) {
		counts := []domain.CountValue{
			{ItemID: "item-1", CurrentValue: 10},
			{ItemID: "item-2", CurrentValue: 20},
		}
		mockUC := &MockCountValueUseCase{
			GetAllFunc: func(ctx context.Context) ([]domain.CountValue, error) {
				return counts, nil
			},
		}
		handler := inbound.NewCountValueHandler(mockUC)
		r := gin.New()
		handler.RegisterRoutes(r)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/counts/values", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var actual []map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &actual)
		assert.Len(t, actual, 2)
	})
}

type MockCountValueUseCase struct {
	InitializeFunc  func(ctx context.Context, itemID string, initialValue int) (*domain.CountValue, error)
	GetFunc         func(ctx context.Context, itemID string) (*domain.CountValue, error)
	GetMultipleFunc func(ctx context.Context, itemIDs []string) ([]domain.CountValue, error)
	GetAllFunc      func(ctx context.Context) ([]domain.CountValue, error)
	DeleteFunc      func(ctx context.Context, itemID string) error
	IncreaseFunc    func(ctx context.Context, itemID string, amount int, source string) (*domain.CountValue, error)
	DecreaseFunc    func(ctx context.Context, itemID string, amount int, source string) (*domain.CountValue, error)
	ResetFunc       func(ctx context.Context, itemID string, source string) (*domain.CountValue, error)
	GetHistoryFunc  func(ctx context.Context, itemID string) ([]domain.CountLog, error)
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

func (m *MockCountValueUseCase) GetAll(ctx context.Context) ([]domain.CountValue, error) {
	return m.GetAllFunc(ctx)
}

func (m *MockCountValueUseCase) Delete(ctx context.Context, itemID string) error {
	return m.DeleteFunc(ctx, itemID)
}

func (m *MockCountValueUseCase) Increase(ctx context.Context, itemID string, amount int, source string) (*domain.CountValue, error) {
	return m.IncreaseFunc(ctx, itemID, amount, source)
}

func (m *MockCountValueUseCase) Decrease(ctx context.Context, itemID string, amount int, source string) (*domain.CountValue, error) {
	return m.DecreaseFunc(ctx, itemID, amount, source)
}

func (m *MockCountValueUseCase) Reset(ctx context.Context, itemID string, source string) (*domain.CountValue, error) {
	return m.ResetFunc(ctx, itemID, source)
}

func (m *MockCountValueUseCase) GetHistory(ctx context.Context, itemID string) ([]domain.CountLog, error) {
	return m.GetHistoryFunc(ctx, itemID)
}
