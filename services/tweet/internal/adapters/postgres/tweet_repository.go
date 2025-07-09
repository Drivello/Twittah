// Package postgres implements TweetRepository using Ent ORM.
package postgres

import (
	"context"

	// "github.com/Drivello/Twittah/services/auth/ent" // Disabled: ent package not generated
	"github.com/Drivello/Twittah/services/tweet/internal/domain"
	"github.com/Drivello/Twittah/services/tweet/internal/ports"
	"github.com/Drivello/Twittah/services/tweet/ent"
	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"go.uber.org/zap"
	"strconv"
	"github.com/Drivello/Twittah/services/tweet/ent/tweet"
)

// EntTweetRepository implements ports.TweetRepository using Ent.
// EntTweetRepository implements ports.TweetRepository using a stub (Ent disabled)
type EntTweetRepository struct {
	client *ent.Client
}

// NewEntTweetRepository creates a new EntTweetRepository.
func NewEntTweetRepository(client *ent.Client) *EntTweetRepository {
	return &EntTweetRepository{client: client}
}

// Save persists a tweet to the database.
func (r *EntTweetRepository) Save(ctx context.Context, tweet *domain.Tweet) error {
	logger := common.Logger()
	_, err := r.client.Tweet.
		Create().
		SetAuthorID(tweet.AuthorID).
		SetContent(tweet.Content).
		SetCreatedAt(tweet.CreatedAt).
		Save(ctx)
	if err != nil {
		logger.Error("failed to persist tweet", zap.Error(err), zap.String("author_id", tweet.AuthorID))
		return err
	}
	logger.Info("tweet persisted", zap.String("author_id", tweet.AuthorID))
	return nil
}

// GetRecentTweets returns recent tweets by user IDs.
func (r *EntTweetRepository) GetRecentTweets(ctx context.Context, userIDs []string, limit int) ([]domain.Tweet, error) {
	logger := common.Logger()
	query := r.client.Tweet.Query().
		Where(tweet.AuthorIDIn(userIDs...)).
		Order(ent.Desc("created_at")).
		Limit(limit)
	tweets, err := query.All(ctx)
	if err != nil {
		logger.Error("failed to query tweets", zap.Error(err))
		return nil, err
	}
	var result []domain.Tweet
	for _, t := range tweets {
		result = append(result, domain.Tweet{
			ID:        strconv.Itoa(t.ID),
			AuthorID:  t.AuthorID,
			Content:   t.Content,
			CreatedAt: t.CreatedAt,
		})
	}
	logger.Info("queried tweets", zap.Int("count", len(result)))
	return result, nil
}

var _ ports.TweetRepository = (*EntTweetRepository)(nil)
