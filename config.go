package fridge

import (
	"time"
)

const (
	defaultBestBy = time.Hour
	defaultUseBy  = 24 * time.Hour
)

// ConfigOption option to configure the fridge client
type ConfigOption func(*Config)

// WithDefaultDurations modifies configuration's default best by and use by durations
func WithDefaultDurations(bestBy time.Duration, useBy time.Duration) ConfigOption {
	return func(config *Config) {
		config.defaultBestBy = bestBy
		config.defaultUseBy = useBy
	}
}

// Config configuration for the fridge client
type Config struct {
	defaultBestBy time.Duration
	defaultUseBy  time.Duration
}

func newConfig(options ...ConfigOption) *Config {
	config := &Config{
		defaultBestBy: defaultBestBy,
		defaultUseBy:  defaultUseBy,
	}

	for _, option := range options {
		option(config)
	}
	return config
}
