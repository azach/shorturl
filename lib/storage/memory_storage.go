package storage

import (
	"errors"
	"fmt"
	"time"
)

type MemoryStorage struct {
	items map[string]string
}

func NewMemoryStorage() Storage {
	items := make(map[string]string)
	return &MemoryStorage{
		items: items,
	}
}

func (c *MemoryStorage) Set(key string, value string) error {
	_, ok := c.items[key]
	if ok {
		return errors.New(fmt.Sprintf("value already exists for key %s", key))
	} else {
		c.items[key] = value
	}

	return nil
}

func (c *MemoryStorage) Get(key string) (value string, exists bool) {
	value, exists = c.items[key]
	return value, exists
}

// Not implemented
func (c *MemoryStorage) Hit(key string, viewedAt time.Time) {
	return
}

// Not implemented
func (c *MemoryStorage) GetHits(key string, asOf time.Time, hitRange HitRange) (int64, error) {
	return 0, nil
}
