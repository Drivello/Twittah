package redis

import (
	"context"
	"time"
	"github.com/redis/go-redis/v9"
)

type TimelineCache struct {
	client *redis.Client
}

func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{Addr: addr})
}

func NewTimelineCache(client *redis.Client) *TimelineCache {
	return &TimelineCache{client: client}
}

func (c *TimelineCache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *TimelineCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}
