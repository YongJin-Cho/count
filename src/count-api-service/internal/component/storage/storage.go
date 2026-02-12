package storage

import (
	"count-api-service/internal/common/event"
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
