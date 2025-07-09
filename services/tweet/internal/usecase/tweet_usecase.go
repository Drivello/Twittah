// Package usecase implements TweetService business logic orchestrating domain and ports.
package usecase

import (
	"context"
	"github.com/Drivello/Twittah/services/tweet/internal/domain"
	"github.com/Drivello/Twittah/services/tweet/internal/ports"
)

// PublishTweet handles publishing a tweet, fan-out to timelines, and emitting events.
type PublishTweet struct {
	Repo  ports.TweetRepository
	Cache ports.Cache
}

// NewPublishTweet creates a new PublishTweet usecase.
func NewPublishTweet(repo ports.TweetRepository, cache ports.Cache) *PublishTweet {
	return &PublishTweet{Repo: repo, Cache: cache}
}

// Execute publishes a tweet and fans out to followers' timelines.
func (uc *PublishTweet) Execute(ctx context.Context, tweet *domain.Tweet, followerIDs []string) error {
	if err := tweet.Validate(); err != nil {
		return err
	}

	if err := uc.Repo.Save(ctx, tweet); err != nil {
		return err
	}

	if err := uc.Cache.AddTweetToTimelines(ctx, *tweet, followerIDs); err != nil {
		return err
	}

	return nil
}

// GetTimeline handles retrieving a user's timeline, using cache or DB.
type GetTimeline struct {
	Repo  ports.TweetRepository
	Cache ports.Cache
}

// NewGetTimeline creates a new GetTimeline usecase.
func NewGetTimeline(repo ports.TweetRepository, cache ports.Cache) *GetTimeline {
	return &GetTimeline{Repo: repo, Cache: cache}
}

// Execute retrieves the timeline for a user.
func (uc *GetTimeline) Execute(ctx context.Context, userID string, limit int) ([]domain.Tweet, error) {
	tweets, err := uc.Cache.GetTimeline(ctx, userID, limit)
	if err == nil && len(tweets) > 0 {
		return tweets, nil
	}

	// Cache miss: rebuild from repo
	tweets, err = uc.Repo.GetRecentTweets(ctx, []string{userID}, limit)
	if err != nil {
		return nil, err
	}

	if len(tweets) > 0 {
		// Cache the timeline (ignore cache errors)
		_ = uc.Cache.AddTweetToTimelines(ctx, tweets[0], []string{userID})
	}

	return tweets, nil
}
