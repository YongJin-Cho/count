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

	var counts []model.CountItem
	scanner := bufio.NewScanner(f)
	matchCount := 0
	for scanner.Scan() {
		var ev event.CountCollectedEvent
		if err := json.Unmarshal(scanner.Bytes(), &ev); err != nil {
			continue // Skip malformed lines
		}

		if filter != "" && ev.ExternalID != filter {
			continue
		}

		if matchCount >= offset && len(counts) < limit {
			counts = append(counts, model.CountItem{
				ExternalID: ev.ExternalID,
				Count:      ev.Count,
				UpdatedAt:  ev.Timestamp,
			})
		}
		matchCount++

		if len(counts) >= limit {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return counts, nil
}

// CountTotal returns the total number of records matching the filter.
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

	count := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if filter == "" {
			count++
			continue
		}

		var ev event.CountCollectedEvent
		if err := json.Unmarshal(scanner.Bytes(), &ev); err != nil {
			continue
		}

		if ev.ExternalID == filter {
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scanner error: %w", err)
	}

	return count, nil
}
