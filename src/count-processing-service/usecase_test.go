package domain_test

import (
	"context"
	"count-processing-service/domain"
	"count-processing-service/mocks"
	"sync"
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
		mockHistoryRepo := &mocks.MockCountHistoryRepository{}
		uc := domain.NewCountValueUseCase(mockRepo, mockHistoryRepo)

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
		mockHistoryRepo := &mocks.MockCountHistoryRepository{}
		uc := domain.NewCountValueUseCase(mockRepo, mockHistoryRepo)

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
		mockHistoryRepo := &mocks.MockCountHistoryRepository{}
		uc := domain.NewCountValueUseCase(mockRepo, mockHistoryRepo)

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
		mockHistoryRepo := &mocks.MockCountHistoryRepository{}
		uc := domain.NewCountValueUseCase(mockRepo, mockHistoryRepo)

		res, err := uc.GetMultiple(ctx, []string{"item-1", "item-2"})
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("empty", func(t *testing.T) {
		mockHistoryRepo := &mocks.MockCountHistoryRepository{}
		uc := domain.NewCountValueUseCase(&mocks.MockCountValueRepository{}, mockHistoryRepo)
		res, err := uc.GetMultiple(ctx, []string{})
		assert.NoError(t, err)
		assert.Empty(t, res)
	})
}

func TestCountValueUseCase_GetAll(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expected := []domain.CountValue{
			{ItemID: "item-1", CurrentValue: 10},
			{ItemID: "item-2", CurrentValue: 20},
		}
		mockRepo := &mocks.MockCountValueRepository{
			GetAllFunc: func(ctx context.Context) ([]domain.CountValue, error) {
				return expected, nil
			},
		}
		mockHistoryRepo := &mocks.MockCountHistoryRepository{}
		uc := domain.NewCountValueUseCase(mockRepo, mockHistoryRepo)

		res, err := uc.GetAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
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
		mockHistoryRepo := &mocks.MockCountHistoryRepository{}
		uc := domain.NewCountValueUseCase(mockRepo, mockHistoryRepo)

		err := uc.Delete(ctx, "item-1")
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := &mocks.MockCountValueRepository{
			GetByIDFunc: func(ctx context.Context, itemID string) (*domain.CountValue, error) {
				return nil, nil
			},
		}
		mockHistoryRepo := &mocks.MockCountHistoryRepository{}
		uc := domain.NewCountValueUseCase(mockRepo, mockHistoryRepo)

		err := uc.Delete(ctx, "item-1")
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestCountValueUseCase_Updates(t *testing.T) {
	ctx := context.Background()

	t.Run("increase success", func(t *testing.T) {
		mockRepo := &mocks.MockCountValueRepository{
			IncreaseFunc: func(ctx context.Context, itemID string, amount int, source string) (*domain.CountValue, error) {
				assert.Equal(t, "source-A", source)
				return &domain.CountValue{ItemID: itemID, CurrentValue: 11}, nil
			},
		}
		mockHistoryRepo := &mocks.MockCountHistoryRepository{}
		uc := domain.NewCountValueUseCase(mockRepo, mockHistoryRepo)

		res, err := uc.Increase(ctx, "item-1", 1, "source-A")
		assert.NoError(t, err)
		assert.Equal(t, 11, res.CurrentValue)
	})

	t.Run("increase not found", func(t *testing.T) {
		mockRepo := &mocks.MockCountValueRepository{
			IncreaseFunc: func(ctx context.Context, itemID string, amount int, source string) (*domain.CountValue, error) {
				return nil, nil
			},
		}
		mockHistoryRepo := &mocks.MockCountHistoryRepository{}
		uc := domain.NewCountValueUseCase(mockRepo, mockHistoryRepo)

		res, err := uc.Increase(ctx, "item-1", 1, "source-A")
		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Nil(t, res)
	})

	t.Run("decrease success", func(t *testing.T) {
		mockRepo := &mocks.MockCountValueRepository{
			DecreaseFunc: func(ctx context.Context, itemID string, amount int, source string) (*domain.CountValue, error) {
				return &domain.CountValue{ItemID: itemID, CurrentValue: 9}, nil
			},
		}
		mockHistoryRepo := &mocks.MockCountHistoryRepository{}
		uc := domain.NewCountValueUseCase(mockRepo, mockHistoryRepo)

		res, err := uc.Decrease(ctx, "item-1", 1, "source-B")
		assert.NoError(t, err)
		assert.Equal(t, 9, res.CurrentValue)
	})

	t.Run("reset success", func(t *testing.T) {
		mockRepo := &mocks.MockCountValueRepository{
			ResetFunc: func(ctx context.Context, itemID string, source string) (*domain.CountValue, error) {
				return &domain.CountValue{ItemID: itemID, CurrentValue: 0}, nil
			},
		}
		mockHistoryRepo := &mocks.MockCountHistoryRepository{}
		uc := domain.NewCountValueUseCase(mockRepo, mockHistoryRepo)

		res, err := uc.Reset(ctx, "item-1", "source-C")
		assert.NoError(t, err)
		assert.Equal(t, 0, res.CurrentValue)
	})
}

func TestCountValueUseCase_Concurrency(t *testing.T) {
	ctx := context.Background()

	t.Run("100 concurrent increase requests", func(t *testing.T) {
		currentVal := 0
		var mu sync.Mutex
		mockRepo := &mocks.MockCountValueRepository{
			IncreaseFunc: func(ctx context.Context, itemID string, amount int, source string) (*domain.CountValue, error) {
				mu.Lock()
				defer mu.Unlock()
				currentVal += amount
				return &domain.CountValue{ItemID: itemID, CurrentValue: currentVal}, nil
			},
		}
		mockHistoryRepo := &mocks.MockCountHistoryRepository{}
		uc := domain.NewCountValueUseCase(mockRepo, mockHistoryRepo)

		var wg sync.WaitGroup
		numRequests := 100
		for i := 0; i < numRequests; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := uc.Increase(ctx, "item-1", 1, "source-X")
				assert.NoError(t, err)
			}()
		}
		wg.Wait()

		assert.Equal(t, 100, currentVal)
	})
}

func TestCountValueUseCase_GetHistory(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := &mocks.MockCountValueRepository{
			GetByIDFunc: func(ctx context.Context, itemID string) (*domain.CountValue, error) {
				return &domain.CountValue{ItemID: "item-1"}, nil
			},
		}
		expectedLogs := []domain.CountLog{
			{ItemID: "item-1", OperationType: "increase", ChangeAmount: 5, Source: "A"},
		}
		mockHistoryRepo := &mocks.MockCountHistoryRepository{
			GetHistoryFunc: func(ctx context.Context, itemID string) ([]domain.CountLog, error) {
				return expectedLogs, nil
			},
		}
		uc := domain.NewCountValueUseCase(mockRepo, mockHistoryRepo)

		res, err := uc.GetHistory(ctx, "item-1")
		assert.NoError(t, err)
		assert.Equal(t, expectedLogs, res)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := &mocks.MockCountValueRepository{
			GetByIDFunc: func(ctx context.Context, itemID string) (*domain.CountValue, error) {
				return nil, nil
			},
		}
		mockHistoryRepo := &mocks.MockCountHistoryRepository{}
		uc := domain.NewCountValueUseCase(mockRepo, mockHistoryRepo)

		res, err := uc.GetHistory(ctx, "item-1")
		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Nil(t, res)
	})
}
