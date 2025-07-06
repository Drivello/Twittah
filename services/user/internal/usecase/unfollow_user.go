package usecase

import (
	"context"
	"errors"

	"github.com/Drivello/Twittah/services/user/internal/ports"
)

type UnfollowUserUseCase struct {
	Repo ports.UserRepository
}

func NewUnfollowUserUseCase(repo ports.UserRepository) *UnfollowUserUseCase {
	return &UnfollowUserUseCase{Repo: repo}
}

func (uc *UnfollowUserUseCase) Execute(ctx context.Context, followerID, followeeID string) error {
	if followerID == followeeID {
		return errors.New("cannot unfollow yourself")
	}
	return uc.Repo.UnfollowUser(ctx, followerID, followeeID)
}
