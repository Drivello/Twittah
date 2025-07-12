package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/gateway/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
)

// UserUseCase implementa la lógica de usuario y cumple con el puerto hexagonal UserUseCasePort

type FollowUseCase struct {
	Producer ports.EventProducerPort[kafka.UserPayload]
}

var _ ports.FollowUserUseCasePort = (*FollowUseCase)(nil)

// NewFollowUseCase crea un nuevo UserUseCase.
func NewFollowUseCase(producer ports.EventProducerPort[kafka.UserPayload]) *FollowUseCase {
	return &FollowUseCase{Producer: producer}
}

// FollowUser handles follow logic.
func (uc *FollowUseCase) Execute(ctx context.Context, followerID, followeeID int64) error {
	request := kafka.KafkaFollowPayload{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}
	return uc.Producer.PublishEvent("users.follow", request)
}
