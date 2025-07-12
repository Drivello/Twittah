package ports

import (
	"context"
	"github.com/Drivello/Twittah/services/gateway/internal/domain"
)

type UserServicePort interface {
	GetFollowers(ctx context.Context, userID string) ([]*domain.User, error)
}
