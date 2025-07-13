package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/tweet/internal/domain"
	"github.com/Drivello/Twittah/services/tweet/internal/ports"
)

type CreateUserUsecase struct {
	repo ports.UserRepository
}

func NewCreateUserUsecase(repo ports.UserRepository) *CreateUserUsecase {
	return &CreateUserUsecase{repo: repo}
}

func (uc *CreateUserUsecase) Execute(ctx context.Context, user *domain.User) error {
	return uc.repo.Save(user)
}
