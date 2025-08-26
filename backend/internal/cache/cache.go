package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	ctx    context.Context
}

func NewCache() *Cache {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &Cache{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (c *Cache) Set(key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(c.ctx, key, data, ttl).Err()
}

func (c *Cache) Get(key string, dest any) (bool, error) {
	data, err := c.client.Get(c.ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return false, err
	}

	return true, nil
}

func (c *Cache) InvalidateUserCategories(userID uint) error {
	pattern := fmt.Sprintf("categories:%d:*", userID)
	iter := c.client.Scan(c.ctx, 0, pattern, 0).Iterator()
	for iter.Next(c.ctx) {
		if err := c.client.Del(c.ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}
