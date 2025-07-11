package ports

import (
	"context"

	"github.com/Drivello/Twittah/services/auth/internal/domain"
)

type AuthUseCasesPort interface {
	RegisterUser(ctx context.Context, user domain.User) (int64, error)
}
