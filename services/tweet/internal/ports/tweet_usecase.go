package ports

import (
	"context"

	"github.com/Drivello/Twittah/services/tweet/internal/domain"
)

type CreateTweetUsecasePort interface {
	CreateTweet(ctx context.Context, tweet *domain.Tweet) error
}

type DeleteTweetUsecasePort interface {
	DeleteTweet(ctx context.Context, tweetID int64) error
}
