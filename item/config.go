package item

import (
	"errors"
	"time"
)

const (
	invalidDurations = "Invalid BestBy and UseBy durations"
)

// ConfigOption a config option
type ConfigOption func(*Config)

// NewConfig creates a new item config
func NewConfig(key string, bestBy time.Duration, useBy time.Duration, options ...ConfigOption) (*Config, error) {
	if bestBy > useBy {
		return nil, errors.New(invalidDurations)
	}

	config := &Config{
		Key:    key,
		BestBy: bestBy,
		UseBy:  useBy,
	}

	for _, option := range options {
		option(config)
	}
	return config, nil
}

// WithRestock option to restock
func WithRestock(restock func() (string, error)) ConfigOption {
	return func(config *Config) {
		config.Restock = restock
	}
}

// Config contains item configuration
type Config struct {
	Key     string
	BestBy  time.Duration
	UseBy   time.Duration
	Restock func() (string, error)
}

// GetUseByInSeconds returns Use By duration in seconds
func (c *Config) GetUseByInSeconds() int64 {
	return int64(c.UseBy.Seconds())
}
