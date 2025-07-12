package ports

import "context"

type FollowUserUseCasePort interface {
	Execute(ctx context.Context, followerID, followeeID int64) error
}

type UnfollowUserUseCasePort interface {
	Execute(ctx context.Context, followerID, followeeID int64) error
}
