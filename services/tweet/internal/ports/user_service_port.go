package ports

import (
	"context"

	"github.com/Drivello/Twittah/services/tweet/internal/domain"
)

type UserServicePort interface {
	GetFollowing(ctx context.Context, userID int64) ([]*domain.User, error)
}
