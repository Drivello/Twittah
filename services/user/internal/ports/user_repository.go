package ports

import "context"

// UserRepository defines the expected behavior for user data storage.
type UserRepository interface {
	InsertUser(ctx context.Context, id, username string) error
	FollowUser(ctx context.Context, followerID, followeeID string) error
	UnfollowUser(ctx context.Context, followerID, followeeID string) error
	GetFollowers(ctx context.Context, userID string) ([]string, error)
	GetFollowing(ctx context.Context, userID string) ([]string, error)
}
