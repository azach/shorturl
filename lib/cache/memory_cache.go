package cache

import (
	"errors"
	"fmt"
	"time"
)

type MemoryCache struct {
	items map[string]string
}

func NewMemoryCache() Cache {
	items := make(map[string]string)
	return &MemoryCache{
		items: items,
	}
}

func (c *MemoryCache) Set(key string, value string) error {
	_, ok := c.items[key]
	if ok {
		return errors.New(fmt.Sprintf("value already exists for key %s", key))
	} else {
		c.items[key] = value
	}

	return nil
}

func (c *MemoryCache) Get(key string) (value string, exists bool) {
	value, exists = c.items[key]
	return value, exists
}

func (c *MemoryCache) Hit(key string, viewedAt time.Time) {
	return
}

func (c *MemoryCache) GetHits(key string, asOf time.Time, hitRange HitRange) (int64, error) {
	return 0, nil
}
