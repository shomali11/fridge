package item

import (
	"time"
)

// ConfigOption a config option
type ConfigOption func(*Config)

// WithDurations modifies configuration's default best by and use by durations
func WithDurations(bestBy time.Duration, useBy time.Duration) ConfigOption {
	return func(config *Config) {
		config.BestBy = bestBy
		config.UseBy = useBy
	}
}

// Config contains item configuration
type Config struct {
	Timestamp time.Time
	BestBy    time.Duration
	UseBy     time.Duration
}

// GetUseByInSeconds returns Use By duration in seconds
func (c *Config) GetUseByInSeconds() int64 {
	return int64(c.UseBy.Seconds())
}

// QueryConfigOption a config option
type QueryConfigOption func(*QueryConfig)

// WithRestock sets configuration's restocking mechanism
func WithRestock(restock func() (string, error)) QueryConfigOption {
	return func(queryConfig *QueryConfig) {
		queryConfig.Restock = restock
	}
}

// QueryConfig contains item configuration
type QueryConfig struct {
	Restock func() (string, error)
}
