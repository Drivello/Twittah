package usecase

import (
	"context"
	"fmt"

	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"github.com/Drivello/Twittah/services/tweet/internal/domain"
	"github.com/Drivello/Twittah/services/tweet/internal/ports"
	"go.uber.org/zap"
)

type CreateUserUsecase struct {
	repo ports.UserRepository
}

func NewCreateUserUsecase(repo ports.UserRepository) *CreateUserUsecase {
	return &CreateUserUsecase{repo: repo}
}

func (uc *CreateUserUsecase) Execute(ctx context.Context, user *domain.User) error {
	common.Logger().Debug("[CreateUserUsecase] Execute called", zap.Any("user", user))
	if user == nil {
		return fmt.Errorf("user is nil")
	}
	if user.ID <= 0 {
		return fmt.Errorf("user id must be greater than zero")
	}
	if user.Username == "" {
		return fmt.Errorf("user username must not be empty")
	}
	return uc.repo.Save(user)
}
