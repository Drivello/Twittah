package usecase

import (
	"context"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
)

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
	return uc.Producer.PublishUserCreateRequest(ctx, username, email, password)
}
