package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/auth/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/auth/internal/common"
	"github.com/Drivello/Twittah/services/auth/internal/domain"
	"github.com/Drivello/Twittah/services/auth/internal/ports"
	"go.uber.org/zap"
)

type registerUserUseCase struct {
	repo     ports.AuthRepositoryPort
	producer ports.EventProducerPort
}

const (
	UserCreatedEventType = "users.created"
)

func NewRegisterUserUseCase(repo ports.AuthRepositoryPort, producer ports.EventProducerPort) ports.RegisterUserUseCasesPort {
	common.Logger().Debug("[AuthUseCase] Creating new auth use case")
	return &registerUserUseCase{repo: repo, producer: producer}
}

func (uc *registerUserUseCase) Execute(ctx context.Context, user domain.User) (int64, error) {
	common.Logger().Info("[AuthUseCase] Starting user registration", zap.String("username", user.Username))

	// Crear usuario o recuperar existente
	userID, created, err := uc.repo.CreateOrGetUser(ctx, user)
	if err != nil {
		common.Logger().Error("[AuthUseCase] Failed to create user in repository", zap.Error(err))
		return 0, err
	}

	if created {
		common.Logger().Info("[AuthUseCase] User successfully persisted", zap.Int64("id", userID))
	} else {
		common.Logger().Debug("[AuthUseCase] User already existed, skipping creation", zap.Int64("id", userID))
	}

	if created && uc.producer != nil {

		payload := kafka.NewUserCreatedPayload(userID, user.Username)
		pubErr := uc.producer.PublishEvent(UserCreatedEventType, payload)
		if pubErr != nil {
			common.Logger().Error("[AuthUseCase] Failed to publish user created event", zap.Error(pubErr))
		} else {
			common.Logger().Info("[AuthUseCase] Published user created event", zap.Int64("id", userID))
		}
	}

	return userID, nil
}
