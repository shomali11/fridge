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

// Contains concurrent contains in map
func (c *ConcurrentMultiMap) Contains(key string) bool {
	_, ok := c.Get(key)
	return ok
}

// Size concurrent size of map
func (c *ConcurrentMultiMap) Size() int {
	c.lock.RLock()
	size := len(c.internalMap)
	c.lock.RUnlock()
	return size
}

// Clear concurrent map
func (c *ConcurrentMultiMap) Clear() {
	c.lock.Lock()
	c.internalMap = make(map[string][]interface{})
	c.lock.Unlock()
}
