package fridge

import (
	"errors"
	"github.com/shomali11/eventbus"
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

	// Refresh is when an item was restocked with a fresher one
	Refresh = "REFRESH"

	// OutOfStock is when an item needs restocking, but no restocking function was provided
	OutOfStock = "OUT_OF_STOCK"

	// Unchanged is when the restocked item is not different from the version in the cache
	Unchanged = "UNCHANGED"
)

const (
	empty                 = ""
	eventsTopic           = "fridge_events"
	invalidDurationsError = "Invalid 'best by' and 'use by' durations"
)

// NewClient returns a client
func NewClient(redisClient *RedisClient, options ...DefaultsOption) *Client {
	client := &Client{
		defaults: newDefaults(options...),
		dao:      newDao(redisClient),
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

// Fridge event
type Event struct {
	Key  string
	Type string
}

// Client fridge client
type Client struct {
	defaults    *Defaults
	dao         *Dao
	bus         *eventbus.Client
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

	err = c.dao.Set(key, value, storageDetails.GetUseByInSeconds())
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
		return c.restockAndCompare(key, cachedValue, restock)
	}

	now := time.Now().UTC()
	if now.Before(storageDetails.Timestamp.Add(storageDetails.BestBy)) {
		c.publish(key, Fresh)
		return cachedValue, true, nil
	}

	if now.Before(storageDetails.Timestamp.Add(storageDetails.UseBy)) {
		c.publish(key, Cold)
		go c.restockAndCompare(key, cachedValue, restock)
		return cachedValue, true, nil
	}

	c.publish(key, Expired)
	return c.restockAndCompare(key, cachedValue, restock)
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

func (c *Client) restockAndCompare(key string, cachedValue string, callback func() (string, error)) (string, bool, error) {
	newValue, found, err := c.restock(key, callback)
	if err != nil || !found {
		return empty, found, err
	}

	if newValue == cachedValue {
		go c.publish(key, Unchanged)
	}
	return newValue, true, nil
}

func (c *Client) restock(key string, callback func() (string, error)) (string, bool, error) {
	if callback == nil {
		go c.publish(key, OutOfStock)
		return empty, false, nil
	}

	result, err := callback()
	if err != nil {
		return empty, false, err
	}

	go c.publish(key, Refresh)

	bestBy, useBy := c.getDurations(key)
	err = c.Put(key, result, WithDurations(bestBy, useBy))
	if err != nil {
		return empty, false, err
	}
	return result, true, nil
}

func (c *Client) getDurations(key string) (time.Duration, time.Duration) {
	bestBy := c.defaults.BestBy
	useBy := c.defaults.UseBy

	itemConfig, found, err := c.dao.GetStorageDetails(key)
	if found && err == nil {
		bestBy = itemConfig.BestBy
		useBy = itemConfig.UseBy
	}
	return bestBy, useBy
}
