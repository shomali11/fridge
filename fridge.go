package fridge

import (
	"errors"
	"github.com/shomali11/eventbus"
	"github.com/shomali11/parallelizer"
	"time"
)

const (
	// Fresh is when an item has not passed its "Best By" duration
	Fresh = "FRESH"

	// Cold is when an item has passed its "Best By" duration but not its "Use By" one
	Cold = "COLD"

	// Expired is when an item has passed its "Use By" duration
	Expired = "EXPIRED"

	// NotFound is when an item was not found due to being removed or was never stored before
	NotFound = "NOT_FOUND"

	// Restock is when an item was restocked with a fresher one
	Restock = "RESTOCK"

	// OutOfStock is when an item needs restocking, but no restocking function was provided
	OutOfStock = "OUT_OF_STOCK"

	// Unchanged is when the restocked item is not different from the version in the cache
	Unchanged = "UNCHANGED"
)

const (
	empty                 = ""
	eventsTopic           = "fridge_events"
	invalidDurationsError = "invalid 'best by' and 'use by' durations"
)

// NewClient returns a client
func NewClient(cache Cache, options ...DefaultsOption) *Client {
	client := &Client{
		defaults: newDefaults(options...),
		dao:      newDao(cache),
		group:    parallelizer.NewGroup(),
	}

	bus := eventbus.NewClient()
	bus.Subscribe(eventsTopic, func(value interface{}) {
		event, ok := value.(*Event)
		if !ok {
			return
		}

		if client.handleEvent == nil {
			return
		}
		client.handleEvent(event)
	})

	client.bus = bus
	return client
}

// Cache is a Fridge cache interface
type Cache interface {
	// Get a value by key
	Get(key string) (string, bool, error)

	// Set a key value pair
	Set(key string, value string, timeout time.Duration) error

	// Remove a key
	Remove(key string) error

	// Ping to test connectivity
	Ping() error

	// Close to close resources
	Close() error
}

// Event is a Fridge event
type Event struct {
	Key  string
	Type string
}

// Client fridge client
type Client struct {
	defaults    *Defaults
	dao         *Dao
	bus         *eventbus.Client
	group       *parallelizer.Group
	handleEvent func(event *Event)
}

// Put an item
func (c *Client) Put(key string, value string, options ...StorageOption) error {
	storageDetails := newStorageDetails(c.defaults, options...)
	if storageDetails.BestBy > storageDetails.UseBy {
		return errors.New(invalidDurationsError)
	}

	err := c.dao.SetStorageDetails(key, storageDetails)
	if err != nil {
		return err
	}

	err = c.dao.Set(key, value, storageDetails.UseBy)
	if err != nil {
		return err
	}
	return nil
}

// Get an item
func (c *Client) Get(key string, options ...RetrievalOption) (string, bool, error) {
	retrievalDetails := newRetrievalDetails(options...)
	restock := retrievalDetails.Restock

	storageDetails, found, err := c.dao.GetStorageDetails(key)
	if err != nil {
		return empty, false, err
	}

	if !found {
		c.publish(key, NotFound)
		return empty, false, err
	}

	cachedValue, found, err := c.dao.Get(key)
	if err != nil {
		return empty, false, err
	}

	if !found {
		c.publish(key, Expired)
		return c.restock(key, cachedValue, storageDetails, restock)
	}

	now := time.Now().UTC()
	if now.Before(storageDetails.Timestamp.Add(storageDetails.BestBy)) {
		c.publish(key, Fresh)
		return cachedValue, true, nil
	}

	if now.Before(storageDetails.Timestamp.Add(storageDetails.UseBy)) {
		c.publish(key, Cold)
		if !storageDetails.Restocking {
			c.group.Add(func() {
				c.restock(key, cachedValue, storageDetails, restock)
			})
		}
		return cachedValue, true, nil
	}

	c.publish(key, Expired)
	return c.restock(key, cachedValue, storageDetails, restock)
}

// Remove an item
func (c *Client) Remove(key string) error {
	return c.dao.Remove(key)
}

// Ping pings redis
func (c *Client) Ping() error {
	return c.dao.Ping()
}

// Close closes resources
func (c *Client) Close() error {
	c.bus.Close()
	return c.dao.Close()
}

// HandleEvent overrides the default handleEvent callback
func (c *Client) HandleEvent(handleEvent func(event *Event)) {
	c.handleEvent = handleEvent
}

func (c *Client) publish(key string, eventType string) {
	c.bus.Publish(eventsTopic, &Event{Key: key, Type: eventType})
}

func (c *Client) restock(key string, cachedValue string, storageDetails *StorageDetails, callback func() (string, error)) (string, bool, error) {
	if callback == nil {
		c.publish(key, OutOfStock)
		return empty, false, nil
	}

	storageDetails.Restocking = true
	err := c.dao.SetStorageDetails(key, storageDetails)
	if err != nil {
		return empty, false, err
	}

	freshValue, err := callback()
	if err != nil {
		storageDetails.Restocking = false
		c.dao.SetStorageDetails(key, storageDetails)
		return empty, false, err
	}

	c.publish(key, Restock)

	bestBy, useBy := storageDetails.BestBy, storageDetails.UseBy
	err = c.Put(key, freshValue, WithDurations(bestBy, useBy))
	if err != nil {
		return empty, false, err
	}

	if freshValue == cachedValue {
		c.publish(key, Unchanged)
	}
	return freshValue, true, nil
}
