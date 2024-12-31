package golang_utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestEvent struct {
	DefaultEvent
}

func newEmptyTestEvent() *TestEvent {
	return newTestEventDetailed("", "", nil)
}

func newTestEventDetailed(name, msg string, data map[string]any) *TestEvent {
	return &TestEvent{
		DefaultEvent: DefaultEvent{TypeName: name, Msg: msg, DataMap: data},
	}
}

var emptyEventHandler = func(event Event) error { return nil }

func TestRegistration(t *testing.T) {
	//Register
	Reset()
	EventBus.Register("suzy", newEmptyTestEvent(), emptyEventHandler)
	assert.Equal(t, 1, len(EventBus.Registrations()), "There should be one registration")
}

func TestMulitRegistrationSameFilter(t *testing.T) {
	//Register
	Reset()
	EventBus.Register("suzy", newEmptyTestEvent(), emptyEventHandler)
	EventBus.Register("suzy", newEmptyTestEvent(), emptyEventHandler)

	registrations := EventBus.Registrations()
	assert.Equal(t, 1, len(registrations), "There should be one registration")
	assert.Equal(
		t,
		2,
		len(EventBus.AllHandlers()),
		"There should two handlers registered for the same registration (filter,type)",
	)
}

func TestReceiveEvent(t *testing.T) {
	Reset()
	EventBus.Register(
		"suzy", newEmptyTestEvent(), func(event Event) error {
			if "suzy" != event.Name() {
				return fmt.Errorf("expected suzyq event, got %s", event.Name())
			}
			if "Reg 1" != event.Message() {
				return fmt.Errorf("expected Reg 1, got %s", event.Message())
			}
			return nil
		},
	)
	EventBus.Register(
		"2suzy", newEmptyTestEvent(), func(event Event) error {
			fmt.Println("Reg 2")
			return nil
		},
	)

	failed := 0
	EventBus.Register(
		"*", NewEmptyDeadLetterEvent(), func(event Event) error {
			failed++
			return nil
		},
	)

	EventBus.Send(newTestEventDetailed("suzy", "Reg 1", nil))

	assert.Equal(t, 0, failed, "All messages should have been successfully sent!")
	assert.Equal(t, 1, EventBus.Sent(), "1 message should have been sent!")

}

func TestFailedSendCapturesMissedEventHandlers(t *testing.T) {
	Reset()
	EventBus.Register(
		"suzy", newEmptyTestEvent(), func(event Event) error {
			if "suzy" != event.Name() {
				return fmt.Errorf("expected suzy event, got %s", event.Name())
			}
			if "Reg 1" != event.Message() {
				return fmt.Errorf("expected Reg 1, got %s", event.Message())
			}
			return nil
		},
	)
	EventBus.Register(
		"suzy", newEmptyTestEvent(), func(event Event) error {
			if "suzy" != event.Name() {
				return fmt.Errorf("expected suzy event, got %s", event.Name())
			}
			if "Reg 2" != event.Message() {
				return fmt.Errorf("expected Reg 2, got %s", event.Message())
			}
			return nil
		},
	)

	failed := 0
	var missedHandlers []Handler
	var dlEvents []Event
	EventBus.Register(
		"*", NewEmptyDeadLetterEvent(), func(event Event) error {
			failed++
			dlEvents = append(dlEvents, event)
			missedHandlers = event.Get("handlers").([]Handler)
			return nil
		},
	)

	EventBus.Send(newTestEventDetailed("suzy", "Reg 1", nil))

	assert.Equal(t, 1, len(missedHandlers), "We should have failed to send the event to 1 handler!")
	assert.Equal(t, 1, failed, "1 messages should have failed!")
	assert.Equal(t, 2, EventBus.Sent(), "1 suzy event & 1 dead-letter event should have been sent!")
	assert.Equal(t, "expected Reg 2, got Reg 1", dlEvents[0].Message(), "Expected the correct error message!")

}
