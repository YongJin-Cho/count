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

func NewPostgresHistoryRepository(db *sqlx.DB) domain.CountHistoryRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Init(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS count_values (
		item_id TEXT PRIMARY KEY,
		current_value INTEGER DEFAULT 0,
		last_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS count_history (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		item_id TEXT NOT NULL,
		operation_type TEXT NOT NULL,
		change_amount INTEGER NOT NULL,
		source TEXT NOT NULL,
		timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (item_id) REFERENCES count_values(item_id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_count_history_item_id_timestamp ON count_history(item_id, timestamp DESC);
	`
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

func (r *postgresRepository) Increase(ctx context.Context, itemID string, amount int, source string) (*domain.CountValue, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `UPDATE count_values SET current_value = current_value + $1, last_updated_at = NOW() WHERE item_id = $2 RETURNING item_id, current_value, last_updated_at`
	var count domain.CountValue
	err = tx.GetContext(ctx, &count, query, amount, itemID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	logQuery := `INSERT INTO count_history (item_id, operation_type, change_amount, source, timestamp) VALUES ($1, $2, $3, $4, NOW())`
	_, err = tx.ExecContext(ctx, logQuery, itemID, "increase", amount, source)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &count, nil
}

func (r *postgresRepository) Decrease(ctx context.Context, itemID string, amount int, source string) (*domain.CountValue, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `UPDATE count_values SET current_value = current_value - $1, last_updated_at = NOW() WHERE item_id = $2 RETURNING item_id, current_value, last_updated_at`
	var count domain.CountValue
	err = tx.GetContext(ctx, &count, query, amount, itemID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	logQuery := `INSERT INTO count_history (item_id, operation_type, change_amount, source, timestamp) VALUES ($1, $2, $3, $4, NOW())`
	_, err = tx.ExecContext(ctx, logQuery, itemID, "decrease", -amount, source)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &count, nil
}

func (r *postgresRepository) Reset(ctx context.Context, itemID string, source string) (*domain.CountValue, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// To calculate diff, we need old value
	var oldVal int
	err = tx.GetContext(ctx, &oldVal, "SELECT current_value FROM count_values WHERE item_id = $1 FOR UPDATE", itemID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	query := `UPDATE count_values SET current_value = 0, last_updated_at = NOW() WHERE item_id = $1 RETURNING item_id, current_value, last_updated_at`
	var count domain.CountValue
	err = tx.GetContext(ctx, &count, query, itemID)
	if err != nil {
		return nil, err
	}

	logQuery := `INSERT INTO count_history (item_id, operation_type, change_amount, source, timestamp) VALUES ($1, $2, $3, $4, NOW())`
	_, err = tx.ExecContext(ctx, logQuery, itemID, "reset", -oldVal, source)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &count, nil
}

func (r *postgresRepository) SaveLog(ctx context.Context, log domain.CountLog) error {
	query := `INSERT INTO count_history (item_id, operation_type, change_amount, source, timestamp) VALUES ($1, $2, $3, $4, NOW())`
	_, err := r.db.ExecContext(ctx, query, log.ItemID, log.OperationType, log.ChangeAmount, log.Source)
	return err
}

func (r *postgresRepository) GetHistory(ctx context.Context, itemID string) ([]domain.CountLog, error) {
	query := `SELECT id, item_id, operation_type, change_amount, source, timestamp FROM count_history WHERE item_id = $1 ORDER BY timestamp DESC`
	var logs []domain.CountLog
	err := r.db.SelectContext(ctx, &logs, query, itemID)
	if err != nil {
		return nil, err
	}
	return logs, nil
}
