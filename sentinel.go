package fridge

import (
	"crypto/tls"
	"github.com/shomali11/xredis"
	"time"
)

// SentinelOption an option for a sentinel option
type SentinelOption func(*SentinelSettings)

// WithSentinelAddresses sets sentinel addresses
func WithSentinelAddresses(addresses []string) SentinelOption {
	return func(redisSettings *SentinelSettings) {
		redisSettings.Addresses = addresses
	}
}

// WithSentinelMasterName sets sentinel master name
func WithSentinelMasterName(masterName string) SentinelOption {
	return func(redisSettings *SentinelSettings) {
		redisSettings.MasterName = masterName
	}
}

// WithRedisPassword sets redis password
func WithRedisPassword(password string) SentinelOption {
	return func(redisSettings *SentinelSettings) {
		redisSettings.Password = password
	}
}

// WithRedisDatabase sets redis database
func WithRedisDatabase(database int) SentinelOption {
	return func(redisSettings *SentinelSettings) {
		redisSettings.Database = database
	}
}

// WithRedisNetwork sets redis network
func WithRedisNetwork(network string) SentinelOption {
	return func(redisSettings *SentinelSettings) {
		redisSettings.Network = network
	}
}

// WithRedisConnectTimeout sets redis connect timeout
func WithRedisConnectTimeout(connectTimeout time.Duration) SentinelOption {
	return func(redisSettings *SentinelSettings) {
		redisSettings.ConnectTimeout = connectTimeout
	}
}

// WithRedisWriteTimeout sets redis write timeout
func WithRedisWriteTimeout(writeTimeout time.Duration) SentinelOption {
	return func(redisSettings *SentinelSettings) {
		redisSettings.WriteTimeout = writeTimeout
	}
}

// WithRedisReadTimeout sets redis read timeout
func WithRedisReadTimeout(readTimeout time.Duration) SentinelOption {
	return func(redisSettings *SentinelSettings) {
		redisSettings.ReadTimeout = readTimeout
	}
}

// WithRedisConnectionIdleTimeout sets redis connection idle timeout
func WithRedisConnectionIdleTimeout(connectionIdleTimeout time.Duration) SentinelOption {
	return func(redisSettings *SentinelSettings) {
		redisSettings.ConnectionIdleTimeout = connectionIdleTimeout
	}
}

// WithRedisConnectionMaxIdle sets redis connection max idle
func WithRedisConnectionMaxIdle(connectionMaxIdle int) SentinelOption {
	return func(redisSettings *SentinelSettings) {
		redisSettings.ConnectionMaxIdle = connectionMaxIdle
	}
}

// WithRedisConnectionMaxActive sets redis connection max active
func WithRedisConnectionMaxActive(connectionMaxActive int) SentinelOption {
	return func(redisSettings *SentinelSettings) {
		redisSettings.ConnectionMaxActive = connectionMaxActive
	}
}

// WithRedisConnectionWait sets redis connection wait
func WithRedisConnectionWait(connectionWait bool) SentinelOption {
	return func(redisSettings *SentinelSettings) {
		redisSettings.ConnectionWait = connectionWait
	}
}

// WithRedisTlsConfig sets redis tls config
func WithRedisTlsConfig(tlsConfig *tls.Config) SentinelOption {
	return func(redisSettings *SentinelSettings) {
		redisSettings.TlsConfig = tlsConfig
	}
}

// WithRedisTlsSkipVerify sets redis tls skip verification
func WithRedisTlsSkipVerify(tlsSkipVerify bool) SentinelOption {
	return func(redisSettings *SentinelSettings) {
		redisSettings.TlsSkipVerify = tlsSkipVerify
	}
}

// SentinelSettings contains redis settings
type SentinelSettings struct {
	Addresses             []string
	MasterName            string
	Password              string
	Database              int
	Network               string
	ConnectTimeout        time.Duration
	WriteTimeout          time.Duration
	ReadTimeout           time.Duration
	ConnectionIdleTimeout time.Duration
	ConnectionMaxIdle     int
	ConnectionMaxActive   int
	ConnectionWait        bool
	TlsConfig             *tls.Config
	TlsSkipVerify         bool
}

// NewSentinelCache creates a new redis client
func NewSentinelCache(options ...SentinelOption) *SentinelCache {
	settings := &SentinelSettings{}
	for _, option := range options {
		option(settings)
	}

	sentinelOptions := &xredis.SentinelOptions{
		Addresses:             settings.Addresses,
		MasterName:            settings.MasterName,
		Password:              settings.Password,
		Database:              settings.Database,
		Network:               settings.Network,
		ConnectTimeout:        settings.ConnectTimeout,
		WriteTimeout:          settings.WriteTimeout,
		ReadTimeout:           settings.ReadTimeout,
		ConnectionIdleTimeout: settings.ConnectionIdleTimeout,
		ConnectionMaxIdle:     settings.ConnectionMaxIdle,
		ConnectionMaxActive:   settings.ConnectionMaxActive,
		ConnectionWait:        settings.ConnectionWait,
		TlsConfig:             settings.TlsConfig,
		TlsSkipVerify:         settings.TlsSkipVerify,
	}
	return &SentinelCache{client: xredis.SetupSentinelClient(sentinelOptions)}
}

// SentinelCache contains redis client
type SentinelCache struct {
	client *xredis.Client
}

// Get a value by key
func (c *SentinelCache) Get(key string) (string, bool, error) {
	return c.client.Get(key)
}

// Set a key value pair
func (c *SentinelCache) Set(key string, value string, timeout time.Duration) error {
	seconds := int64(timeout.Seconds())
	if seconds == 0 {
		_, err := c.client.Set(key, value)
		return err
	}

	_, err := c.client.SetEx(key, value, seconds)
	return err
}

// Remove a key
func (c *SentinelCache) Remove(key string) error {
	_, err := c.client.Del(key)
	return err
}

// Ping to test connectivity
func (c *SentinelCache) Ping() error {
	_, err := c.client.Ping()
	return err
}

// Close to close resources
func (c *SentinelCache) Close() error {
	return c.client.Close()
}
