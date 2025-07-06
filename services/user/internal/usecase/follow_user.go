package usecase

import (
	"context"
	"errors"

	"github.com/Drivello/Twittah/services/user/internal/ports"
)

type FollowUserUseCase struct {
	Repo ports.UserRepository
}

func NewFollowUserUseCase(repo ports.UserRepository) *FollowUserUseCase {
	return &FollowUserUseCase{Repo: repo}
}

func (uc *FollowUserUseCase) Execute(ctx context.Context, followerID, followeeID string) error {
	if followerID == followeeID {
		return errors.New("cannot follow yourself")
	}
	return uc.Repo.FollowUser(ctx, followerID, followeeID)
}
