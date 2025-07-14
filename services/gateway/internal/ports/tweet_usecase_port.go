package ports

import (
	"context"

	"github.com/Drivello/Twittah/services/gateway/internal/domain"
)

type CreateTweetUseCasePort interface {
	Execute(ctx context.Context, authorID int64, content string) error
}

type DeleteTweetUseCasePort interface {
	Execute(ctx context.Context, tweetID int64) error
}

type GetTimelinePort interface {
	Execute(ctx context.Context, userID int64) ([]*domain.Tweet, error)
}

type GetTweetsFromIDsPort interface {
	Execute(ctx context.Context, userIDs []int64) ([]*domain.Tweet, error)
}

type GetTweetsFromUserIDPort interface {
	Execute(ctx context.Context, userID int64) ([]*domain.Tweet, error)
}
