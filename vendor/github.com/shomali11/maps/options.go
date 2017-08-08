package maps

const (
	defaultShards = 16
)

// ShardOption options for a map
type ShardOption func(*Config)

func WithNumberOfShards(shards uint32) ShardOption {
	return func(config *Config) {
		if shards < 1 {
			config.shards = defaultShards
		} else {
			config.shards = shards
		}
	}
}

// Config concurrent map
type Config struct {
	shards uint32
}

func getNumberOfShards(options ...ShardOption) uint32 {
	config := &Config{
		shards: defaultShards,
	}

	for _, option := range options {
		option(config)
	}

	return config.shards
}
