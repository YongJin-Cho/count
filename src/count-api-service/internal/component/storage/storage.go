package storage

import (
	"bufio"
	"count-api-service/internal/common/event"
	"count-api-service/internal/common/model"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type FileStorage struct {
	filePath string
	mu       sync.Mutex
}

func NewFileStorage(filePath string) *FileStorage {
	return &FileStorage{
		filePath: filePath,
	}
}

func (s *FileStorage) Start(subscriber event.Subscriber) {
	ch := subscriber.Subscribe()
	go func() {
		for ev := range ch {
			if err := s.persist(ev); err != nil {
				log.Printf("Failed to persist event: %v", err)
			}
		}
	}()
}

func (s *FileStorage) persist(ev event.CountCollectedEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.OpenFile(s.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer f.Close()

	data, err := json.Marshal(ev)
	if err != nil {
		return fmt.Errorf("could not marshal event: %w", err)
	}

	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}

	return nil
}

// FindAll retrieves paginated count records, optionally filtered by external_id.
// It returns the latest value for each external_id.
func (s *FileStorage) FindAll(filter string, limit int, offset int) ([]model.CountItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.Open(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []model.CountItem{}, nil
		}
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer f.Close()

	// Use a map to keep track of the latest record for each external_id
	latestCounts := make(map[string]model.CountItem)
	var orderedIDs []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var ev event.CountCollectedEvent
		if err := json.Unmarshal(scanner.Bytes(), &ev); err != nil {
			continue // Skip malformed lines
		}

		if filter != "" && ev.ExternalID != filter {
			continue
		}

		item := model.CountItem{
			ExternalID: ev.ExternalID,
			Count:      ev.Count,
			UpdatedAt:  ev.Timestamp,
		}

		if _, exists := latestCounts[ev.ExternalID]; !exists {
			orderedIDs = append(orderedIDs, ev.ExternalID)
		}
		latestCounts[ev.ExternalID] = item
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	// Apply pagination on the unique list
	var result []model.CountItem
	for i := offset; i < len(orderedIDs) && len(result) < limit; i++ {
		result = append(result, latestCounts[orderedIDs[i]])
	}

	return result, nil
}

// FindById retrieves a specific count record by its Source ID (ExternalID).
func (s *FileStorage) FindById(id string) (*model.CountItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.Open(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("record not found: %s", id)
		}
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer f.Close()

	var latest *model.CountItem
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var ev event.CountCollectedEvent
		if err := json.Unmarshal(scanner.Bytes(), &ev); err != nil {
			continue
		}

		if ev.ExternalID == id {
			latest = &model.CountItem{
				ExternalID: ev.ExternalID,
				Count:      ev.Count,
				UpdatedAt:  ev.Timestamp,
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	if latest == nil {
		return nil, fmt.Errorf("record not found: %s", id)
	}

	return latest, nil
}

// CountTotal returns the total number of unique records matching the filter.
func (s *FileStorage) CountTotal(filter string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.Open(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("could not open file: %w", err)
	}
	defer f.Close()

	uniqueIDs := make(map[string]bool)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var ev event.CountCollectedEvent
		if err := json.Unmarshal(scanner.Bytes(), &ev); err != nil {
			continue
		}

		if filter != "" && ev.ExternalID != filter {
			continue
		}

		uniqueIDs[ev.ExternalID] = true
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scanner error: %w", err)
	}

	return len(uniqueIDs), nil
}

// Create persists a new count record.
func (s *FileStorage) Create(item model.CountItem) error {
	ev := event.CountCollectedEvent{
		ExternalID: item.ExternalID,
		Count:      item.Count,
		Timestamp:  time.Now().Format(time.RFC3339),
	}
	if item.UpdatedAt != "" {
		ev.Timestamp = item.UpdatedAt
	}
	return s.persist(ev)
}

// UpdateValue updates the count value for a specific Source ID.
func (s *FileStorage) UpdateValue(id string, value int) error {
	ev := event.CountCollectedEvent{
		ExternalID: id,
		Count:      value,
		Timestamp:  time.Now().Format(time.RFC3339),
	}
	return s.persist(ev)
}
