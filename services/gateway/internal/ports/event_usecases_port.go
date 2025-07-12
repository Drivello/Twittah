package ports

import (
	"context"
)

type RegisterUserUseCasePort interface {
	Execute(ctx context.Context, username, email, password string) error
}

type FollowUserUseCasePort interface {
	Execute(ctx context.Context, followerID, followeeID int64) error
}

type UnfollowUserUseCasePort interface {
	Execute(ctx context.Context, followerID, followeeID int64) error
}

type CreateTweetUseCasePort interface {
	Execute(ctx context.Context, authorID int64, content string) error
}
