// Package redis implements Cache for timelines using go-redis v8.
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Drivello/Twittah/services/tweet/internal/domain"
	"github.com/Drivello/Twittah/services/tweet/internal/ports" // redis "github.com/go-redis/redis/v8" // Uncomment when go-redis is available

	redis "github.com/go-redis/redis/v8"
)

// TimelineCache implements ports.Cache using go-redis.
type TimelineCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewTimelineCache creates a new TimelineCache.
func NewTimelineCache(client *redis.Client, ttlHours int) *TimelineCache {
	return &TimelineCache{client: client, ttl: time.Duration(ttlHours) * time.Hour}
}

// AddTweetToTimelines adds a tweet to each follower's timeline in Redis (ZADD per follower).
func (c *TimelineCache) AddTweetToTimelines(ctx context.Context, tweet domain.Tweet, followerIDs []string) error {
	data, err := json.Marshal(tweet)
	if err != nil {
		return err
	}
	for _, fid := range followerIDs {
		key := fmt.Sprintf("timeline:%s", fid)
		z := &redis.Z{
			Score:  float64(tweet.CreatedAt.Unix()),
			Member: data,
		}
		if err := c.client.ZAdd(ctx, key, z).Err(); err != nil {
			return err
		}
		if err := c.client.Expire(ctx, key, c.ttl).Err(); err != nil {
			return err
		}
	}
	return nil
}

// GetTimeline retrieves a user's timeline from Redis (ZRANGE with scores, newest first).
func (c *TimelineCache) GetTimeline(ctx context.Context, userID string, limit int) ([]domain.Tweet, error) {
	key := fmt.Sprintf("timeline:%s", userID)
	vals, err := c.client.ZRevRange(ctx, key, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}
	var tweets []domain.Tweet
	for _, v := range vals {
		var t domain.Tweet
		if err := json.Unmarshal([]byte(v), &t); err == nil {
			tweets = append(tweets, t)
		}
	}
	return tweets, nil
}

// RemoveTweetsFromTimeline removes all tweets by authorID from userID's timeline.
func (c *TimelineCache) RemoveTweetsFromTimeline(ctx context.Context, userID, authorID string) error {
	key := fmt.Sprintf("timeline:%s", userID)
	vals, err := c.client.ZRange(ctx, key, 0, -1).Result()
	if err != nil {
		return err
	}
	for _, v := range vals {
		var t domain.Tweet
		if err := json.Unmarshal([]byte(v), &t); err == nil && t.AuthorID == authorID {
			c.client.ZRem(ctx, key, v)
		}
	}
	return nil
}

var _ ports.Cache = (*TimelineCache)(nil)
