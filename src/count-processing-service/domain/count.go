package domain

import "time"

type CountValue struct {
	ItemID        string    `json:"itemId" db:"item_id"`
	CurrentValue  int       `json:"currentValue" db:"current_value"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt" db:"last_updated_at"`
}
