package kafka

import (
	"context"

	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"github.com/Drivello/Twittah/services/tweet/internal/domain"
	"github.com/Drivello/Twittah/services/tweet/internal/ports"
	"go.uber.org/zap"
)

func HandleUserCreated(ctx context.Context, createUserUC ports.CreateUserUseCasePort, payload KafkaUserCreatedPayload) error {
	common.Logger().Debug("[UserEventConsumer] User created event received", zap.Int64("user_id", payload.UserID))
	user := &domain.User{
		ID:       payload.UserID,
		Username: payload.Username,
	}
	return createUserUC.Execute(ctx, user)
}
