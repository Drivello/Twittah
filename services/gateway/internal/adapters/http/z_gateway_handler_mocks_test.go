package http_test

import (
	"context"

	"github.com/Drivello/Twittah/services/gateway/internal/domain"
)

type mockRegisterUserUseCase struct {
	ExecuteFunc func(ctx context.Context, username, email, password string) error
}

func (m *mockRegisterUserUseCase) Execute(ctx context.Context, username, email, password string) error {
	return m.ExecuteFunc(ctx, username, email, password)
}

type mockCreateTweetUseCase struct {
	ExecuteFunc func(ctx context.Context, authorID int64, content string) error
}

func (m *mockCreateTweetUseCase) Execute(ctx context.Context, authorID int64, content string) error {
	return m.ExecuteFunc(ctx, authorID, content)
}

type mockDeleteTweetUseCase struct {
	ExecuteFunc func(ctx context.Context, tweetID int64) error
}

func (m *mockDeleteTweetUseCase) Execute(ctx context.Context, tweetID int64) error {
	return m.ExecuteFunc(ctx, tweetID)
}

type mockGetTimelineUseCase struct {
	ExecuteFunc func(ctx context.Context, userID int64) ([]*domain.Tweet, error)
}

func (m *mockGetTimelineUseCase) Execute(ctx context.Context, userID int64) ([]*domain.Tweet, error) {
	return m.ExecuteFunc(ctx, userID)
}

type mockGetTweetsFromIDsUseCase struct {
	ExecuteFunc func(ctx context.Context, ids []int64) ([]*domain.Tweet, error)
}

func (m *mockGetTweetsFromIDsUseCase) Execute(ctx context.Context, ids []int64) ([]*domain.Tweet, error) {
	return m.ExecuteFunc(ctx, ids)
}

type mockGetTweetsFromUserIDUseCase struct {
	ExecuteFunc func(ctx context.Context, userID int64) ([]*domain.Tweet, error)
}

func (m *mockGetTweetsFromUserIDUseCase) Execute(ctx context.Context, userID int64) ([]*domain.Tweet, error) {
	return m.ExecuteFunc(ctx, userID)
}

type mockFollowUserUseCase struct {
	ExecuteFunc func(ctx context.Context, followerID, followeeID int64) error
}

func (m *mockFollowUserUseCase) Execute(ctx context.Context, followerID, followeeID int64) error {
	return m.ExecuteFunc(ctx, followerID, followeeID)
}

type mockUnfollowUserUseCase struct {
	ExecuteFunc func(ctx context.Context, followerID, followeeID int64) error
}

func (m *mockUnfollowUserUseCase) Execute(ctx context.Context, followerID, followeeID int64) error {
	return m.ExecuteFunc(ctx, followerID, followeeID)
}

type mockGetFollowersUseCase struct {
	ExecuteFunc func(ctx context.Context, userID int64) ([]*domain.User, error)
}

func (m *mockGetFollowersUseCase) Execute(ctx context.Context, userID int64) ([]*domain.User, error) {
	return m.ExecuteFunc(ctx, userID)
}
