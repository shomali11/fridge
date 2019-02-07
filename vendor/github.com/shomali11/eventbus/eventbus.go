package eventbus

import (
	"github.com/shomali11/maps"
)

const (
	eventsBuffer = 10000
)

// NewClient creates a new eventbus client
func NewClient() *Client {
	events := make(chan *Event, eventsBuffer)
	client := &Client{
		events:   events,
		registry: maps.NewShardedConcurrentMultiMap(),
	}

	go func(client *Client) {
		for event := range client.events {
			handlers, ok := client.registry.Get(event.Topic)
			if !ok {
				continue
			}

			for _, iHandler := range handlers {
				handler, ok := iHandler.(EventHandler)
				if !ok {
					continue
				}

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
	registry *maps.ShardedConcurrentMultiMap
}

// Publish publishes an event
func (c *Client) Publish(topic string, value interface{}) {
	event := &Event{Topic: topic, Value: value}

	// Attempt to publish an event to the channel, if the channel's buffer was full, discard
	select {
	case c.events <- event:
	default:
	}
}

// Subscribe subscribes to a topic
func (c *Client) Subscribe(topic string, handler EventHandler) {
	c.registry.Append(topic, handler)
}

// Close closes eventbus
func (c *Client) Close() {
	close(c.events)
	c.registry.Clear()
}
