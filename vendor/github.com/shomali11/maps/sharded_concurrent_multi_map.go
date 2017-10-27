package maps

import (
	"github.com/shomali11/util/xhashes"
)

// NewShardedConcurrentMultiMap creates a new sharded concurrent map
func NewShardedConcurrentMultiMap(options ...ShardOption) *ShardedConcurrentMultiMap {
	shardedConcurrentMultiMap := &ShardedConcurrentMultiMap{
		shards: getNumberOfShards(options...),
	}

	concurrentMaps := make([]*ConcurrentMultiMap, shardedConcurrentMultiMap.shards)
	for i := uint32(0); i < shardedConcurrentMultiMap.shards; i++ {
		concurrentMaps[i] = NewConcurrentMultiMap()
	}

	shardedConcurrentMultiMap.concurrentMaps = concurrentMaps
	return shardedConcurrentMultiMap
}

// ShardedConcurrentMultiMap concurrent map
type ShardedConcurrentMultiMap struct {
	shards         uint32
	concurrentMaps []*ConcurrentMultiMap
}

// Set concurrent set to map
func (c *ShardedConcurrentMultiMap) Set(key string, values []interface{}) {
	shard := c.getShard(key)
	concurrentMap := c.concurrentMaps[shard]
	concurrentMap.Set(key, values)
}

// Append concurrent append to map
func (c *ShardedConcurrentMultiMap) Append(key string, value interface{}) {
	shard := c.getShard(key)
	concurrentMap := c.concurrentMaps[shard]
	concurrentMap.Append(key, value)
}

// Get concurrent get from map
func (c *ShardedConcurrentMultiMap) Get(key string) ([]interface{}, bool) {
	shard := c.getShard(key)
	concurrentMap := c.concurrentMaps[shard]
	return concurrentMap.Get(key)
}

// Remove concurrent remove from map
func (c *ShardedConcurrentMultiMap) Remove(key string) {
	shard := c.getShard(key)
	concurrentMap := c.concurrentMaps[shard]
	concurrentMap.Remove(key)
}

// ContainsKey concurrent contains key in map
func (c *ShardedConcurrentMultiMap) ContainsKey(key string) bool {
	_, ok := c.Get(key)
	return ok
}

// ContainsEntry concurrent contains entry in map
func (c *ShardedConcurrentMultiMap) ContainsEntry(key string, value interface{}) bool {
	shard := c.getShard(key)
	concurrentMap := c.concurrentMaps[shard]
	return concurrentMap.ContainsEntry(key, value)
}

// Size concurrent size of map
func (c *ShardedConcurrentMultiMap) Size() int {
	sum := 0
	for _, concurrentMap := range c.concurrentMaps {
		sum += concurrentMap.Size()
	}
	return sum
}

// IsEmpty concurrent check of map's emptiness
func (c *ShardedConcurrentMultiMap) IsEmpty() bool {
	for _, concurrentMap := range c.concurrentMaps {
		if !concurrentMap.IsEmpty() {
			return false
		}
	}
	return true
}

// Keys concurrent retrieval of keys from map
func (c *ShardedConcurrentMultiMap) Keys() []string {
	keys := []string{}
	for _, concurrentMap := range c.concurrentMaps {
		keys = append(keys, concurrentMap.Keys()...)
	}
	return keys
}

// Remove concurrent remove from map
func (c *ShardedConcurrentMultiMap) Clear() {
	for _, concurrentMap := range c.concurrentMaps {
		concurrentMap.Clear()
	}
}

func (c *ShardedConcurrentMultiMap) getShard(key string) uint32 {
	return xhashes.FNV32(key) % uint32(c.shards)
}
