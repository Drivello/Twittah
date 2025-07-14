package ports

import (
	"context"

	"github.com/Drivello/Twittah/services/tweet/internal/domain"
)

type GetTimelinePort interface {
	Execute(ctx context.Context, userID int64) ([]*domain.Tweet, error)
}

type GetTweetsFromMultipleUserIDsPort interface {
	Execute(ctx context.Context, userIDs []int64) ([]*domain.Tweet, error)
}

type GetUserTweetsPort interface {
	Execute(ctx context.Context, userID int64) ([]*domain.Tweet, error)
}
