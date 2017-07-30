package item

import (
	"errors"
	"time"
)

const (
	invalidDurations = "Invalid BestBy and UseBy durations"
)

// Config contains item configuration
type Config struct {
	Key            string
	StockTimestamp time.Time
	BestBy         time.Duration
	UseBy          time.Duration
}

// NewConfig creates a new item config
func NewConfig(key string, bestBy time.Duration, useBy time.Duration) (*Config, error) {
	if bestBy > useBy {
		return nil, errors.New(invalidDurations)
	}

	config := &Config{
		Key:    key,
		BestBy: bestBy,
		UseBy:  useBy,
	}
	return config, nil
}
