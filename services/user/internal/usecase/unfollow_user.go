package usecase

import (
	"context"
	"errors"

	"github.com/Drivello/Twittah/services/user/internal/ports"
)

// UnfollowUserUseCase orchestrates the process of unfollowing a user.
type UnfollowUserUseCase struct {
	Repo ports.UserRepository // Outbound port for persistence
}

// NewUnfollowUserUseCase constructs an UnfollowUserUseCase with the given ports.
func NewUnfollowUserUseCase(repo ports.UserRepository) *UnfollowUserUseCase {
	return &UnfollowUserUseCase{Repo: repo}
}

// Execute performs the unfollow operation and publishes the event.
func (uc *UnfollowUserUseCase) Execute(ctx context.Context, followerID, followeeID int64) error {
	if followerID == followeeID {
		return errors.New("cannot unfollow yourself")
	}
	err := uc.Repo.UnfollowUser(ctx, followerID, followeeID)
	if err != nil {
		return err
	}
	return nil
}
