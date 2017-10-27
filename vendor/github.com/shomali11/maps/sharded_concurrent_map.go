package maps

import (
	"github.com/shomali11/util/xhashes"
)

// NewShardedConcurrentMap creates a new sharded concurrent map
func NewShardedConcurrentMap(options ...ShardOption) *ShardedConcurrentMap {
	shardedConcurrentMap := &ShardedConcurrentMap{
		shards: getNumberOfShards(options...),
	}

	concurrentMaps := make([]*ConcurrentMap, shardedConcurrentMap.shards)
	for i := uint32(0); i < shardedConcurrentMap.shards; i++ {
		concurrentMaps[i] = NewConcurrentMap()
	}

	shardedConcurrentMap.concurrentMaps = concurrentMaps
	return shardedConcurrentMap
}

// ShardedConcurrentMap concurrent map
type ShardedConcurrentMap struct {
	shards         uint32
	concurrentMaps []*ConcurrentMap
}

// Set concurrent set to map
func (c *ShardedConcurrentMap) Set(key string, value interface{}) {
	shard := c.getShard(key)
	concurrentMap := c.concurrentMaps[shard]
	concurrentMap.Set(key, value)
}

// Get concurrent get from map
func (c *ShardedConcurrentMap) Get(key string) (interface{}, bool) {
	shard := c.getShard(key)
	concurrentMap := c.concurrentMaps[shard]
	return concurrentMap.Get(key)
}

// Remove concurrent remove from map
func (c *ShardedConcurrentMap) Remove(key string) {
	shard := c.getShard(key)
	concurrentMap := c.concurrentMaps[shard]
	concurrentMap.Remove(key)
}

// ContainsKey concurrent contains key in map
func (c *ShardedConcurrentMap) ContainsKey(key string) bool {
	_, ok := c.Get(key)
	return ok
}

// ContainsEntry concurrent contains entry in map
func (c *ShardedConcurrentMap) ContainsEntry(key string, value interface{}) bool {
	shard := c.getShard(key)
	concurrentMap := c.concurrentMaps[shard]
	return concurrentMap.ContainsEntry(key, value)
}

// Size concurrent size of map
func (c *ShardedConcurrentMap) Size() int {
	sum := 0
	for _, concurrentMap := range c.concurrentMaps {
		sum += concurrentMap.Size()
	}
	return sum
}

// IsEmpty concurrent check of map's emptiness
func (c *ShardedConcurrentMap) IsEmpty() bool {
	for _, concurrentMap := range c.concurrentMaps {
		if !concurrentMap.IsEmpty() {
			return false
		}
	}
	return true
}

// Keys concurrent retrieval of keys from map
func (c *ShardedConcurrentMap) Keys() []string {
	keys := []string{}
	for _, concurrentMap := range c.concurrentMaps {
		keys = append(keys, concurrentMap.Keys()...)
	}
	return keys
}

// Remove concurrent remove from map
func (c *ShardedConcurrentMap) Clear() {
	for _, concurrentMap := range c.concurrentMaps {
		concurrentMap.Clear()
	}
}

func (c *ShardedConcurrentMap) getShard(key string) uint32 {
	return xhashes.FNV32(key) % uint32(c.shards)
}
