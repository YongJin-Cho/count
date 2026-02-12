package storage

import (
	"bufio"
	"count-api-service/internal/common/event"
	"encoding/json"
	"os"
	"testing"
	"time"
)

type mockSubscriber struct {
	ch chan event.CountCollectedEvent
}

func (m *mockSubscriber) Subscribe() <-chan event.CountCollectedEvent {
	return m.ch
}

func TestFileStorage(t *testing.T) {
	tempFile := "test_counts.log"
	defer os.Remove(tempFile)

	s := NewFileStorage(tempFile)
	mSub := &mockSubscriber{ch: make(chan event.CountCollectedEvent, 1)}
	s.Start(mSub)

	ev := event.CountCollectedEvent{
		ExternalID: "test-id",
		Count:      100,
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	mSub.ch <- ev

	// Give it a moment to process
	time.Sleep(100 * time.Millisecond)

	f, err := os.Open(tempFile)
	if err != nil {
		t.Fatalf("Could not open storage file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		t.Fatal("Storage file is empty")
	}

	var stored event.CountCollectedEvent
	if err := json.Unmarshal(scanner.Bytes(), &stored); err != nil {
		t.Fatalf("Could not unmarshal stored data: %v", err)
	}

	if stored.ExternalID != ev.ExternalID || stored.Count != ev.Count {
		t.Errorf("Stored data mismatch. Got %+v, want %+v", stored, ev)
	}
}

func TestFileStorage_FindAll(t *testing.T) {
	tempFile := "test_find_all.log"
	defer os.Remove(tempFile)

	s := NewFileStorage(tempFile)
	events := []event.CountCollectedEvent{
		{ExternalID: "A", Count: 1, Timestamp: "2023-01-01T00:00:00Z"},
		{ExternalID: "A", Count: 2, Timestamp: "2023-01-02T00:00:00Z"},
		{ExternalID: "B", Count: 3, Timestamp: "2023-01-03T00:00:00Z"},
	}

	for _, ev := range events {
		if err := s.persist(ev); err != nil {
			t.Fatalf("Failed to persist event: %v", err)
		}
	}

	tests := []struct {
		name          string
		filter        string
		limit         int
		offset        int
		expectedCount int
	}{
		{"All", "", 10, 0, 3},
		{"Filter A", "A", 10, 0, 2},
		{"Limit 1", "", 1, 0, 1},
		{"Offset 1", "", 10, 1, 2},
		{"Offset 2", "", 10, 2, 1},
		{"Empty results", "C", 10, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := s.FindAll(tt.filter, tt.limit, tt.offset)
			if err != nil {
				t.Fatalf("FindAll failed: %v", err)
			}
			if len(results) != tt.expectedCount {
				t.Errorf("Expected %d results, got %d", tt.expectedCount, len(results))
			}
			if tt.filter != "" {
				for _, r := range results {
					if r.ExternalID != tt.filter {
						t.Errorf("Expected ExternalID %s, got %s", tt.filter, r.ExternalID)
					}
				}
			}
		})
	}
}

func TestFileStorage_CountTotal(t *testing.T) {
	tempFile := "test_count_total.log"
	defer os.Remove(tempFile)

	s := NewFileStorage(tempFile)
	events := []event.CountCollectedEvent{
		{ExternalID: "A", Count: 1, Timestamp: "2023-01-01T00:00:00Z"},
		{ExternalID: "A", Count: 2, Timestamp: "2023-01-02T00:00:00Z"},
		{ExternalID: "B", Count: 3, Timestamp: "2023-01-03T00:00:00Z"},
	}

	for _, ev := range events {
		if err := s.persist(ev); err != nil {
			t.Fatalf("Failed to persist event: %v", err)
		}
	}

	tests := []struct {
		name          string
		filter        string
		expectedCount int
	}{
		{"All", "", 3},
		{"Filter A", "A", 2},
		{"Filter B", "B", 1},
		{"Filter C", "C", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := s.CountTotal(tt.filter)
			if err != nil {
				t.Fatalf("CountTotal failed: %v", err)
			}
			if count != tt.expectedCount {
				t.Errorf("Expected total %d, got %d", tt.expectedCount, count)
			}
		})
	}
}
