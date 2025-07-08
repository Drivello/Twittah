package usecase

import (
	"context"
	"errors"

	"github.com/Drivello/Twittah/services/user/internal/ports"
)

type UnfollowUserUseCase struct {
	Repo   ports.UserRepository
	Events ports.UserEventPublisher
}

func NewUnfollowUserUseCase(repo ports.UserRepository, events ports.UserEventPublisher) *UnfollowUserUseCase {
	return &UnfollowUserUseCase{Repo: repo, Events: events}
}

func (uc *UnfollowUserUseCase) Execute(ctx context.Context, followerID, followeeID string) error {
	if followerID == followeeID {
		return errors.New("cannot unfollow yourself")
	}
	err := uc.Repo.UnfollowUser(ctx, followerID, followeeID)
	if err != nil {
		return err
	}
	// Publish event (fire and forget, but log error)
	if pubErr := uc.Events.PublishUnfollow(followerID, followeeID); pubErr != nil {
		// TODO: Connect with timeline
		// Log error if needed, or ignore for now
	}
	return nil
}
