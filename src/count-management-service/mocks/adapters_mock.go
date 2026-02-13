package mocks

import (
	"context"
	"count-management-service/domain"
)

type MockCountItemRepository struct {
	SaveFunc       func(ctx context.Context, item *domain.CountItem) error
	FindAllFunc    func(ctx context.Context) ([]domain.CountItem, error)
	FindByIDFunc   func(ctx context.Context, id string) (*domain.CountItem, error)
	DeleteFunc     func(ctx context.Context, id string) error
	UpdateFunc     func(ctx context.Context, item *domain.CountItem) error
	FindByNameFunc func(ctx context.Context, name string) (*domain.CountItem, error)
}

func (m *MockCountItemRepository) Save(ctx context.Context, item *domain.CountItem) error {
	return m.SaveFunc(ctx, item)
}
func (m *MockCountItemRepository) FindAll(ctx context.Context) ([]domain.CountItem, error) {
	return m.FindAllFunc(ctx)
}
func (m *MockCountItemRepository) FindByID(ctx context.Context, id string) (*domain.CountItem, error) {
	return m.FindByIDFunc(ctx, id)
}
func (m *MockCountItemRepository) Delete(ctx context.Context, id string) error {
	return m.DeleteFunc(ctx, id)
}
func (m *MockCountItemRepository) Update(ctx context.Context, item *domain.CountItem) error {
	return m.UpdateFunc(ctx, item)
}
func (m *MockCountItemRepository) FindByName(ctx context.Context, name string) (*domain.CountItem, error) {
	if m.FindByNameFunc == nil {
		return nil, nil
	}
	return m.FindByNameFunc(ctx, name)
}

type MockValueServiceClient struct {
	InitializeValueFunc func(ctx context.Context, itemId string, initialValue int) error
	DeleteValueFunc     func(ctx context.Context, itemId string) error
}

func (m *MockValueServiceClient) InitializeValue(ctx context.Context, itemId string, initialValue int) error {
	return m.InitializeValueFunc(ctx, itemId, initialValue)
}
func (m *MockValueServiceClient) DeleteValue(ctx context.Context, itemId string) error {
	return m.DeleteValueFunc(ctx, itemId)
}
