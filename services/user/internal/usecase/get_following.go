package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/user/internal/domain"
	"github.com/Drivello/Twittah/services/user/internal/ports"
	"github.com/Drivello/Twittah/services/user/internal/common"
	"go.uber.org/zap"
)

type GetFollowingUseCase struct {
	repo ports.UserRepository
}

func NewGetFollowingUseCase(repo ports.UserRepository) *GetFollowingUseCase {
	return &GetFollowingUseCase{repo: repo}
}

func (uc *GetFollowingUseCase) Execute(ctx context.Context, userID int64) ([]*domain.User, error) {
	common.Logger().Debug("[GetFollowingUseCase] Execute called", zap.Int64("user_id", userID))
	following, err := uc.repo.GetFollowing(ctx, userID)
	if err != nil {
		common.Logger().Debug("[GetFollowingUseCase] Error", zap.Int64("user_id", userID), zap.Error(err))
		return nil, err
	}
	common.Logger().Debug("[GetFollowingUseCase] Success", zap.Int64("user_id", userID), zap.Any("following", following))
	return following, nil
}
