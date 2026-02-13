package mocks

import (
	"context"
	"count-processing-service/domain"
)

type MockCountValueRepository struct {
	InitFunc     func(ctx context.Context) error
	CreateFunc   func(ctx context.Context, count *domain.CountValue) error
	GetByIDFunc  func(ctx context.Context, itemID string) (*domain.CountValue, error)
	GetByIDsFunc func(ctx context.Context, itemIDs []string) ([]domain.CountValue, error)
	DeleteFunc   func(ctx context.Context, itemID string) error
	IncreaseFunc func(ctx context.Context, itemID string, amount int) (*domain.CountValue, error)
	DecreaseFunc func(ctx context.Context, itemID string, amount int) (*domain.CountValue, error)
	ResetFunc    func(ctx context.Context, itemID string) (*domain.CountValue, error)
}

func (m *MockCountValueRepository) Init(ctx context.Context) error {
	if m.InitFunc != nil {
		return m.InitFunc(ctx)
	}
	return nil
}

func (m *MockCountValueRepository) Create(ctx context.Context, count *domain.CountValue) error {
	return m.CreateFunc(ctx, count)
}

func (m *MockCountValueRepository) GetByID(ctx context.Context, itemID string) (*domain.CountValue, error) {
	return m.GetByIDFunc(ctx, itemID)
}

func (m *MockCountValueRepository) GetByIDs(ctx context.Context, itemIDs []string) ([]domain.CountValue, error) {
	return m.GetByIDsFunc(ctx, itemIDs)
}

func (m *MockCountValueRepository) Delete(ctx context.Context, itemID string) error {
	return m.DeleteFunc(ctx, itemID)
}

func (m *MockCountValueRepository) Increase(ctx context.Context, itemID string, amount int) (*domain.CountValue, error) {
	return m.IncreaseFunc(ctx, itemID, amount)
}

func (m *MockCountValueRepository) Decrease(ctx context.Context, itemID string, amount int) (*domain.CountValue, error) {
	return m.DecreaseFunc(ctx, itemID, amount)
}

func (m *MockCountValueRepository) Reset(ctx context.Context, itemID string) (*domain.CountValue, error) {
	return m.ResetFunc(ctx, itemID)
}
