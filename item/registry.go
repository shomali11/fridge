package item

import (
	"errors"
	"fmt"
	"sync"
)

const (
	unregisteredItemFormat = "Unregistered item '%s'"
)

// NewRegistry creates a new item registry
func NewRegistry() *Registry {
	return &Registry{registry: make(map[string]*Config)}
}

// Registry contains registered items
type Registry struct {
	registry map[string]*Config
	mutex    sync.RWMutex
}

// Get retrieves an item configuration by key
func (r *Registry) Get(key string) (*Config, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	config, found := r.registry[key]
	if !found {
		return nil, errors.New(fmt.Sprintf(unregisteredItemFormat, key))
	}
	return config, nil
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
