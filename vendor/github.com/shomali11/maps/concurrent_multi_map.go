package maps

import "sync"

// NewConcurrentMultiMap creates a new concurrent multi map
func NewConcurrentMultiMap() *ConcurrentMultiMap {
	return &ConcurrentMultiMap{
		internalMap: make(map[string][]interface{}),
		lock:        &sync.RWMutex{},
	}
}

// ConcurrentMultiMap concurrent map
type ConcurrentMultiMap struct {
	internalMap map[string][]interface{}
	lock        *sync.RWMutex
}

// Set concurrent set to map
func (c *ConcurrentMultiMap) Set(key string, values []interface{}) {
	c.lock.Lock()
	c.internalMap[key] = values
	c.lock.Unlock()
}

// Append concurrent append to map
func (c *ConcurrentMultiMap) Append(key string, value interface{}) {
	c.lock.Lock()
	values, ok := c.internalMap[key]
	if !ok {
		values = []interface{}{}
	}
	values = append(values, value)
	c.internalMap[key] = values
	c.lock.Unlock()
}

// Get concurrent get from map
func (c *ConcurrentMultiMap) Get(key string) ([]interface{}, bool) {
	c.lock.RLock()
	value, ok := c.internalMap[key]
	c.lock.RUnlock()
	return value, ok
}

// Remove concurrent remove from map
func (c *ConcurrentMultiMap) Remove(key string) {
	c.lock.Lock()
	delete(c.internalMap, key)
	c.lock.Unlock()
}

// ContainsKey concurrent contains key in map
func (c *ConcurrentMultiMap) ContainsKey(key string) bool {
	_, ok := c.Get(key)
	return ok
}

// ContainsEntry concurrent contains entry in map
func (c *ConcurrentMultiMap) ContainsEntry(key string, value interface{}) bool {
	existingValues, ok := c.Get(key)
	if !ok {
		return false
	}

	for _, existingValue := range existingValues {
		if existingValue == value {
			return true
		}
	}
	return false
}

// Size concurrent size of map
func (c *ConcurrentMultiMap) Size() int {
	c.lock.RLock()
	size := len(c.internalMap)
	c.lock.RUnlock()
	return size
}

// IsEmpty concurrent check of map's emptiness
func (c *ConcurrentMultiMap) IsEmpty() bool {
	return c.Size() == 0
}

// Keys concurrent retrieval of keys from map
func (c *ConcurrentMultiMap) Keys() []string {
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
func (c *ConcurrentMultiMap) Clear() {
	c.lock.Lock()
	c.internalMap = make(map[string][]interface{})
	c.lock.Unlock()
}
