package domain

import (
	"context"
	"github.com/google/uuid"
)

type countItemUseCase struct {
	repo        CountItemRepository
	valueClient ValueServiceClient
}

func NewCountItemUseCase(repo CountItemRepository, valueClient ValueServiceClient) *countItemUseCase {
	return &countItemUseCase{
		repo:        repo,
		valueClient: valueClient,
	}
}

func (u *countItemUseCase) RegisterItem(ctx context.Context, name, description string) (*CountItem, error) {
	if name == "" {
		return nil, ErrEmptyName
	}

	existing, _ := u.repo.FindByName(ctx, name)
	if existing != nil {
		return nil, ErrDuplicateName
	}

	item := &CountItem{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
	}

	if err := u.repo.Save(ctx, item); err != nil {
		return nil, err
	}

	if err := u.valueClient.InitializeValue(ctx, item.ID, 0); err != nil {
		_ = u.repo.Delete(ctx, item.ID)
		return nil, err
	}

	return item, nil
}

func (u *countItemUseCase) ListItem(ctx context.Context) ([]CountItem, error) {
	return u.repo.FindAll(ctx)
}

func (u *countItemUseCase) ListItemWithValues(ctx context.Context) ([]CountItemWithValue, error) {
	items, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return []CountItemWithValue{}, nil
	}

	itemIds := make([]string, len(items))
	for i, item := range items {
		itemIds[i] = item.ID
	}

	values, err := u.valueClient.GetValues(ctx, itemIds)
	if err != nil {
		// Even if we fail to get values, we might want to return items with 0 or error.
		// For now, let's return error as per usual.
		return nil, err
	}

	result := make([]CountItemWithValue, len(items))
	for i, item := range items {
		result[i] = CountItemWithValue{
			CountItem: item,
			Value:     values[item.ID],
		}
	}

	return result, nil
}

func (u *countItemUseCase) GetItemValue(ctx context.Context, id string) (int, error) {
	return u.valueClient.GetValue(ctx, id)
}

func (u *countItemUseCase) DeleteItem(ctx context.Context, id string) error {
	item, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if item == nil {
		return ErrItemNotFound
	}

	if err := u.repo.Delete(ctx, id); err != nil {
		return err
	}
	return u.valueClient.DeleteValue(ctx, id)
}

func (u *countItemUseCase) UpdateItem(ctx context.Context, id, name, description string) (*CountItem, error) {
	if name == "" {
		return nil, ErrEmptyName
	}

	item, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrItemNotFound
	}

	if name != item.Name {
		existing, _ := u.repo.FindByName(ctx, name)
		if existing != nil {
			return nil, ErrDuplicateName
		}
	}

	item.Name = name
	item.Description = description

	if err := u.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}
