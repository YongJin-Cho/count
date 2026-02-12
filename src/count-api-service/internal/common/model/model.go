package model

import (
	"errors"
	"time"
)

// CountRequest is the DTO for incoming API requests.
type CountRequest struct {
	ExternalID string `json:"external_id"`
	Count      *int   `json:"count"`
}

// Validate checks if the request is valid.
func (r *CountRequest) Validate() error {
	if r.ExternalID == "" {
		return errors.New("missing external_id")
	}
	if r.Count == nil {
		return errors.New("missing count")
	}
	if *r.Count < 0 {
		return errors.New("invalid count value")
	}
	return nil
}

// CountData is the internal domain model for processing and storage.
type CountData struct {
	ExternalID string    `json:"external_id"`
	Count      int       `json:"count"`
	Timestamp  time.Time `json:"timestamp"`
}
