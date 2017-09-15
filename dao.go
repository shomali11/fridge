package fridge

import (
	"fmt"
	"github.com/shomali11/util/conversions"
	"time"
)

const (
	configKeyFormat = "%s.config"
)

// Dao controls access to redis
type Dao struct {
	cache Cache
}

// Get retrieves an item
func (d *Dao) Get(key string) (string, bool, error) {
	return d.cache.Get(key)
}

// Set stores a value
func (d *Dao) Set(key string, value string, timeout time.Duration) error {
	return d.cache.Set(key, value, timeout)
}

// SetStorageDetails stores a key's defaults
func (d *Dao) SetStorageDetails(key string, storageDetails *StorageDetails) error {
	storageDetails.Timestamp = time.Now().UTC()
	timestampString, err := conversions.Stringify(storageDetails)
	if err != nil {
		return err
	}

	configKey := fmt.Sprintf(configKeyFormat, key)
	return d.cache.Set(configKey, timestampString, 0)
}

// GetStorageDetails retrieves a key's storage details
func (d *Dao) GetStorageDetails(key string) (*StorageDetails, bool, error) {
	configKey := fmt.Sprintf(configKeyFormat, key)
	configString, found, err := d.cache.Get(configKey)
	if err != nil {
		return nil, false, err
	}

	if !found {
		return nil, false, nil
	}

	var storageDetails *StorageDetails
	err = conversions.Structify(configString, &storageDetails)
	if err != nil {
		return nil, false, err
	}
	return storageDetails, true, nil
}

// Remove an item
func (d *Dao) Remove(key string) error {
	timestampKey := fmt.Sprintf(configKeyFormat, key)
	err := d.cache.Remove(key)
	if err != nil {
		return err
	}
	return d.cache.Remove(timestampKey)
}

// Ping pings redis
func (d *Dao) Ping() error {
	return d.cache.Ping()
}

// Close closes resources
func (d *Dao) Close() error {
	return d.cache.Close()
}

func newDao(cache Cache) *Dao {
	return &Dao{cache: cache}
}
