package ports

import "context"

type UseCaseInterface interface {
	Execute(ctx context.Context, payload interface{}) error
}
