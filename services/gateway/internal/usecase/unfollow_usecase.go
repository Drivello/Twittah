package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/gateway/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
)

// UserUseCase implementa la lógica de usuario y cumple con el puerto hexagonal UserUseCasePort

type UnfollowUseCase struct {
	Producer ports.EventProducerPort[kafka.UserPayload]
}

var _ ports.UnfollowUserUseCasePort = (*UnfollowUseCase)(nil)

// NewUnFollowUseCase crea un nuevo UserUseCase.
func NewUnfollowUseCase(producer ports.EventProducerPort[kafka.UserPayload]) *UnfollowUseCase {
	return &UnfollowUseCase{Producer: producer}
}

// Execute handles unfollow logic.
func (uc *UnfollowUseCase) Execute(ctx context.Context, followerID, followeeID int64) error {
	request := kafka.KafkaFollowPayload{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}
	return uc.Producer.PublishEvent("users.unfollow", request)
}
