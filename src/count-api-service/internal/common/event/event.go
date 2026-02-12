package event

import (
	"sync"
)

// CountCollectedEvent is the internal event emitted when a count is successfully collected.
type CountCollectedEvent struct {
	ExternalID string `json:"external_id"`
	Count      int    `json:"count"`
	Timestamp  string `json:"timestamp"`
}

// Publisher is the interface for publishing events.
type Publisher interface {
	Publish(event CountCollectedEvent)
}

// Subscriber is the interface for subscribing to events.
type Subscriber interface {
	Subscribe() <-chan CountCollectedEvent
}

// EventBus is an implementation of Publisher and Subscriber using Go channels.
type EventBus struct {
	mu          sync.RWMutex
	subscribers []chan CountCollectedEvent
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make([]chan CountCollectedEvent, 0),
	}
}

func (b *EventBus) Publish(event CountCollectedEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, ch := range b.subscribers {
		// Non-blocking send to avoid hanging if subscriber is slow
		select {
		case ch <- event:
		default:
		}
	}
}

func (b *EventBus) Subscribe() <-chan CountCollectedEvent {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan CountCollectedEvent, 100)
	b.subscribers = append(b.subscribers, ch)
	return ch
}
