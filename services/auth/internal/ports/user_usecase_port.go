package ports

import (
	"context"

	"github.com/Drivello/Twittah/services/auth/internal/domain"
)

type UserUseCasePort interface {
	RegisterUser(ctx context.Context, user *domain.User) (int64, error)
}
