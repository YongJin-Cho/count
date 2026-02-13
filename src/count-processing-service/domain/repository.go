package domain

import (
	"context"
)

type CountValueRepository interface {
	Init(ctx context.Context) error
	Create(ctx context.Context, count *CountValue) error
	GetByID(ctx context.Context, itemID string) (*CountValue, error)
	GetByIDs(ctx context.Context, itemIDs []string) ([]CountValue, error)
	Delete(ctx context.Context, itemID string) error
	Increase(ctx context.Context, itemID string, amount int) (*CountValue, error)
	Decrease(ctx context.Context, itemID string, amount int) (*CountValue, error)
	Reset(ctx context.Context, itemID string) (*CountValue, error)
}
