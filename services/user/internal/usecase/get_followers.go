package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/user/internal/domain"
	"github.com/Drivello/Twittah/services/user/internal/ports"
	"github.com/Drivello/Twittah/services/user/internal/common"
	"go.uber.org/zap"
)

type GetFollowersUseCase struct {
	repo ports.UserRepository
}

func NewGetFollowersUseCase(repo ports.UserRepository) *GetFollowersUseCase {
	return &GetFollowersUseCase{repo: repo}
}

func (uc *GetFollowersUseCase) Execute(ctx context.Context, userID int64) ([]*domain.User, error) {
	common.Logger().Debug("[GetFollowersUseCase] Execute called", zap.Int64("user_id", userID))
	followers, err := uc.repo.GetFollowers(ctx, userID)
	if err != nil {
		common.Logger().Debug("[GetFollowersUseCase] Error", zap.Int64("user_id", userID), zap.Error(err))
		return nil, err
	}
	common.Logger().Debug("[GetFollowersUseCase] Success", zap.Int64("user_id", userID), zap.Any("followers", followers))
	return followers, nil
}
