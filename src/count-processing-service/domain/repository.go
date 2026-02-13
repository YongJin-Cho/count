package domain

import (
	"context"
)

type CountValueRepository interface {
	Init(ctx context.Context) error
	Create(ctx context.Context, count *CountValue) error
	GetByID(ctx context.Context, itemID string) (*CountValue, error)
	GetByIDs(ctx context.Context, itemIDs []string) ([]CountValue, error)
	GetAll(ctx context.Context) ([]CountValue, error)
	Delete(ctx context.Context, itemID string) error
	Increase(ctx context.Context, itemID string, amount int, source string) (*CountValue, error)
	Decrease(ctx context.Context, itemID string, amount int, source string) (*CountValue, error)
	Reset(ctx context.Context, itemID string, source string) (*CountValue, error)
}

type CountHistoryRepository interface {
	SaveLog(ctx context.Context, log CountLog) error
	GetHistory(ctx context.Context, itemID string) ([]CountLog, error)
}
