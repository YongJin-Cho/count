package ports

import (
	"context"
	"count-management-service/domain"
)

type CountItemService interface {
	RegisterItem(ctx context.Context, name, description string) (*domain.CountItem, error)
	ListItem(ctx context.Context) ([]domain.CountItem, error)
	ListItemWithValues(ctx context.Context) ([]domain.CountItemWithValue, error)
	DeleteItem(ctx context.Context, id string) error
	UpdateItem(ctx context.Context, id, name, description string) (*domain.CountItem, error)
	GetItemValue(ctx context.Context, id string) (int, error)
}
