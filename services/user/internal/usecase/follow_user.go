package usecase

import (
	"context"
	"errors"

	"github.com/Drivello/Twittah/services/user/internal/ports"
)

type FollowUserUseCase struct {
	Repo   ports.UserRepository
	Events ports.UserEventPublisher
}

func NewFollowUserUseCase(repo ports.UserRepository, events ports.UserEventPublisher) *FollowUserUseCase {
	return &FollowUserUseCase{Repo: repo, Events: events}
}

func (uc *FollowUserUseCase) Execute(ctx context.Context, followerID, followeeID string) error {
	if followerID == followeeID {
		return errors.New("cannot follow yourself")
	}
	err := uc.Repo.FollowUser(ctx, followerID, followeeID)
	if err != nil {
		return err
	}
	// Publish event (fire and forget, but log error)
	if pubErr := uc.Events.PublishFollow(followerID, followeeID); pubErr != nil {
		// TODO: Connect with timeline
		// Log error if needed, or ignore for now
	}
	return nil
}
