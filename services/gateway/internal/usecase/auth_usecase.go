package usecase

import (
	"context"
	"encoding/json"

	"github.com/Drivello/Twittah/services/gateway/internal/ports"
)

type UserEventRequest struct {
	EventType string                 `json:"event_type"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
}

// AuthUseCase implementa la lógica de autenticación y cumple con el puerto hexagonal AuthUseCasePort

type AuthUseCase struct {
	Producer ports.AuthEventProducerPort
}

var _ ports.AuthUseCasePort = (*AuthUseCase)(nil)

// NewAuthUseCase creates a new AuthUseCase.
func NewAuthUseCase(producer ports.AuthEventProducerPort) *AuthUseCase {
	return &AuthUseCase{Producer: producer}
}

// RegisterUser handles user registration logic.
func (uc *AuthUseCase) RegisterUser(ctx context.Context, username, email, password string) error {
	event := UserEventRequest{
		EventType: "users.create",
		Payload: map[string]interface{}{
			"username": username,
			"email":    email,
			"password": password,
		},
	}

eventBytes, err := json.Marshal(event)
if err != nil {
	return err
}
return uc.Producer.PublishUserCreateRequest(ctx, eventBytes)
}
