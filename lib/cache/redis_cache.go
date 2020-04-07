package cache

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache() Cache {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
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
		} else {
			logrus.Errorf("error getting key: %v", err)
			return "", false
		}
	}
	return value, true
}

// TODO: this can be parallelized or backgrounded
func (c *RedisCache) Hit(key string, viewedAt time.Time) {
	for _, hitRange := range []HitRange{AllTime, Weekly, Daily, Minute} {
		precision := toPrecision(hitRange)
		hashKey := toHashKey(key, precision)
		bucket := toBucket(viewedAt, precision)
		c.client.HIncrBy(hashKey, bucket, 1)
	}
}

func (c *RedisCache) GetHits(key string, asOf time.Time, hitRange HitRange) (int64, error) {
	precision := toPrecision(hitRange)
	hashKey := toHashKey(key, precision)
	bucket := toBucket(asOf, precision)
	res, err := c.client.HGet(hashKey, bucket).Result()
	if err == redis.Nil {
		return 0, nil
	}
	count, err := strconv.ParseInt(res, 10, 64)
	return count, err
}

func toHashKey(key string, precision int64) string {
	return fmt.Sprintf("count:%v:%v", precision, key)
}

func toBucket(timestamp time.Time, precision int64) string {
	if precision == 0 {
		return "0"
	}
	return strconv.FormatInt(int64(timestamp.Unix()/precision)*precision, 10)
}
