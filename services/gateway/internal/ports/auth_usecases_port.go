package ports

import (
	"context"
)

type RegisterUserUseCasePort interface {
	Execute(ctx context.Context, username, email, password string) error
}
