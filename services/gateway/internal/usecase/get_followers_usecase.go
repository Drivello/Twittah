package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/gateway/internal/ports"
)

type GetFollowersUseCase struct {
	UserService ports.UserServicePort
}

func NewGetFollowersUseCase(userService ports.UserServicePort) *GetFollowersUseCase {
	return &GetFollowersUseCase{UserService: userService}
}

func (uc *GetFollowersUseCase) Execute(ctx context.Context, userID string) ([]string, error) {
	return uc.UserService.GetFollowers(ctx, userID)
}
