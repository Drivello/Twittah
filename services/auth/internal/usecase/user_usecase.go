package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/auth/internal/common"
	"github.com/Drivello/Twittah/services/auth/internal/domain"
	"github.com/Drivello/Twittah/services/auth/internal/ports"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type authUseCase struct {
	repo                     ports.AuthRepositoryPort
	producer                 ports.UserEventProducerPort
	authUserCreatedEventType string
}

func NewAuthUseCasesPort(repo ports.AuthRepositoryPort, producer ports.UserEventProducerPort, eventString string) ports.AuthUseCasesPort {
	common.Logger().Debug("[AuthUseCase] Creating new auth use case", zap.String("event_string", eventString))
	return &authUseCase{repo: repo, producer: producer, authUserCreatedEventType: eventString}
}

func (uc *authUseCase) RegisterUser(ctx context.Context, user domain.User) (int64, error) {
	common.Logger().Info("[AuthUseCase] Starting user registration", zap.String("username", user.Username))

	// Hash de la contraseña
	hashed, err := hashPassword(user.Password)
	if err != nil {
		common.Logger().Error("[AuthUseCase] Failed to hash password", zap.Error(err))
		return 0, err
	}
	user.Password = hashed

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
		pubErr := uc.producer.PublishUserCreated(userID, user.Username, uc.authUserCreatedEventType)
		if pubErr != nil {
			common.Logger().Error("[AuthUseCase] Failed to publish user created event", zap.Error(pubErr))
		} else {
			common.Logger().Info("[AuthUseCase] Published user created event", zap.Int64("id", userID))
		}
	}

	return userID, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
