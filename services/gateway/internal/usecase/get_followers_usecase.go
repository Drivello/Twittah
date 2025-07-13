package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"github.com/Drivello/Twittah/services/gateway/internal/domain"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
	"go.uber.org/zap"
)

type GetFollowersUseCase struct {
	UserService ports.UserServicePort
}

func NewGetFollowersUseCase(userService ports.UserServicePort) *GetFollowersUseCase {
	return &GetFollowersUseCase{UserService: userService}
}

func (uc *GetFollowersUseCase) Execute(ctx context.Context, userID int64) ([]*domain.User, error) {
	common.Logger().Debug("[Gateway] GetFollowersUseCase.Execute called", zap.Int64("user_id", userID))
	followers, err := uc.UserService.GetFollowers(ctx, userID)
	if err != nil {
		common.Logger().Debug("[Gateway] GetFollowersUseCase.Execute error", zap.Int64("user_id", userID), zap.Error(err))
		return nil, err
	}
	common.Logger().Debug("[Gateway] GetFollowersUseCase.Execute success", zap.Int64("user_id", userID), zap.Any("followers", followers))
	return followers, nil
}
