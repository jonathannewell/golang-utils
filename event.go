package golang_utils

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"reflect"
	"strings"
)

var (
	EventBus = newBus()
)

type Event interface {
	Message() string
	Name() string //?????
	Data() map[string]any
	Get(key string) any
	Error() error
	Domain() string
	GetDomain() any
	Matches(filter string) bool
}

type DefaultEvent struct {
	Msg      string
	TypeName string
	DataMap  map[string]any
	Err      error
	Dmn      string
}

func (e *DefaultEvent) Name() string {
	return e.TypeName
}
func (e *DefaultEvent) Data() map[string]any {
	return e.DataMap
}
func (e *DefaultEvent) Message() string {
	return e.Msg
}
func (e *DefaultEvent) Get(key string) any {
	if v, ok := e.DataMap[key]; ok {
		return v
	}
	return "???"
}

func (e *DefaultEvent) Error() error {
	return e.Err
}
func (e *DefaultEvent) Domain() string {
	return e.Dmn
}

func (e *DefaultEvent) GetDomain() any {
	return e.Get("domain")
}

func (e *DefaultEvent) Matches(filter string) bool {
	if filter == "*" {
		return true
	}

	if strings.HasPrefix(filter, "!") {
		if strings.HasPrefix(e.Name(), strings.TrimLeft(filter, "!")) {
			return false
		}
		return true
	}

	return strings.HasPrefix(e.Name(), filter)
}

type DeadLetter struct {
	event    Event
	err      error
	handlers []Handler
}

func NewDeadLetter(event Event, err error, handlers []Handler) *DeadLetter {
	return &DeadLetter{
		event:    event,
		err:      err,
		handlers: handlers,
	}
}

type Registration struct {
	filter   string
	event    Event
	handlers []Handler
}

func NewRegistration(filter string, event Event, handler Handler) *Registration {
	return &Registration{
		filter:   filter,
		event:    event,
		handlers: []Handler{handler},
	}
}

func (r *Registration) uniqueName() string {
	if r.event == nil {
		return r.filter
	}
	return r.filter + "-" + reflect.TypeOf(r.event).String()
}

type Handler func(event Event) error
type RegistrationHandlers map[string]*Registration
type Bus struct {
	handlers RegistrationHandlers
	sent     int
}

func newBus() *Bus {
	return &Bus{
		handlers: make(RegistrationHandlers),
	}
}

func (b *Bus) Register(filter string, emptyEvent Event, handler Handler) {
	b.RegisterHandler(NewRegistration(filter, emptyEvent, handler))
}

func (b *Bus) RegisterHandler(registration *Registration) {

	if reg, ok := b.handlers[registration.uniqueName()]; ok {
		reg.handlers = append(reg.handlers, registration.handlers...)
	} else {
		b.handlers[registration.uniqueName()] = registration
	}
}

func (b *Bus) Send(event Event) {
	if event == nil {
		return
	}

	//All events should get sent to at least two places. The handling target and the event tab!
	sentCnt := 0

	for _, registration := range b.handlers {
		if registration.event == nil || reflect.TypeOf(registration.event) == reflect.TypeOf(event) {
			if event.Matches(registration.filter) {
				eg := new(errgroup.Group)
				missedHandlers := make([]Handler, 0, len(registration.handlers))
				missedHandlers = append(missedHandlers, registration.handlers...)
				for i, handler := range registration.handlers {
					//Send Event to each registered consumer in separate goroutine!
					eg.Go(
						func() error {
							err := handler(event)
							if err == nil {
								b.sent++
								sentCnt++
								return nil
							}
							missedHandlers = missedHandlers[i:]
							return err
						},
					)
				} //End handler loop
				if err := eg.Wait(); err != nil {
					b.Send(NewDeadLetterEvent(event, err, missedHandlers))
				}
			}
		}
	}

	if sentCnt < 2 {
		b.Send(NewDeadLetterEvent(event, fmt.Errorf("No handler(s) for event %s found", event.Name()), nil))
	}

}

func (b *Bus) Registrations() []*Registration {
	registrations := make([]*Registration, 0)
	for _, registration := range b.handlers {
		registrations = append(registrations, registration)
	}
	return registrations
}

func (b *Bus) ClearRegistrations() {
	b.handlers = make(RegistrationHandlers)
}

func Reset() {
	EventBus = newBus()
}

func (b *Bus) AllHandlers() []Handler {
	handlers := make([]Handler, 0)
	for _, registration := range b.handlers {
		handlers = append(handlers, registration.handlers...)
	}
	return handlers
}

func (b *Bus) Sent() int {
	return b.sent
}

//***************************  BASIC BUILT IN EVENTS **************************************************************//

type LogEvent struct {
	DefaultEvent
}

func NewEmptyLogEvent() *LogEvent {
	return NewLogEventDetailed("", "", nil)
}

func SendAppLogEvent(format string, args ...interface{}) {
	EventBus.Send(NewLogEventDetailed("app", fmt.Sprintf(format, args...), nil))
}

func NewLogEventDetailed(name, msg string, data map[string]any) *LogEvent {
	return &LogEvent{
		DefaultEvent{TypeName: name, Msg: msg, DataMap: data},
	}
}

type DeadLetterEvent struct {
	DefaultEvent
}

func NewEmptyDeadLetterEvent() *DeadLetterEvent {
	return &DeadLetterEvent{}
}

func NewDeadLetterEvent(event Event, err error, handlers []Handler) *DeadLetterEvent {
	return &DeadLetterEvent{
		DefaultEvent{
			TypeName: "dead-letter",
			Msg:      err.Error(),
			DataMap:  map[string]any{"event": event, "handlers": handlers},
		},
	}
}

type ErrorEvent struct {
	DefaultEvent
}

func NewEmptyErrorEvent() *ErrorEvent {
	return &ErrorEvent{
		DefaultEvent{TypeName: "error"},
	}
}

func NewErrorEvent(source string, err error, fmtMsg string, args ...interface{}) *ErrorEvent {
	return &ErrorEvent{
		DefaultEvent{
			TypeName: "error",
			Msg:      fmt.Sprintf(fmtMsg, args...),
			Dmn:      source,
			Err:      err,
		},
	}
}

func SendErrorEvent(source string, err error, msg string, args ...interface{}) {
	EventBus.Send(NewErrorEvent(source, err, msg, args...))
	SendAppLogEvent(msg, args...)
}
