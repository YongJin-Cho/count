package outbound

import (
	"context"
	"count-management-service/domain"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) domain.CountItemRepository {
	_ = db.AutoMigrate(&domain.CountItem{})
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Save(ctx context.Context, item *domain.CountItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *postgresRepository) FindAll(ctx context.Context) ([]domain.CountItem, error) {
	var items []domain.CountItem
	err := r.db.WithContext(ctx).Find(&items).Error
	return items, err
}

func (r *postgresRepository) FindByID(ctx context.Context, id string) (*domain.CountItem, error) {
	var item domain.CountItem
	err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &item, err
}

func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.CountItem{}, "id = ?", id).Error
}

func (r *postgresRepository) Update(ctx context.Context, item *domain.CountItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *postgresRepository) FindByName(ctx context.Context, name string) (*domain.CountItem, error) {
	var item domain.CountItem
	err := r.db.WithContext(ctx).First(&item, "name = ?", name).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &item, err
}
