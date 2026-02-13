package domain

import "time"

type CountValue struct {
	ItemID        string    `json:"itemId" db:"item_id"`
	CurrentValue  int       `json:"currentValue" db:"current_value"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt" db:"last_updated_at"`
}

type CountLog struct {
	ID            string    `json:"id" db:"id"`
	ItemID        string    `json:"itemId" db:"item_id"`
	OperationType string    `json:"type" db:"operation_type"`
	ChangeAmount  int       `json:"change" db:"change_amount"`
	Source        string    `json:"source" db:"source"`
	Timestamp     time.Time `json:"timestamp" db:"timestamp"`
}
