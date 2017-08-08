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

// Contains concurrent contains in map
func (c *ConcurrentMap) Contains(key string) bool {
	_, ok := c.Get(key)
	return ok
}

// Size concurrent size of map
func (c *ConcurrentMap) Size() int {
	c.lock.RLock()
	size := len(c.internalMap)
	c.lock.RUnlock()
	return size
}

// Clear concurrent map
func (c *ConcurrentMap) Clear() {
	c.lock.Lock()
	c.internalMap = make(map[string]interface{})
	c.lock.Unlock()
}
