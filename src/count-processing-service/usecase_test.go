package domain_test

import (
	"context"
	"count-processing-service/domain"
	"count-processing-service/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountValueUseCase_Initialize(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := &mocks.MockCountValueRepository{
			GetByIDFunc: func(ctx context.Context, itemID string) (*domain.CountValue, error) {
				return nil, nil
			},
			CreateFunc: func(ctx context.Context, count *domain.CountValue) error {
				assert.Equal(t, "item-1", count.ItemID)
				assert.Equal(t, 10, count.CurrentValue)
				return nil
			},
		}
		uc := domain.NewCountValueUseCase(mockRepo)

		res, err := uc.Initialize(ctx, "item-1", 10)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "item-1", res.ItemID)
	})

	t.Run("already exists", func(t *testing.T) {
		mockRepo := &mocks.MockCountValueRepository{
			GetByIDFunc: func(ctx context.Context, itemID string) (*domain.CountValue, error) {
				return &domain.CountValue{ItemID: "item-1"}, nil
			},
		}
		uc := domain.NewCountValueUseCase(mockRepo)

		res, err := uc.Initialize(ctx, "item-1", 10)
		assert.ErrorIs(t, err, domain.ErrAlreadyExists)
		assert.Nil(t, res)
	})

	t.Run("repository error on GetByID", func(t *testing.T) {
		mockRepo := &mocks.MockCountValueRepository{
			GetByIDFunc: func(ctx context.Context, itemID string) (*domain.CountValue, error) {
				return nil, assert.AnError
			},
		}
		uc := domain.NewCountValueUseCase(mockRepo)

		res, err := uc.Initialize(ctx, "item-1", 10)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, res)
	})
}

func TestCountValueUseCase_GetMultiple(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expected := []domain.CountValue{
			{ItemID: "item-1", CurrentValue: 10},
			{ItemID: "item-2", CurrentValue: 20},
		}
		mockRepo := &mocks.MockCountValueRepository{
			GetByIDsFunc: func(ctx context.Context, itemIDs []string) ([]domain.CountValue, error) {
				assert.Equal(t, []string{"item-1", "item-2"}, itemIDs)
				return expected, nil
			},
		}
		uc := domain.NewCountValueUseCase(mockRepo)

		res, err := uc.GetMultiple(ctx, []string{"item-1", "item-2"})
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("empty", func(t *testing.T) {
		uc := domain.NewCountValueUseCase(&mocks.MockCountValueRepository{})
		res, err := uc.GetMultiple(ctx, []string{})
		assert.NoError(t, err)
		assert.Empty(t, res)
	})
}

func TestCountValueUseCase_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := &mocks.MockCountValueRepository{
			GetByIDFunc: func(ctx context.Context, itemID string) (*domain.CountValue, error) {
				return &domain.CountValue{ItemID: "item-1"}, nil
			},
			DeleteFunc: func(ctx context.Context, itemID string) error {
				assert.Equal(t, "item-1", itemID)
				return nil
			},
		}
		uc := domain.NewCountValueUseCase(mockRepo)

		err := uc.Delete(ctx, "item-1")
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := &mocks.MockCountValueRepository{
			GetByIDFunc: func(ctx context.Context, itemID string) (*domain.CountValue, error) {
				return nil, nil
			},
		}
		uc := domain.NewCountValueUseCase(mockRepo)

		err := uc.Delete(ctx, "item-1")
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}
