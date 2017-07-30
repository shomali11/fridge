package item

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewConfig(t *testing.T) {
	config, err := NewConfig("name", time.Second, time.Minute)
	assert.Nil(t, err)
	assert.Equal(t, config.Key, "name")
	assert.Equal(t, config.BestBy, time.Second)
	assert.Equal(t, config.UseBy, time.Minute)

	_, err = NewConfig("name", time.Minute, time.Second)
	assert.NotNil(t, err)
}
