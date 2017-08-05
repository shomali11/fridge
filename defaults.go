package fridge

import (
	"time"
)

const (
	defaultBestBy = time.Hour
	defaultUseBy  = 24 * time.Hour
)

// DefaultsOption an option for default values
type DefaultsOption func(*Defaults)

// WithDefaultDurations sets default best by and use by durations
func WithDefaultDurations(bestBy time.Duration, useBy time.Duration) DefaultsOption {
	return func(defaults *Defaults) {
		defaults.BestBy = bestBy
		defaults.UseBy = useBy
	}
}

// Defaults configuration for the fridge client
type Defaults struct {
	BestBy time.Duration
	UseBy  time.Duration
}

func newDefaults(options ...DefaultsOption) *Defaults {
	config := &Defaults{
		BestBy: defaultBestBy,
		UseBy:  defaultUseBy,
	}

	for _, option := range options {
		option(config)
	}
	return config
}
