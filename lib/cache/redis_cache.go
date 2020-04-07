package cache

type RedisCache struct{}

func NewRedisCache() Cache {
	// not implemented
	return &RedisCache{}
}

func (c *RedisCache) Set(key string, value string) error {
	// not implemented
	return nil
}

func (c *RedisCache) Get(key string) (value string, exists bool) {
	// not implemented
	return "", true
}
