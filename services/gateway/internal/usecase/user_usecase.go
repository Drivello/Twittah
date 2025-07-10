package usecase

import (
	"context"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
)

// UserUseCase implementa la lógica de usuario y cumple con el puerto hexagonal UserUseCasePort

type UserUseCase struct {
	Producer ports.UserEventProducerPort
}

var _ ports.UserUseCasePort = (*UserUseCase)(nil)

// NewUserUseCase crea un nuevo UserUseCase.
func NewUserUseCase(producer ports.UserEventProducerPort) *UserUseCase {
	return &UserUseCase{Producer: producer}
}

// FollowUser handles follow logic.
func (uc *UserUseCase) FollowUser(ctx context.Context, followerID, followeeID string) error {
	return uc.Producer.PublishFollow(ctx, followerID, followeeID)
}

// UnfollowUser handles unfollow logic.
func (uc *UserUseCase) UnfollowUser(ctx context.Context, followerID, followeeID string) error {
	return uc.Producer.PublishUnfollow(ctx, followerID, followeeID)
}
