package fridge

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/shomali11/fridge/item"
	"github.com/shomali11/util/conversions"
	"github.com/shomali11/xredis"
	"time"
)

const (
	unregisteredItemFormat = "Unregistered item '%s'"
	itemConfigKeyFormat    = "%s.config"
	empty                  = ""
)

// DefaultClient returns a client with default options
func DefaultClient() *Client {
	client := xredis.DefaultClient()
	return &Client{xredisClient: client}
}

// SetupClient returns a client with provided options
func SetupClient(options *xredis.Options) *Client {
	client := xredis.SetupClient(options)
	return &Client{xredisClient: client}
}

// NewClient returns a client using provided redis.Pool
func NewClient(pool *redis.Pool) *Client {
	client := xredis.NewClient(pool)
	return &Client{xredisClient: client}
}

// Client fridge client
type Client struct {
	xredisClient *xredis.Client
}

// Register an item
func (c *Client) Register(key string, bestBy time.Duration, useBy time.Duration) error {
	itemConfig, err := item.NewConfig(key, bestBy, useBy)
	if err != nil {
		return err
	}

	itemConfigKey := fmt.Sprintf(itemConfigKeyFormat, key)
	itemConfigString, err := conversions.Stringify(itemConfig)
	if err != nil {
		return err
	}

	_, err = c.xredisClient.Set(itemConfigKey, itemConfigString)
	return err
}

// Put an item
func (c *Client) Put(key string, value string) error {
	itemConfig, err := c.retrieveItemConfig(key)
	if err != nil {
		return err
	}

	useByInSeconds := int64(itemConfig.UseBy.Seconds())
	_, err = c.xredisClient.SetEx(key, value, useByInSeconds)
	if err != nil {
		return err
	}

	itemConfig.StockTimestamp = time.Now().UTC()
	err = c.saveItemInfo(itemConfig)
	if err != nil {
		return err
	}
	return nil
}

// Get an item
func (c *Client) Get(key string, restock func() (string, error)) (string, bool, error) {
	itemConfig, err := c.retrieveItemConfig(key)
	if err != nil {
		return empty, false, err
	}

	value, ok, err := c.xredisClient.Get(key)
	if err != nil {
		return empty, false, err
	}

	if !ok {
		return c.callRestock(itemConfig, restock)
	}

	now := time.Now().UTC()
	if now.Before(itemConfig.StockTimestamp.Add(itemConfig.BestBy)) {
		return value, true, nil
	}

	if now.Before(itemConfig.StockTimestamp.Add(itemConfig.UseBy)) {
		go c.callRestock(itemConfig, restock)
		return value, true, nil
	}
	return c.callRestock(itemConfig, restock)
}

// Remove an item
func (c *Client) Remove(key string) error {
	_, err := c.xredisClient.Del(key)
	return err
}

// Deregister an item
func (c *Client) Deregister(key string) error {
	itemConfigKey := fmt.Sprintf(itemConfigKeyFormat, key)
	_, err := c.xredisClient.Del(itemConfigKey)
	return err
}

// Ping pings redis
func (c *Client) Ping() error {
	_, err := c.xredisClient.Ping()
	return err
}

// Close closes resources
func (c *Client) Close() error {
	return c.xredisClient.Close()
}

func (c *Client) callRestock(itemConfig *item.Config, restock func() (string, error)) (string, bool, error) {
	if restock == nil {
		return empty, false, nil
	}

	result, err := restock()
	if err != nil {
		return empty, false, err
	}

	err = c.Put(itemConfig.Key, result)
	if err != nil {
		return empty, false, err
	}
	return result, true, nil
}

func (c *Client) saveItemInfo(itemConfig *item.Config) error {
	itemConfigString, err := conversions.Stringify(itemConfig)
	if err != nil {
		return err
	}

	itemConfigKey := fmt.Sprintf(itemConfigKeyFormat, itemConfig.Key)
	_, err = c.xredisClient.Set(itemConfigKey, itemConfigString)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) retrieveItemConfig(key string) (*item.Config, error) {
	itemConfigKey := fmt.Sprintf(itemConfigKeyFormat, key)
	itemConfigString, ok, err := c.xredisClient.Get(itemConfigKey)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New(fmt.Sprintf(unregisteredItemFormat, key))
	}

	itemConfig := item.Config{}
	err = conversions.Structify(itemConfigString, &itemConfig)
	if err != nil {
		return nil, err
	}
	return &itemConfig, nil
}