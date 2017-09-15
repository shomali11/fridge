package fridge

import (
	"time"
)

// StorageOption an option for a storage
type StorageOption func(*StorageDetails)

// WithDurations sets storage best by and use by durations
func WithDurations(bestBy time.Duration, useBy time.Duration) StorageOption {
	return func(storageInfo *StorageDetails) {
		storageInfo.BestBy = bestBy
		storageInfo.UseBy = useBy
	}
}

// StorageDetails contains storage information
type StorageDetails struct {
	Timestamp  time.Time
	Restocking bool
	BestBy     time.Duration
	UseBy      time.Duration
}

func newStorageDetails(defaults *Defaults, options ...StorageOption) *StorageDetails {
	storageDetails := &StorageDetails{BestBy: defaults.BestBy, UseBy: defaults.UseBy}
	for _, option := range options {
		option(storageDetails)
	}
	return storageDetails
}
