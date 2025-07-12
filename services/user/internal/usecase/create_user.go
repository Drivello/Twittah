package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/user/internal/domain"
	"github.com/Drivello/Twittah/services/user/internal/ports"
)

// CreateUserUseCase handles the user creation logic.
type CreateUserUseCase struct {
	repo ports.UserRepository
}

// NewCreateUserUseCase constructs a new CreateUserUseCase.
func NewCreateUserUseCase(repo ports.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{repo: repo}
}

// Execute creates a new user.
func (uc *CreateUserUseCase) Execute(ctx context.Context, user *domain.User) error {
	return uc.repo.InsertUser(ctx, user)
}
