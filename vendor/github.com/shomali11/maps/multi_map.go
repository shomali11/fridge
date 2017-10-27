package maps

// NewMultiMap creates a new multi map
func NewMultiMap() *MultiMap {
	return &MultiMap{
		internalMap: make(map[string][]interface{}),
	}
}

// MultiMap map
type MultiMap struct {
	internalMap map[string][]interface{}
}

// Set set to map
func (c *MultiMap) Set(key string, values []interface{}) {
	c.internalMap[key] = values
}

// Append append to map
func (c *MultiMap) Append(key string, value interface{}) {
	values, ok := c.internalMap[key]
	if !ok {
		values = []interface{}{}
	}
	values = append(values, value)
	c.internalMap[key] = values
}

// Get get from map
func (c *MultiMap) Get(key string) ([]interface{}, bool) {
	value, ok := c.internalMap[key]
	return value, ok
}

// Remove remove from map
func (c *MultiMap) Remove(key string) {
	delete(c.internalMap, key)
}

// ContainsKey contains key in map
func (c *MultiMap) ContainsKey(key string) bool {
	_, ok := c.Get(key)
	return ok
}

// ContainsEntry concurrent contains entry in map
func (c *MultiMap) ContainsEntry(key string, value interface{}) bool {
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

// Size size of map
func (c *MultiMap) Size() int {
	return len(c.internalMap)
}

// IsEmpty check of map's emptiness
func (c *MultiMap) IsEmpty() bool {
	return c.Size() == 0
}

// Keys retrieval of keys from map
func (c *MultiMap) Keys() []string {
	keys := make([]string, len(c.internalMap))
	i := 0
	for key := range c.internalMap {
		keys[i] = key
		i++
	}
	return keys
}

// Clear map
func (c *MultiMap) Clear() {
	c.internalMap = make(map[string][]interface{})
}
