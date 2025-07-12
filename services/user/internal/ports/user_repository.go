// Package ports defines interfaces (ports) for driving and driven adapters in the hexagonal architecture.
package ports

import (
	"context"

	"github.com/Drivello/Twittah/services/user/internal/domain"
)

// UserRepository defines the contract for user data storage and retrieval.
// This is an outbound port for persistence.
type UserRepository interface {
	// InsertUser stores a new user if not exists.
	InsertUser(ctx context.Context, user *domain.User) error
	// FollowUser creates a follow relationship.
	FollowUser(ctx context.Context, followerID, followeeID int64) error
	// UnfollowUser removes a follow relationship.
	UnfollowUser(ctx context.Context, followerID, followeeID int64) error
	// GetFollowers returns the IDs of users who follow the given user.
	GetFollowers(ctx context.Context, userID int64) ([]int64, error)
	// GetFollowing returns the IDs of users followed by the given user.
	GetFollowing(ctx context.Context, userID int64) ([]int64, error)
}
