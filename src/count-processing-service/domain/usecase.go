package domain

import (
	"context"
)

type CountValueUseCase interface {
	Initialize(ctx context.Context, itemID string, initialValue int) (*CountValue, error)
	Get(ctx context.Context, itemID string) (*CountValue, error)
	GetMultiple(ctx context.Context, itemIDs []string) ([]CountValue, error)
	GetAll(ctx context.Context) ([]CountValue, error)
	Delete(ctx context.Context, itemID string) error
	Increase(ctx context.Context, itemID string, amount int, source string) (*CountValue, error)
	Decrease(ctx context.Context, itemID string, amount int, source string) (*CountValue, error)
	Reset(ctx context.Context, itemID string, source string) (*CountValue, error)
	GetHistory(ctx context.Context, itemID string) ([]CountLog, error)
}

type countValueUseCase struct {
	repo        CountValueRepository
	historyRepo CountHistoryRepository
}

func NewCountValueUseCase(repo CountValueRepository, historyRepo CountHistoryRepository) CountValueUseCase {
	return &countValueUseCase{repo: repo, historyRepo: historyRepo}
}

func (u *countValueUseCase) Initialize(ctx context.Context, itemID string, initialValue int) (*CountValue, error) {
	existing, err := u.repo.GetByID(ctx, itemID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrAlreadyExists
	}
	count := &CountValue{
		ItemID:       itemID,
		CurrentValue: initialValue,
	}
	err = u.repo.Create(ctx, count)
	if err != nil {
		return nil, err
	}
	return count, nil
}

func (u *countValueUseCase) Get(ctx context.Context, itemID string) (*CountValue, error) {
	val, err := u.repo.GetByID(ctx, itemID)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, ErrNotFound
	}
	return val, nil
}

func (u *countValueUseCase) GetMultiple(ctx context.Context, itemIDs []string) ([]CountValue, error) {
	if len(itemIDs) == 0 {
		return []CountValue{}, nil
	}
	return u.repo.GetByIDs(ctx, itemIDs)
}

func (u *countValueUseCase) GetAll(ctx context.Context) ([]CountValue, error) {
	return u.repo.GetAll(ctx)
}

func (u *countValueUseCase) Delete(ctx context.Context, itemID string) error {
	existing, err := u.repo.GetByID(ctx, itemID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNotFound
	}
	return u.repo.Delete(ctx, itemID)
}

func (u *countValueUseCase) Increase(ctx context.Context, itemID string, amount int, source string) (*CountValue, error) {
	val, err := u.repo.Increase(ctx, itemID, amount, source)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, ErrNotFound
	}
	return val, nil
}

func (u *countValueUseCase) Decrease(ctx context.Context, itemID string, amount int, source string) (*CountValue, error) {
	val, err := u.repo.Decrease(ctx, itemID, amount, source)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, ErrNotFound
	}
	return val, nil
}

func (u *countValueUseCase) Reset(ctx context.Context, itemID string, source string) (*CountValue, error) {
	val, err := u.repo.Reset(ctx, itemID, source)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, ErrNotFound
	}
	return val, nil
}

func (u *countValueUseCase) GetHistory(ctx context.Context, itemID string) ([]CountLog, error) {
	existing, err := u.repo.GetByID(ctx, itemID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrNotFound
	}
	return u.historyRepo.GetHistory(ctx, itemID)
}
