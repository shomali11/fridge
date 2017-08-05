package fridge

import (
	"crypto/tls"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisClient_WithHost(t *testing.T) {
	redisSettings := &RedisSettings{}

	redisOption := WithHost("host")
	redisOption(redisSettings)

	assert.Equal(t, redisSettings.Host, "host")
}

func TestRedisClient_WithPort(t *testing.T) {
	redisSettings := &RedisSettings{}

	redisOption := WithPort(1111)
	redisOption(redisSettings)

	assert.Equal(t, redisSettings.Port, 1111)
}

func TestRedisClient_WithPassword(t *testing.T) {
	redisSettings := &RedisSettings{}

	redisOption := WithPassword("password")
	redisOption(redisSettings)

	assert.Equal(t, redisSettings.Password, "password")
}

func TestRedisClient_WithDatabase(t *testing.T) {
	redisSettings := &RedisSettings{}

	redisOption := WithDatabase(1)
	redisOption(redisSettings)

	assert.Equal(t, redisSettings.Database, 1)
}

func TestRedisClient_WithNetwork(t *testing.T) {
	redisSettings := &RedisSettings{}

	redisOption := WithNetwork("network")
	redisOption(redisSettings)

	assert.Equal(t, redisSettings.Network, "network")
}

func TestRedisClient_WithConnectTimeout(t *testing.T) {
	redisSettings := &RedisSettings{}

	redisOption := WithConnectTimeout(time.Second)
	redisOption(redisSettings)

	assert.Equal(t, redisSettings.ConnectTimeout, time.Second)
}

func TestRedisClient_WithReadTimeout(t *testing.T) {
	redisSettings := &RedisSettings{}

	redisOption := WithReadTimeout(time.Second)
	redisOption(redisSettings)

	assert.Equal(t, redisSettings.ReadTimeout, time.Second)
}

func TestRedisClient_WithWriteTimeout(t *testing.T) {
	redisSettings := &RedisSettings{}

	redisOption := WithWriteTimeout(time.Second)
	redisOption(redisSettings)

	assert.Equal(t, redisSettings.WriteTimeout, time.Second)
}

func TestRedisClient_WithConnectionIdleTimeout(t *testing.T) {
	redisSettings := &RedisSettings{}

	redisOption := WithConnectionIdleTimeout(time.Second)
	redisOption(redisSettings)

	assert.Equal(t, redisSettings.ConnectionIdleTimeout, time.Second)
}

func TestRedisClient_WithConnectionMaxIdle(t *testing.T) {
	redisSettings := &RedisSettings{}

	redisOption := WithConnectionMaxIdle(11)
	redisOption(redisSettings)

	assert.Equal(t, redisSettings.ConnectionMaxIdle, 11)
}

func TestRedisClient_WithConnectionMaxActive(t *testing.T) {
	redisSettings := &RedisSettings{}

	redisOption := WithConnectionMaxActive(11)
	redisOption(redisSettings)

	assert.Equal(t, redisSettings.ConnectionMaxActive, 11)
}

func TestRedisClient_WithConnectionWait(t *testing.T) {
	redisSettings := &RedisSettings{}

	redisOption := WithConnectionWait(true)
	redisOption(redisSettings)

	assert.Equal(t, redisSettings.ConnectionWait, true)
}

func TestRedisClient_WithTlsConfig(t *testing.T) {
	redisSettings := &RedisSettings{}

	redisOption := WithTlsConfig(&tls.Config{})
	redisOption(redisSettings)

	assert.NotNil(t, redisSettings.TlsConfig)
}

func TestRedisClient_WithTlsSkipVerify(t *testing.T) {
	redisSettings := &RedisSettings{}

	redisOption := WithTlsSkipVerify(true)
	redisOption(redisSettings)

	assert.Equal(t, redisSettings.TlsSkipVerify, true)
}

func TestRedisClient_WithTestOnBorrowPeriod(t *testing.T) {
	redisSettings := &RedisSettings{}

	redisOption := WithTestOnBorrowPeriod(time.Second)
	redisOption(redisSettings)

	assert.Equal(t, redisSettings.TestOnBorrowPeriod, time.Second)
}
