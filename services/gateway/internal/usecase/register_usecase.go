package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/gateway/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
)

// AuthUseCase implementa la lógica de autenticación y cumple con el puerto hexagonal AuthUseCasePort

type RegisterUserUseCase struct {
	Producer ports.EventProducerPort[kafka.AuthPayload]
}

var _ ports.RegisterUserUseCasePort = (*RegisterUserUseCase)(nil)

// NewRegisterUserUseCase creates a new AuthUseCase.
func NewRegisterUserUseCase(producer ports.EventProducerPort[kafka.AuthPayload]) *RegisterUserUseCase {
	return &RegisterUserUseCase{Producer: producer}
}

// Execute handles user registration logic.
func (uc *RegisterUserUseCase) Execute(ctx context.Context, username, email, password string) error {
	event := kafka.KafkaEventRequest[kafka.AuthPayload]{
		EventType: "users.create",
		Payload: kafka.KafkaUserCreatePayload{
			Username: username,
			Email:    email,
			Password: password,
		},
	}
	return uc.Producer.PublishEvent(event)
}
