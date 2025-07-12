package kafka

import (
	"context"

	"github.com/Drivello/Twittah/services/auth/internal/domain"
	"github.com/Drivello/Twittah/services/auth/internal/ports"
)

// Handler for user creation events (modularizado)
func HandleUserCreate(ctx context.Context, authUC ports.RegisterUserUseCasesPort, payload KafkaUserCreatePayload) error {

	if payload.Username == "" || payload.Email == "" || payload.Password == "" {
		return domain.ErrInvalidPayload
	}

	user := domain.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	}
	_, err := authUC.Execute(ctx, user)
	return err
}
