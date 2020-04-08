package storage

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
)

type RedisStorage struct {
	client *redis.Client
	cache  Storage
}

func NewRedisStorage() Storage {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisStorage{
		client: redisClient,
		cache:  NewMemoryStorage(),
	}
}

func (c *RedisStorage) Set(key string, value string) error {
	err := c.cache.Set(key, value)
	if err != nil {
		logrus.Errorf("failed to set in cache: %v", err)
	}

	return c.client.Set(key, value, 0).Err()
}

func (c *RedisStorage) Get(key string) (value string, exists bool) {
	// Try fetch from memory cache first, then fall back to redis
	value, exists = c.cache.Get(key)
	if exists {
		logrus.Errorf("fetched value for %s from cache", key)
		return value, exists
	}

	value, err := c.client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", false
		} else {
			logrus.Errorf("error getting key: %v", err)
			return "", false
		}
	}

	// Set in memory cache for subsequent lookups
	err = c.cache.Set(key, value)
	if err != nil {
		logrus.Errorf("failed to set in cache: %v", err)
	}

	return value, true
}

func (c *RedisStorage) Hit(key string, viewedAt time.Time) {
	var wg sync.WaitGroup

	for _, hitRange := range []HitRange{AllTime, Weekly, Daily, Minute} {
		precision := toPrecision(hitRange)
		wg.Add(1)
		go c.hit(&wg, key, precision, viewedAt)
	}
	wg.Wait()
}

func (c *RedisStorage) hit(wg *sync.WaitGroup, key string, precision int64, viewedAt time.Time) {
	defer wg.Done()

	hashKey := toHashKey(key, precision)
	bucket := toBucket(viewedAt, precision)
	c.client.HIncrBy(hashKey, bucket, 1)
}

func (c *RedisStorage) GetHits(key string, asOf time.Time, hitRange HitRange) (int64, error) {
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
