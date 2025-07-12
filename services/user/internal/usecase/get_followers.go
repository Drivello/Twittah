package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/user/internal/domain"
	"github.com/Drivello/Twittah/services/user/internal/ports"
)

type GetFollowersUseCase struct {
	repo ports.UserRepository
}

func NewGetFollowersUseCase(repo ports.UserRepository) *GetFollowersUseCase {
	return &GetFollowersUseCase{repo: repo}
}

func (uc *GetFollowersUseCase) Execute(ctx context.Context, userID int64) ([]*domain.User, error) {
	return uc.repo.GetFollowers(ctx, userID)
}
