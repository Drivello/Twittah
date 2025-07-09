// Package ports defines the interfaces (ports) for TweetService adapters.
package ports

import (
	"context"
	"github.com/Drivello/Twittah/services/tweet/internal/domain"
)

// TweetRepository provides persistence operations for tweets.
type TweetRepository interface {
	// Save persists a tweet to storage.
	Save(ctx context.Context, tweet *domain.Tweet) error
	// GetRecentTweets returns recent tweets by user IDs.
	GetRecentTweets(ctx context.Context, userIDs []string, limit int) ([]domain.Tweet, error)
}

// Cache provides timeline cache operations for users.
type Cache interface {
	// AddTweetToTimelines adds a tweet to the timelines of given followers.
	AddTweetToTimelines(ctx context.Context, tweet domain.Tweet, followerIDs []string) error
	// GetTimeline returns the timeline for a user from cache.
	GetTimeline(ctx context.Context, userID string, limit int) ([]domain.Tweet, error)
	// RemoveTweetsFromTimeline removes all tweets by an author from a user's timeline.
	RemoveTweetsFromTimeline(ctx context.Context, userID, authorID string) error
}
