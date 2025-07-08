package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/auth/internal/domain"
	"github.com/Drivello/Twittah/services/auth/internal/ports"
	"go.uber.org/zap"
)

type userUseCase struct {
	repo ports.UserRepository
}

func NewUserUseCasePort(repo ports.UserRepository) ports.UserUseCasePort {
	return &userUseCase{repo: repo}
}

func (uc *userUseCase) RegisterUser(ctx context.Context, user *domain.User) (int64, error) {
	zap.L().Info("[UseCase] Starting user registration", zap.String("username", user.Username))
	user.Password = hashPassword(user.Password)
	zap.L().Debug("[UseCase] Password hashed")
	userID, err := uc.repo.CreateUser(ctx, user)
	if err != nil {
		zap.L().Error("[UseCase] Failed to create user in repository", zap.Error(err))
		return 0, err
	}
	zap.L().Info("[UseCase] User successfully persisted", zap.Int64("id", userID))
	return userID, nil
}
