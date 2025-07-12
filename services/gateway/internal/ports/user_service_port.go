package ports

import "context"

type UserServicePort interface {
	GetFollowers(ctx context.Context, userID string) ([]string, error)
}
