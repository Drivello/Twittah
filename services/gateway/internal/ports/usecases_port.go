package ports

import (
	"context"

	"github.com/Drivello/Twittah/services/gateway/internal/domain"
)

// AuthUseCasePort define las operaciones de negocio para autenticación
type AuthUseCasePort interface {
	RegisterUser(ctx context.Context, username, email, password string) error
}

// UserQueryPort defines the contract for querying followers
type UserQueryPort interface {
	GetFollowers(c context.Context, userID string) ([]string, error)
}

// UserUseCasePort define las operaciones de negocio para usuarios
type UserUseCasePort interface {
	FollowUser(ctx context.Context, followerID, followeeID string) error
	UnfollowUser(ctx context.Context, followerID, followeeID string) error
}

// TweetUseCasePort define las operaciones de negocio para tweets
type TweetUseCasePort interface {
	PublishTweet(ctx context.Context, event domain.KafkaEventRequest) error
}
