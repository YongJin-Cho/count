package outbound

import (
	"context"
	"count-processing-service/domain"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) domain.CountValueRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Init(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS count_values (
		item_id TEXT PRIMARY KEY,
		current_value INTEGER DEFAULT 0,
		last_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *postgresRepository) Create(ctx context.Context, count *domain.CountValue) error {
	query := `INSERT INTO count_values (item_id, current_value, last_updated_at) VALUES ($1, $2, NOW())`
	_, err := r.db.ExecContext(ctx, query, count.ItemID, count.CurrentValue)
	return err
}

func (r *postgresRepository) GetByID(ctx context.Context, itemID string) (*domain.CountValue, error) {
	query := `SELECT item_id, current_value, last_updated_at FROM count_values WHERE item_id = $1`
	var count domain.CountValue
	err := r.db.GetContext(ctx, &count, query, itemID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &count, nil
}

func (r *postgresRepository) GetByIDs(ctx context.Context, itemIDs []string) ([]domain.CountValue, error) {
	if len(itemIDs) == 0 {
		return []domain.CountValue{}, nil
	}

	query, args, err := sqlx.In(`SELECT item_id, current_value, last_updated_at FROM count_values WHERE item_id IN (?)`, itemIDs)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)
	var counts []domain.CountValue
	err = r.db.SelectContext(ctx, &counts, query, args...)
	if err != nil {
		return nil, err
	}

	return counts, nil
}

func (r *postgresRepository) GetAll(ctx context.Context) ([]domain.CountValue, error) {
	query := `SELECT item_id, current_value, last_updated_at FROM count_values`
	var counts []domain.CountValue
	err := r.db.SelectContext(ctx, &counts, query)
	if err != nil {
		return nil, err
	}
	return counts, nil
}

func (r *postgresRepository) Delete(ctx context.Context, itemID string) error {
	query := `DELETE FROM count_values WHERE item_id = $1`
	_, err := r.db.ExecContext(ctx, query, itemID)
	return err
}

func (r *postgresRepository) Increase(ctx context.Context, itemID string, amount int) (*domain.CountValue, error) {
	query := `UPDATE count_values SET current_value = current_value + $1, last_updated_at = NOW() WHERE item_id = $2 RETURNING item_id, current_value, last_updated_at`
	var count domain.CountValue
	err := r.db.GetContext(ctx, &count, query, amount, itemID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &count, nil
}

func (r *postgresRepository) Decrease(ctx context.Context, itemID string, amount int) (*domain.CountValue, error) {
	query := `UPDATE count_values SET current_value = current_value - $1, last_updated_at = NOW() WHERE item_id = $2 RETURNING item_id, current_value, last_updated_at`
	var count domain.CountValue
	err := r.db.GetContext(ctx, &count, query, amount, itemID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &count, nil
}

func (r *postgresRepository) Reset(ctx context.Context, itemID string) (*domain.CountValue, error) {
	query := `UPDATE count_values SET current_value = 0, last_updated_at = NOW() WHERE item_id = $1 RETURNING item_id, current_value, last_updated_at`
	var count domain.CountValue
	err := r.db.GetContext(ctx, &count, query, itemID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &count, nil
}
