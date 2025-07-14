package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Drivello/Twittah/services/tweet/internal/domain"
	"github.com/Drivello/Twittah/services/tweet/internal/common"
	redis "github.com/redis/go-redis/v9"
)

// TimelineCache implementa ports.TimelineCachePort usando Redis.
type TimelineCache struct {
	client *redis.Client
	defaultTTL time.Duration
}

func NewTimelineCache(client *redis.Client, ttlHours int) *TimelineCache {
	return &TimelineCache{
		client: client,
		defaultTTL: time.Duration(ttlHours) * time.Hour,
	}
}

func (c *TimelineCache) GetTimeline(ctx context.Context, userID string) ([]domain.Tweet, error) {
	key := fmt.Sprintf("timeline:%s", userID)
	var lastErr error
	var data string
	for attempt, backoff := range []time.Duration{100 * time.Millisecond, 200 * time.Millisecond, 400 * time.Millisecond} {
		res, err := c.client.Get(ctx, key).Result()
		if err == redis.Nil {
			common.Logger().Infof("Cache MISS for user %s", userID)
			return nil, nil
		} else if err != nil {
			lastErr = err
			common.Logger().Warnf("Redis GET failed (attempt %d): %v", attempt+1, err)
			time.Sleep(backoff)
			continue
		}
		data = res
		break
	}
	if data == "" && lastErr != nil {
		common.Logger().Errorf("Redis GET failed after retries for user %s: %v", userID, lastErr)
		return nil, lastErr
	}
	common.Logger().Infof("Cache HIT for user %s", userID)
	var tweets []domain.Tweet
	if err := json.Unmarshal([]byte(data), &tweets); err != nil {
		common.Logger().Warnf("Failed to unmarshal cached timeline for user %s: %v", userID, err)
		return nil, err
	}
	return tweets, nil
}

func (c *TimelineCache) SetTimeline(ctx context.Context, userID string, tweets []domain.Tweet) error {
	key := fmt.Sprintf("timeline:%s", userID)
	bytes, err := json.Marshal(tweets)
	if err != nil {
		common.Logger().Warnf("Failed to marshal timeline for user %s: %v", userID, err)
		return err
	}
	var lastErr error
	for attempt, backoff := range []time.Duration{100 * time.Millisecond, 200 * time.Millisecond, 400 * time.Millisecond} {
		err = c.client.Set(ctx, key, bytes, c.defaultTTL).Err()
		if err == nil {
			return nil
		}
		lastErr = err
		common.Logger().Warnf("Redis SET failed (attempt %d) for user %s: %v", attempt+1, userID, err)
		time.Sleep(backoff)
	}
	common.Logger().Errorf("Redis SET failed after retries for user %s: %v", userID, lastErr)
	return lastErr
}

func (c *TimelineCache) InvalidateTimeline(ctx context.Context, userID string) error {
	key := fmt.Sprintf("timeline:%s", userID)
	var lastErr error
	for attempt, backoff := range []time.Duration{100 * time.Millisecond, 200 * time.Millisecond, 400 * time.Millisecond} {
		err := c.client.Del(ctx, key).Err()
		if err == nil {
			common.Logger().Infof("Invalidated timeline cache for user %s", userID)
			return nil
		}
		lastErr = err
		common.Logger().Warnf("Redis DEL failed (attempt %d) for user %s: %v", attempt+1, userID, err)
		time.Sleep(backoff)
	}
	common.Logger().Errorf("Redis DEL failed after retries for user %s: %v", userID, lastErr)
	return lastErr
}

// NewRedisClient crea un cliente Redis listo para usar.
func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{Addr: addr})
}
