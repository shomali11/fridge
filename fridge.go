package fridge

import (
	"github.com/shomali11/fridge/item"
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
	empty = ""
)

// NewClient returns a client using an xredis client
func NewClient(options ...ConfigOption) *Client {
	config := NewConfig(options...)

	client := &Client{
		itemRegistry: item.NewRegistry(config.defaultBestBy, config.defaultUseBy),
		itemDao:      item.NewDao(config.xredisClient),
		eventBus:     NewEventBus(),
	}
	return client
}

// Client fridge client
type Client struct {
	itemRegistry *item.Registry
	itemDao      *item.Dao
	eventBus     *EventBus
}

// Register an item
func (c *Client) Register(key string, options ...item.ConfigOption) {
	itemConfig := item.NewConfig(key, options...)
	c.itemRegistry.Set(itemConfig)
}

// Deregister an item
func (c *Client) Deregister(key string) {
	c.itemRegistry.Remove(key)
}

// Put an item
func (c *Client) Put(key string, value string) error {
	itemConfig := c.itemRegistry.Get(key)
	err := c.itemDao.Set(key, value, itemConfig.GetUseByInSeconds())
	if err != nil {
		return err
	}
	return nil
}

// Get an item
func (c *Client) Get(key string) (string, bool, error) {
	itemConfig := c.itemRegistry.Get(key)
	cachedValue, found, stockTimestamp, err := c.itemDao.Get(key)
	if err != nil {
		return empty, false, err
	}

	if !found {
		if stockTimestamp.IsZero() {
			go c.publish(key, NotFound)
		} else {
			go c.publish(key, Expired)
		}
		return c.restock(itemConfig)
	}

	now := time.Now().UTC()
	if now.Before(stockTimestamp.Add(itemConfig.BestBy)) {
		go c.publish(key, Fresh)
		return cachedValue, true, nil
	}

	if now.Before(stockTimestamp.Add(itemConfig.UseBy)) {
		go c.publish(key, Cold)
		go c.restockAndCompare(cachedValue, itemConfig)
		return cachedValue, true, nil
	}

	go c.publish(key, Expired)
	return c.restockAndCompare(cachedValue, itemConfig)
}

// Remove an item
func (c *Client) Remove(key string) error {
	return c.itemDao.Remove(key)
}

// Ping pings redis
func (c *Client) Ping() error {
	return c.itemDao.Ping()
}

// Close closes resources
func (c *Client) Close() error {
	return c.itemDao.Close()
}

// HandleEvent overrides the default handleEvent callback
func (c *Client) HandleEvent(handleEvent func(event *Event)) {
	c.eventBus.HandleEvent(handleEvent)
}

func (c *Client) publish(key string, eventType string) {
	c.eventBus.Publish(key, eventType)
}

func (c *Client) restockAndCompare(cachedValue string, itemConfig *item.Config) (string, bool, error) {
	newValue, found, err := c.restock(itemConfig)
	if err != nil || !found {
		return empty, found, err
	}

	if newValue == cachedValue {
		go c.publish(itemConfig.Key, Unchanged)
	}
	return newValue, true, nil
}

func (c *Client) restock(itemConfig *item.Config) (string, bool, error) {
	if itemConfig.Restock == nil {
		go c.publish(itemConfig.Key, OutOfStock)
		return empty, false, nil
	}

	result, err := itemConfig.Restock()
	if err != nil {
		return empty, false, err
	}

	go c.publish(itemConfig.Key, Refresh)

	err = c.Put(itemConfig.Key, result)
	if err != nil {
		return empty, false, err
	}
	return result, true, nil
}
