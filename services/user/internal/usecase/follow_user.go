package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Drivello/Twittah/services/user/internal/ports"
)

// FollowUserUseCase orchestrates the process of following a user.
type FollowUserUseCase struct {
	Repo   ports.UserRepository      // Outbound port for persistence
	Events ports.UserEventPublisher  // Outbound port for event publishing
}

// NewFollowUserUseCase constructs a FollowUserUseCase with the given ports.
func NewFollowUserUseCase(repo ports.UserRepository, events ports.UserEventPublisher) *FollowUserUseCase {
	return &FollowUserUseCase{Repo: repo, Events: events}
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
	// Publish event (fire and forget, but log error)
	if pubErr := uc.Events.PublishFollow(fmt.Sprint(followerID), fmt.Sprint(followeeID)); pubErr != nil {
		// TODO: Connect with timeline
		// Log error if needed, or ignore for now
	}
	return nil
}
