package ports

import (
	"context"

	"github.com/Drivello/Twittah/services/gateway/internal/domain"
)

type TweetServicePort interface {
	GetTimeline(ctx context.Context, userID int64) ([]*domain.Tweet, error)
	GetTweetsFromMultipleUserIDs(ctx context.Context, userIDs []int64) ([]*domain.Tweet, error)
	GetTweetsFromUserID(ctx context.Context, userID int64) ([]*domain.Tweet, error)
}
