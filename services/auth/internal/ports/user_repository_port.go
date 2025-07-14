package ports

import (
	"context"

	"github.com/Drivello/Twittah/services/auth/internal/domain"
)

type AuthRepositoryPort interface {
	CreateOrGetUser(ctx context.Context, user domain.User) (int64, bool, error)
}
