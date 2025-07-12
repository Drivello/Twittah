package ports

import (
	"context"

	"github.com/Drivello/Twittah/services/user/internal/domain"
)

type CreateUserUseCasePort interface {
	Execute(ctx context.Context, user *domain.User) error
}
