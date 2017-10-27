package maps

import "sync"

// NewConcurrentMap creates a new concurrent map
func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		internalMap: make(map[string]interface{}),
		lock:        &sync.RWMutex{},
	}
}

// ConcurrentMap concurrent map
type ConcurrentMap struct {
	internalMap map[string]interface{}
	lock        *sync.RWMutex
}

// Set concurrent set to map
func (c *ConcurrentMap) Set(key string, value interface{}) {
	c.lock.Lock()
	c.internalMap[key] = value
	c.lock.Unlock()
}

// Get concurrent get from map
func (c *ConcurrentMap) Get(key string) (interface{}, bool) {
	c.lock.RLock()
	value, ok := c.internalMap[key]
	c.lock.RUnlock()
	return value, ok
}

// Remove concurrent remove from map
func (c *ConcurrentMap) Remove(key string) {
	c.lock.Lock()
	delete(c.internalMap, key)
	c.lock.Unlock()
}

// ContainsKey concurrent contains key in map
func (c *ConcurrentMap) ContainsKey(key string) bool {
	_, ok := c.Get(key)
	return ok
}

// ContainsEntry concurrent contains entry in map
func (c *ConcurrentMap) ContainsEntry(key string, value interface{}) bool {
	existingValue, ok := c.Get(key)
	return ok && existingValue == value
}

// Size concurrent size of map
func (c *ConcurrentMap) Size() int {
	c.lock.RLock()
	size := len(c.internalMap)
	c.lock.RUnlock()
	return size
}

// IsEmpty concurrent check of map's emptiness
func (c *ConcurrentMap) IsEmpty() bool {
	return c.Size() == 0
}

// Keys concurrent retrieval of keys from map
func (c *ConcurrentMap) Keys() []string {
	c.lock.RLock()
	keys := make([]string, len(c.internalMap))
	i := 0
	for key := range c.internalMap {
		keys[i] = key
		i++
	}
	c.lock.RUnlock()
	return keys
}

// Clear concurrent map
func (c *ConcurrentMap) Clear() {
	c.lock.Lock()
	c.internalMap = make(map[string]interface{})
	c.lock.Unlock()
}
