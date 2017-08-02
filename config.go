package fridge

import (
	"github.com/shomali11/xredis"
	"time"
)

const (
	defaultBestBy = time.Hour
	defaultUseBy  = 24 * time.Hour
)

// ConfigOption option to configure the fridge client
type ConfigOption func(*Config)

// NewConfig returns a client configuration
func NewConfig(options ...ConfigOption) *Config {
	config := &Config{
		defaultBestBy: defaultBestBy,
		defaultUseBy:  defaultUseBy,
	}

	for _, option := range options {
		option(config)
	}

	if config.xredisClient == nil {
		config.xredisClient = xredis.DefaultClient()
	}
	return config
}

// WithRedisClient modifies configuration using the xredisClient
func WithRedisClient(xredisClient *xredis.Client) ConfigOption {
	return func(config *Config) {
		config.xredisClient = xredisClient
	}
}

// WithDefaultDurations modifies configuration's default best by and use by durations
func WithDefaultDurations(bestBy time.Duration, useBy time.Duration) ConfigOption {
	return func(config *Config) {
		if bestBy <= useBy {
			config.defaultBestBy = bestBy
			config.defaultUseBy = useBy
		} else {
			config.defaultBestBy = defaultBestBy
			config.defaultUseBy = defaultUseBy
		}
	}
}

// Config configuration for the fridge client
type Config struct {
	xredisClient  *xredis.Client
	defaultBestBy time.Duration
	defaultUseBy  time.Duration
}
