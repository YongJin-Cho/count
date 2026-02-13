package domain_test

import (
	"context"
	"count-management-service/domain"
	"count-management-service/mocks"
	"errors"
	"testing"
)

func TestRegisterItem(t *testing.T) {
	mockRepo := &mocks.MockCountItemRepository{}
	mockClient := &mocks.MockValueServiceClient{}
	uc := domain.NewCountItemUseCase(mockRepo, mockClient)

	t.Run("success", func(t *testing.T) {
		mockRepo.FindByNameFunc = func(ctx context.Context, name string) (*domain.CountItem, error) {
			return nil, nil
		}
		mockRepo.SaveFunc = func(ctx context.Context, item *domain.CountItem) error {
			return nil
		}
		mockClient.InitializeValueFunc = func(ctx context.Context, itemId string, initialValue int) error {
			return nil
		}

		item, err := uc.RegisterItem(context.Background(), "test", "desc")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if item.Name != "test" {
			t.Errorf("expected name test, got %s", item.Name)
		}
	})

	t.Run("empty name", func(t *testing.T) {
		_, err := uc.RegisterItem(context.Background(), "", "desc")
		if err != domain.ErrEmptyName {
			t.Errorf("expected ErrEmptyName, got %v", err)
		}
	})

	t.Run("duplicate name", func(t *testing.T) {
		mockRepo.FindByNameFunc = func(ctx context.Context, name string) (*domain.CountItem, error) {
			return &domain.CountItem{Name: name}, nil
		}
		_, err := uc.RegisterItem(context.Background(), "test", "desc")
		if err != domain.ErrDuplicateName {
			t.Errorf("expected ErrDuplicateName, got %v", err)
		}
	})

	t.Run("initialization failure should rollback", func(t *testing.T) {
		mockRepo.FindByNameFunc = func(ctx context.Context, name string) (*domain.CountItem, error) {
			return nil, nil
		}
		mockRepo.SaveFunc = func(ctx context.Context, item *domain.CountItem) error {
			return nil
		}
		mockClient.InitializeValueFunc = func(ctx context.Context, itemId string, initialValue int) error {
			return errors.New("remote error")
		}

		deleted := false
		mockRepo.DeleteFunc = func(ctx context.Context, id string) error {
			deleted = true
			return nil
		}

		_, err := uc.RegisterItem(context.Background(), "test", "desc")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !deleted {
			t.Error("expected rollback (delete) to be called")
		}
	})
}

func TestListItem(t *testing.T) {
	mockRepo := &mocks.MockCountItemRepository{}
	mockClient := &mocks.MockValueServiceClient{}
	uc := domain.NewCountItemUseCase(mockRepo, mockClient)

	t.Run("success", func(t *testing.T) {
		mockRepo.FindAllFunc = func(ctx context.Context) ([]domain.CountItem, error) {
			return []domain.CountItem{
				{ID: "1", Name: "Item 1"},
				{ID: "2", Name: "Item 2"},
			}, nil
		}

		items, err := uc.ListItem(context.Background())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(items) != 2 {
			t.Errorf("expected 2 items, got %d", len(items))
		}
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo.FindAllFunc = func(ctx context.Context) ([]domain.CountItem, error) {
			return []domain.CountItem{}, nil
		}

		items, err := uc.ListItem(context.Background())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(items) != 0 {
			t.Errorf("expected 0 items, got %d", len(items))
		}
	})
}

func TestUpdateItem(t *testing.T) {
	mockRepo := &mocks.MockCountItemRepository{}
	mockClient := &mocks.MockValueServiceClient{}
	uc := domain.NewCountItemUseCase(mockRepo, mockClient)

	t.Run("success", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id string) (*domain.CountItem, error) {
			return &domain.CountItem{ID: id, Name: "Old Name"}, nil
		}
		mockRepo.UpdateFunc = func(ctx context.Context, item *domain.CountItem) error {
			return nil
		}

		item, err := uc.UpdateItem(context.Background(), "123", "New Name", "New Desc")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if item.Name != "New Name" {
			t.Errorf("expected name New Name, got %s", item.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id string) (*domain.CountItem, error) {
			return nil, nil
		}

		_, err := uc.UpdateItem(context.Background(), "non-existent", "New Name", "desc")
		if err != domain.ErrItemNotFound {
			t.Errorf("expected ErrItemNotFound, got %v", err)
		}
	})

	t.Run("empty name", func(t *testing.T) {
		_, err := uc.UpdateItem(context.Background(), "123", "", "desc")
		if err != domain.ErrEmptyName {
			t.Errorf("expected ErrEmptyName, got %v", err)
		}
	})

	t.Run("duplicate name", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id string) (*domain.CountItem, error) {
			return &domain.CountItem{ID: id, Name: "Old Name"}, nil
		}
		mockRepo.FindByNameFunc = func(ctx context.Context, name string) (*domain.CountItem, error) {
			return &domain.CountItem{ID: "other", Name: name}, nil
		}

		_, err := uc.UpdateItem(context.Background(), "123", "Existing Name", "desc")
		if err != domain.ErrDuplicateName {
			t.Errorf("expected ErrDuplicateName, got %v", err)
		}
	})

	t.Run("success with same name", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id string) (*domain.CountItem, error) {
			return &domain.CountItem{ID: id, Name: "Old Name"}, nil
		}
		mockRepo.UpdateFunc = func(ctx context.Context, item *domain.CountItem) error {
			return nil
		}

		_, err := uc.UpdateItem(context.Background(), "123", "Old Name", "New Desc")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

func TestDeleteItem(t *testing.T) {
	mockRepo := &mocks.MockCountItemRepository{}
	mockClient := &mocks.MockValueServiceClient{}
	uc := domain.NewCountItemUseCase(mockRepo, mockClient)

	t.Run("success", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id string) (*domain.CountItem, error) {
			return &domain.CountItem{ID: id}, nil
		}
		mockRepo.DeleteFunc = func(ctx context.Context, id string) error {
			return nil
		}
		mockClient.DeleteValueFunc = func(ctx context.Context, itemId string) error {
			return nil
		}

		err := uc.DeleteItem(context.Background(), "123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id string) (*domain.CountItem, error) {
			return nil, nil
		}

		err := uc.DeleteItem(context.Background(), "non-existent")
		if err != domain.ErrItemNotFound {
			t.Errorf("expected ErrItemNotFound, got %v", err)
		}
	})
}
