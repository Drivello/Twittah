package ports

import (
	"context"

	"github.com/Drivello/Twittah/services/auth/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (int64, error)
}
