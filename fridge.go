package fridge

import (
	"errors"
	"github.com/shomali11/fridge/item"
	"github.com/shomali11/xredis"
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
	invalidDurationsError = "Invalid 'best by' and 'use by' durations"
)

// NewClient returns a client using an xredis client
func NewClient(xredisClient *xredis.Client, options ...ConfigOption) *Client {
	client := &Client{
		config:   newConfig(options...),
		itemDao:  item.NewDao(xredisClient),
		eventBus: newEventBus(),
	}
	return client
}

// Client fridge client
type Client struct {
	config   *Config
	itemDao  *item.Dao
	eventBus *EventBus
}

// Put an item
func (c *Client) Put(key string, value string, options ...item.ConfigOption) error {
	itemConfig := newItemConfig(c.config, options...)
	if itemConfig.BestBy > itemConfig.UseBy {
		return errors.New(invalidDurationsError)
	}

	err := c.itemDao.SetConfig(key, itemConfig)
	if err != nil {
		return err
	}

	err = c.itemDao.Set(key, value, itemConfig.GetUseByInSeconds())
	if err != nil {
		return err
	}
	return nil
}

// Get an item
func (c *Client) Get(key string, options ...item.QueryConfigOption) (string, bool, error) {
	queryConfig := newQueryConfig(options...)
	restock := queryConfig.Restock

	itemConfig, found, err := c.itemDao.GetConfig(key)
	if err != nil {
		return empty, false, err
	}

	if !found {
		go c.publish(key, NotFound)
		return empty, false, err
	}

	cachedValue, found, err := c.itemDao.Get(key)
	if err != nil {
		return empty, false, err
	}

	if !found {
		go c.publish(key, Expired)
		return c.restockAndCompare(key, cachedValue, restock)
	}

	now := time.Now().UTC()
	if now.Before(itemConfig.Timestamp.Add(itemConfig.BestBy)) {
		go c.publish(key, Fresh)
		return cachedValue, true, nil
	}

	if now.Before(itemConfig.Timestamp.Add(itemConfig.UseBy)) {
		go c.publish(key, Cold)
		go c.restockAndCompare(key, cachedValue, restock)
		return cachedValue, true, nil
	}

	go c.publish(key, Expired)
	return c.restockAndCompare(key, cachedValue, restock)
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
	err = c.Put(key, result, item.WithDurations(bestBy, useBy))
	if err != nil {
		return empty, false, err
	}
	return result, true, nil
}

func (c *Client) getDurations(key string) (time.Duration, time.Duration) {
	bestBy := c.config.defaultBestBy
	useBy := c.config.defaultUseBy

	itemConfig, found, err := c.itemDao.GetConfig(key)
	if found && err == nil {
		bestBy = itemConfig.BestBy
		useBy = itemConfig.UseBy
	}
	return bestBy, useBy
}

func newItemConfig(config *Config, options ...item.ConfigOption) *item.Config {
	itemConfig := &item.Config{BestBy: config.defaultBestBy, UseBy: config.defaultUseBy}
	for _, option := range options {
		option(itemConfig)
	}
	return itemConfig
}

func newQueryConfig(options ...item.QueryConfigOption) *item.QueryConfig {
	queryConfig := &item.QueryConfig{}
	for _, option := range options {
		option(queryConfig)
	}
	return queryConfig
}
