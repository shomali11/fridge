package fridge

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStorageDetails_New(t *testing.T) {
	storageDetails := &StorageDetails{}

	storageOption := WithDurations(time.Second, 2*time.Second)
	storageOption(storageDetails)

	assert.Equal(t, storageDetails.BestBy, time.Second)
	assert.Equal(t, storageDetails.UseBy, 2*time.Second)

	assert.Equal(t, storageDetails.GetUseByInSeconds(), int64(2))
}

func TestStorageDetails_Defaults(t *testing.T) {
	defaults := &Defaults{
		BestBy: time.Minute,
		UseBy:  2 * time.Minute,
	}

	storageDetails := newStorageDetails(defaults)

	assert.Equal(t, storageDetails.BestBy, time.Minute)
	assert.Equal(t, storageDetails.UseBy, 2*time.Minute)

	assert.Equal(t, storageDetails.GetUseByInSeconds(), int64(120))
}

func TestStorageDetails_Override(t *testing.T) {
	defaults := &Defaults{
		BestBy: time.Minute,
		UseBy:  2 * time.Minute,
	}

	storageDetails := newStorageDetails(defaults, WithDurations(time.Second, 2*time.Second))

	assert.Equal(t, storageDetails.BestBy, time.Second)
	assert.Equal(t, storageDetails.UseBy, 2*time.Second)

	assert.Equal(t, storageDetails.GetUseByInSeconds(), int64(2))
}
