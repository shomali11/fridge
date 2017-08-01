package item

import (
	"fmt"
	"github.com/shomali11/util/conversions"
	"github.com/shomali11/xredis"
	"time"
)

const (
	timestampKeyFormat = "%s.timestamp"
)

// NewDao creates a new dao object
func NewDao(client *xredis.Client) *Dao {
	return &Dao{xredisClient: client}
}

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

// GetTimestamp retrieves when an item was stocked
func (d *Dao) GetTimestamp(key string) (time.Time, error) {
	timestampKey := fmt.Sprintf(timestampKeyFormat, key)
	timestampString, found, err := d.xredisClient.Get(timestampKey)
	if err != nil {
		return time.Time{}, err
	}

	if !found {
		return time.Time{}, nil
	}

	var timestamp time.Time
	err = conversions.Structify(timestampString, &timestamp)
	if err != nil {
		return time.Time{}, err
	}
	return timestamp, nil
}

// SetTimestamp stores a new timestamp for the value
func (d *Dao) SetTimestamp(key string, timestamp time.Time) error {
	timestampString, err := conversions.Stringify(timestamp)
	if err != nil {
		return err
	}

	timestampKey := fmt.Sprintf(timestampKeyFormat, key)
	_, err = d.xredisClient.Set(timestampKey, timestampString)
	if err != nil {
		return err
	}
	return nil
}

// Remove an item
func (d *Dao) Remove(key string) error {
	_, err := d.xredisClient.Del(key)
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
