package item

import (
	"sync"
	"time"
)

// NewRegistry creates a new item registry
func NewRegistry(defaultBestBy time.Duration, defaultUseBy time.Duration) *Registry {
	return &Registry{
		registry:      make(map[string]*Config),
		defaultBestBy: defaultBestBy,
		defaultUseBy:  defaultUseBy,
	}
}

// Registry contains registered items
type Registry struct {
	registry      map[string]*Config
	defaultBestBy time.Duration
	defaultUseBy  time.Duration
	mutex         sync.RWMutex
}

// Get retrieves an item configuration by key
func (r *Registry) Get(key string) *Config {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	config, found := r.registry[key]
	if !found {
		return &Config{Key: key, BestBy: r.defaultBestBy, UseBy: r.defaultUseBy}
	}
	return config
}

// Set registers an item configuration
func (r *Registry) Set(config *Config) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.registry[config.Key] = config
}

// Remove deregisters an item configuration
func (r *Registry) Remove(key string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.registry, key)
}
