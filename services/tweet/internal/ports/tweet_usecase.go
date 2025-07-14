package ports

import (
	"context"

	"github.com/Drivello/Twittah/services/tweet/internal/domain"
)

type CreateTweetUsecasePort interface {
	Execute(ctx context.Context, tweet *domain.Tweet) error
}

type DeleteTweetUsecasePort interface {
	Execute(ctx context.Context, tweetID int64) error
}

type GetTweetsFromIDsPort interface {
	Execute(ctx context.Context, userIDs []int64) ([]*domain.Tweet, error)
}
