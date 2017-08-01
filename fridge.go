package fridge

import (
	"github.com/garyburd/redigo/redis"
	"github.com/shomali11/fridge/item"
	"github.com/shomali11/xredis"
	"time"
)

const (
	// Fresh is when an item has not passed its "Best By" duration
	Fresh = "FRESH"

	// Cold is when an item has passed its "Best By" duration but not its "Use By" one
	Cold = "COLD"

	// Expird is when an item has passed its "Use By" duration
	Expired = "EXPIRED"

	// NotFound is when an item was never stored before
	NotFound = "NOT_FOUND"

	// Refresh is when an item was restocked with a fresher one
	Refresh = "REFRESH"

	// OutOfStock is when an item needs restocking, but no restocking function was provided
	OutOfStock = "OUT_OF_STOCK"
)

const (
	empty = ""
)

// DefaultClient returns a client with default options
func DefaultClient() *Client {
	client := xredis.DefaultClient()
	return newClient(client)
}

// SetupClient returns a client with provided options
func SetupClient(options *xredis.Options) *Client {
	client := xredis.SetupClient(options)
	return newClient(client)
}

// NewClient returns a client using provided redis.Pool
func NewClient(pool *redis.Pool) *Client {
	client := xredis.NewClient(pool)
	return newClient(client)
}

// Client fridge client
type Client struct {
	itemRegistry *item.Registry
	itemDao      *item.Dao
	eventBus     *EventBus
}

// Register an item
func (c *Client) Register(key string, bestBy time.Duration, useBy time.Duration, options ...item.ConfigOption) error {
	itemConfig, err := item.NewConfig(key, bestBy, useBy, options...)
	if err != nil {
		return err
	}

	c.itemRegistry.Set(itemConfig)
	return nil
}

// Deregister an item
func (c *Client) Deregister(key string) {
	c.itemRegistry.Remove(key)
}

// Put an item
func (c *Client) Put(key string, value string) error {
	itemConfig, err := c.itemRegistry.Get(key)
	if err != nil {
		return err
	}

	err = c.itemDao.Set(key, value, itemConfig.GetUseByInSeconds())
	if err != nil {
		return err
	}
	return c.itemDao.SetTimestamp(key, time.Now().UTC())
}

// Get an item
func (c *Client) Get(key string) (string, bool, error) {
	itemConfig, err := c.itemRegistry.Get(key)
	if err != nil {
		return empty, false, err
	}

	value, found, err := c.itemDao.Get(key)
	if err != nil {
		return empty, false, err
	}

	stockTimestamp, err := c.itemDao.GetTimestamp(key)
	if err != nil {
		return empty, false, err
	}

	if !found {
		if stockTimestamp.IsZero() {
			go c.publish(key, NotFound)
		} else {
			go c.publish(key, Expired)
		}
		return c.callRestock(itemConfig)
	}

	now := time.Now().UTC()
	if now.Before(stockTimestamp.Add(itemConfig.BestBy)) {
		go c.publish(key, Fresh)
		return value, true, nil
	}

	if now.Before(stockTimestamp.Add(itemConfig.UseBy)) {
		go c.publish(key, Cold)
		go c.callRestock(itemConfig)
		return value, true, nil
	}

	go c.publish(key, Expired)
	return c.callRestock(itemConfig)
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

func (c *Client) callRestock(itemConfig *item.Config) (string, bool, error) {
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

func newClient(xredisClient *xredis.Client) *Client {
	client := &Client{
		itemRegistry: item.NewRegistry(),
		itemDao:      item.NewDao(xredisClient),
		eventBus:     NewEventBus(),
	}
	return client
}
