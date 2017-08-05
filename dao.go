package fridge

import (
	"fmt"
	"github.com/shomali11/util/conversions"
	"github.com/shomali11/xredis"
	"time"
)

const (
	configKeyFormat = "%s.config"
)

// Dao controls access to redis
type Dao struct {
	xredisClient *xredis.Client
}

// Get retrieves an item
func (d *Dao) Get(key string) (string, bool, error) {
	return d.xredisClient.Get(key)
}

// Set stores a value
func (d *Dao) Set(key string, value string, timeout int64) error {
	_, err := d.xredisClient.SetEx(key, value, timeout)
	return err
}

// SetStorageDetails stores a key's defaults
func (d *Dao) SetStorageDetails(key string, storageDetails *StorageDetails) error {
	storageDetails.Timestamp = time.Now().UTC()
	timestampString, err := conversions.Stringify(storageDetails)
	if err != nil {
		return err
	}

	configKey := fmt.Sprintf(configKeyFormat, key)
	_, err = d.xredisClient.Set(configKey, timestampString)
	return err
}

// GetStorageDetails retrieves a key's storage details
func (d *Dao) GetStorageDetails(key string) (*StorageDetails, bool, error) {
	configKey := fmt.Sprintf(configKeyFormat, key)
	configString, found, err := d.xredisClient.Get(configKey)
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
	_, err := d.xredisClient.Del(key, timestampKey)
	return err
}

// Ping pings redis
func (d *Dao) Ping() error {
	_, err := d.xredisClient.Ping()
	return err
}

// Close closes resources
func (d *Dao) Close() error {
	return d.xredisClient.Close()
}

func newDao(redisClient *RedisClient) *Dao {
	return &Dao{xredisClient: redisClient.xredisClient}
}
