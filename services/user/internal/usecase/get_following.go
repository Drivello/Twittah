package usecase

import (
	"context"
	"github.com/Drivello/Twittah/services/user/internal/ports"
)

type GetFollowingUseCase struct {
	repo ports.UserRepository
}

func NewGetFollowingUseCase(repo ports.UserRepository) *GetFollowingUseCase {
	return &GetFollowingUseCase{repo: repo}
}

func (uc *GetFollowingUseCase) Execute(ctx context.Context, userID int64) ([]int64, error) {
	return uc.repo.GetFollowing(ctx, userID)
}
