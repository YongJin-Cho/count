package mocks

import (
	"context"
	"count-processing-service/domain"
)

type MockCountValueRepository struct {
	CreateFunc   func(ctx context.Context, count *domain.CountValue) error
	GetByIDFunc  func(ctx context.Context, itemID string) (*domain.CountValue, error)
	GetByIDsFunc func(ctx context.Context, itemIDs []string) ([]domain.CountValue, error)
	DeleteFunc   func(ctx context.Context, itemID string) error
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
