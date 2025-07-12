package ports

import "context"

type GetFollowersUseCasePort interface {
	Execute(ctx context.Context, userID string) ([]string, error)
}
