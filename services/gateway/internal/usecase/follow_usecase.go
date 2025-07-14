package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/gateway/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
)

// FollowUseCase implements the user logic and fulfills the hexagonal port FollowUserUseCasePort.

type FollowUseCase struct {
	Producer ports.EventProducerPort[kafka.UserPayload]
}

var _ ports.FollowUserUseCasePort = (*FollowUseCase)(nil)

// NewFollowUseCase creates a new FollowUseCase.
func NewFollowUseCase(producer ports.EventProducerPort[kafka.UserPayload]) *FollowUseCase {
	return &FollowUseCase{Producer: producer}
}

// FollowUser handles follow logic.
func (uc *FollowUseCase) Execute(ctx context.Context, followerID, followeeID int64) error {
	event := kafka.KafkaEventRequest[kafka.UserPayload]{
		EventType: "users.follow",
		Payload: kafka.KafkaFollowPayload{
			FollowerID: followerID,
			FolloweeID: followeeID,
		},
	}
	return uc.Producer.PublishEvent(event)
}
