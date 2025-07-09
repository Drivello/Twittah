// Package postgres implements TweetRepository using Ent ORM.
package postgres

import (
	"context"

	// "github.com/Drivello/Twittah/services/auth/ent" // Disabled: ent package not generated
	"github.com/Drivello/Twittah/services/tweet/internal/domain"
	"github.com/Drivello/Twittah/services/tweet/internal/ports"
	//"github.com/Drivello/Twittah/services/tweet/ent" // Uncomment when ent is generated
)

// EntTweetRepository implements ports.TweetRepository using Ent.
// EntTweetRepository implements ports.TweetRepository using a stub (Ent disabled)
type EntTweetRepository struct {
	client interface{} // TODO: Replace with *ent.Client when Ent is generated
}

// NewEntTweetRepository creates a new EntTweetRepository.
func NewEntTweetRepository(client interface{}) *EntTweetRepository {
	return &EntTweetRepository{client: client}
} // TODO: Use *ent.Client when Ent is generated

// Save persists a tweet to the database.
func (r *EntTweetRepository) Save(ctx context.Context, tweet *domain.Tweet) error {
	// TODO: Implement with Ent ORM
	return nil // stub: does nothing
}

// GetRecentTweets returns recent tweets by user IDs.
func (r *EntTweetRepository) GetRecentTweets(ctx context.Context, userIDs []string, limit int) ([]domain.Tweet, error) {
	// TODO: Implement with Ent ORM
	return []domain.Tweet{}, nil // stub: returns empty
}

var _ ports.TweetRepository = (*EntTweetRepository)(nil)
