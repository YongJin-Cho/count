package event

import (
	"testing"
	"time"
)

func TestEventBus(t *testing.T) {
	bus := NewEventBus()
	sub := bus.Subscribe()

	event := CountCollectedEvent{
		ExternalID: "test-id",
		Count:      42,
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	bus.Publish(event)

	select {
	case received := <-sub:
		if received.ExternalID != event.ExternalID || received.Count != event.Count {
			t.Errorf("Received event does not match. Got %+v, want %+v", received, event)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Timed out waiting for event")
	}
}

func TestEventBus_MultipleSubscribers(t *testing.T) {
	bus := NewEventBus()
	sub1 := bus.Subscribe()
	sub2 := bus.Subscribe()

	event := CountCollectedEvent{
		ExternalID: "test-id",
		Count:      10,
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	bus.Publish(event)

	for i, sub := range []<-chan CountCollectedEvent{sub1, sub2} {
		select {
		case received := <-sub:
			if received.ExternalID != event.ExternalID {
				t.Errorf("Subscriber %d: received event does not match", i+1)
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Subscriber %d: timed out", i+1)
		}
	}
}
