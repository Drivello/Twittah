package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/auth/internal/domain"
	"github.com/Drivello/Twittah/services/auth/internal/ports"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	repo                 ports.UserRepository
	producer             ports.UserEventProducerPort
	userCreatedEventType string
}

func NewUserUseCasePort(repo ports.UserRepository, producer ports.UserEventProducerPort, eventString string) ports.UserUseCasePort {
	return &userUseCase{repo: repo, producer: producer, userCreatedEventType: eventString}
}

func (uc *userUseCase) RegisterUser(ctx context.Context, user *domain.User) (int64, error) {
	zap.L().Info("[UseCase] Starting user registration", zap.String("username", user.Username))
	hashed, err := hashPassword(user.Password)
	if err != nil {
		zap.L().Error("[UseCase] Failed to hash password", zap.Error(err))
		return 0, err
	}
	user.Password = hashed
	zap.L().Debug("[UseCase] Password hashed")
	userID, err := uc.repo.CreateUser(ctx, user)
	if err != nil {
		zap.L().Error("[UseCase] Failed to create user in repository", zap.Error(err))
		return 0, err
	}
	zap.L().Info("[UseCase] User successfully persisted", zap.Int64("id", userID))

	// Publicar evento solo si el usuario fue creado exitosamente
	if uc.producer != nil {
		err := uc.producer.PublishUserCreated(userID, user.Username, uc.userCreatedEventType)
		if err != nil {
			zap.L().Error("[UseCase] Failed to publish user created event", zap.Error(err))
		} else {
			zap.L().Debug("[UseCase] Published user created event", zap.Int64("id", userID))
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
