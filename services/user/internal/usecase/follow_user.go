package usecase

import (
	"context"
	"errors"

	"github.com/Drivello/Twittah/services/user/internal/ports"
)

// FollowUserUseCase orchestrates the process of following a user.
type FollowUserUseCase struct {
	Repo ports.UserRepository // Outbound port for persistence
}

// NewFollowUserUseCase constructs a FollowUserUseCase with the given ports.
func NewFollowUserUseCase(repo ports.UserRepository) *FollowUserUseCase {
	return &FollowUserUseCase{Repo: repo}
}

// Execute performs the follow operation and publishes the event.
func (uc *FollowUserUseCase) Execute(ctx context.Context, followerID, followeeID int64) error {
	if followerID == followeeID {
		return errors.New("cannot follow yourself")
	}
	err := uc.Repo.FollowUser(ctx, followerID, followeeID)
	if err != nil {
		return err
	}
	return nil
}
