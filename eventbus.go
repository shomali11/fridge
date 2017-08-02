package fridge

const (
	eventsBuffer = 1000
)

// NewEventBus creates a new EventBus
func NewEventBus() *EventBus {
	events := make(chan *Event, eventsBuffer)
	eventBus := &EventBus{events: events}

	go func(eventBus *EventBus) {
		for event := range events {
			if eventBus.handleEvent == nil {
				continue
			}

			eventBus.handleEvent(event)
		}
	}(eventBus)

	return eventBus
}

// Event contains information about an event
type Event struct {
	Key  string
	Type string
}

// EventBus is an events channel
type EventBus struct {
	events      chan *Event
	handleEvent func(event *Event)
}

// Publish publishes an event
func (e *EventBus) Publish(key string, eventType string) {
	e.events <- &Event{Key: key, Type: eventType}
}

// HandleEvent overrides the default handleEvent callback
func (e *EventBus) HandleEvent(handleEvent func(event *Event)) {
	e.handleEvent = handleEvent
}
