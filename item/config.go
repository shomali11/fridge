package item

import (
	"time"
)

// ConfigOption a config option
type ConfigOption func(*Config)

// NewConfig creates a new item config
func NewConfig(key string, options ...ConfigOption) *Config {
	config := &Config{Key: key}
	for _, option := range options {
		option(config)
	}
	return config
}

// WithRestock option to restock
func WithRestock(restock func() (string, error)) ConfigOption {
	return func(config *Config) {
		config.Restock = restock
	}
}

// WithDurations modifies configuration's default best by and use by durations
func WithDurations(bestBy time.Duration, useBy time.Duration) ConfigOption {
	return func(config *Config) {
		if bestBy <= useBy {
			config.BestBy = bestBy
			config.UseBy = useBy
		}
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
