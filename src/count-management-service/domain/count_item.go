package domain

import (
	"context"
	"errors"
)

var (
	ErrEmptyName     = errors.New("name is required")
	ErrItemNotFound  = errors.New("item not found")
	ErrDuplicateName = errors.New("item name already exists")
)

type CountItem struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"uniqueIndex"`
	Description string `json:"description"`
}

type CountItemRepository interface {
	Save(ctx context.Context, item *CountItem) error
	FindAll(ctx context.Context) ([]CountItem, error)
	FindByID(ctx context.Context, id string) (*CountItem, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, item *CountItem) error
	FindByName(ctx context.Context, name string) (*CountItem, error)
}

type ValueServiceClient interface {
	InitializeValue(ctx context.Context, itemId string, initialValue int) error
	DeleteValue(ctx context.Context, itemId string) error
}
