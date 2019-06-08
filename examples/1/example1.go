package main

import (
	"fmt"
	"github.com/shomali11/fridge"
	"time"
)

type SimpleCache struct {
	memory map[string]string
}

func (c *SimpleCache) Get(key string) (string, bool, error) {
	value, ok := c.memory[key]
	return value, ok, nil
}

func (c *SimpleCache) Set(key string, value string, timeout time.Duration) error {
	// We are not implementing the expiration to keep the example simple
	c.memory[key] = value
	return nil
}

func (c *SimpleCache) Remove(key string) error {
	delete(c.memory, key)
	return nil
}

func (c *SimpleCache) Ping() error {
	return nil
}

func (c *SimpleCache) Close() error {
	return nil
}

func main() {
	simpleCache := &SimpleCache{memory: make(map[string]string)}
	client := fridge.NewClient(simpleCache)

	fmt.Println(client.Ping())
}
