package eventbus

import (
	"github.com/shomali11/cmap"
)

const (
	eventsBuffer = 10000
)

// NewClient creates a new eventbus client
func NewClient() *Client {
	events := make(chan *Event, eventsBuffer)
	client := &Client{
		events:   events,
		registry: cmap.NewShardedConcurrentMap(),
	}

	go func(client *Client) {
		for event := range client.events {
			handlers, ok := client.registry.Get(event.Topic)
			if !ok {
				continue
			}

			for _, handler := range handlers.([]EventHandler) {
				handler(event.Value)
			}
		}
	}(client)

	return client
}

// EventHandler an event handler
type EventHandler func(value interface{})

// Event contains information about an event
type Event struct {
	Topic string
	Value interface{}
}

// Client publishes and subscribes to events
type Client struct {
	events   chan *Event
	registry *cmap.ShardedConcurrentMap
}

// Publish publishes an event
func (c *Client) Publish(topic string, value interface{}) {
	event := &Event{Topic: topic, Value: value}

	// Attempt to publish an event to the channel, if the channel's buffer was full
	// Create a go routine to push the event to the channel (which will block until a read on the channel occurs)
	select {
	case c.events <- event:
	default:
		go func(event *Event) {
			c.events <- event
		}(event)
	}
}

// Subscribe subscribes to a topic
func (c *Client) Subscribe(topic string, handler EventHandler) {
	handlers := []EventHandler{}
	existingHandlers, ok := c.registry.Get(topic)
	if ok {
		handlers = existingHandlers.([]EventHandler)
	}

	handlers = append(handlers, handler)
	c.registry.Set(topic, handlers)
}

// Close closes eventbus
func (c *Client) Close() {
	close(c.events)
	c.registry.Clear()
}
