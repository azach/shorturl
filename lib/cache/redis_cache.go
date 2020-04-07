package cache

import "github.com/go-redis/redis/v7"

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache() Cache {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisCache{
		client: redisClient,
	}
}

func (c *RedisCache) Set(key string, value string) error {
	return c.client.Set(key, value, 0).Err()
}

func (c *RedisCache) Get(key string) (value string, exists bool) {
	value, err := c.client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", false
		}
		panic(err)
	}
	return value, true
}
