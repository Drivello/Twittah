package ports

import (
	"context"
	"github.com/Drivello/Twittah/services/gateway/internal/domain"
)

type GetFollowersUseCasePort interface {
	Execute(ctx context.Context, userID string) ([]*domain.User, error)
}
