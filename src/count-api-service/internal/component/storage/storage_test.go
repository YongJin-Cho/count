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
