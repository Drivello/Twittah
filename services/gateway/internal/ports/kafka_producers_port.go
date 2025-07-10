package ports

import (
	"context"

	"github.com/Drivello/Twittah/services/gateway/internal/domain"
)

// AuthEventProducerPort defines the contract for producing user events.
type AuthEventProducerPort interface {
	PublishUserCreateRequest(ctx context.Context, username, email, password string) error
}

// UserEventProducerPort defines the contract for producing follow events.
type UserEventProducerPort interface {
	PublishFollow(ctx context.Context, followerID, followeeID string) error
	PublishUnfollow(ctx context.Context, followerID, followeeID string) error
}

type TweetEventProducerPort interface {
	PublishTweet(ctx context.Context, event domain.TweetPublishEventDTO) error
}
