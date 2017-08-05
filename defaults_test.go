package fridge

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDefaults_New(t *testing.T) {
	defaults := &Defaults{}

	storageOption := WithDefaultDurations(time.Second, 2*time.Second)
	storageOption(defaults)

	assert.Equal(t, defaults.BestBy, time.Second)
	assert.Equal(t, defaults.UseBy, 2*time.Second)
}

func TestDefaults_Override(t *testing.T) {
	defaults := newDefaults(WithDefaultDurations(time.Minute, 2*time.Minute))

	assert.Equal(t, defaults.BestBy, time.Minute)
	assert.Equal(t, defaults.UseBy, 2*time.Minute)
}
