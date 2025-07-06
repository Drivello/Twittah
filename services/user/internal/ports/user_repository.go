package ports

import "context"

// UserRepository defines the expected behavior for user data storage.
type UserRepository interface {
	FollowUser(ctx context.Context, followerID, followeeID string) error
	UnfollowUser(ctx context.Context, followerID, followeeID string) error
	GetFollowers(ctx context.Context, userID string) ([]string, error)
	GetFollowing(ctx context.Context, userID string) ([]string, error)
	Exists(ctx context.Context, userID string) (bool, error)
}
