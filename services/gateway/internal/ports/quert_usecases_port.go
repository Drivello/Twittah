package ports

import "context"

type UserQueryPort interface {
	GetFollowers(ctx context.Context, userID string) ([]string, error)
}
