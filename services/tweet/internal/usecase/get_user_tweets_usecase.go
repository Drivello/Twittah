package usecase

import (
	"context"
	"fmt"

	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"github.com/Drivello/Twittah/services/tweet/internal/domain"
	"github.com/Drivello/Twittah/services/tweet/internal/ports"
	"go.uber.org/zap"
)

type GetUserTweetsUsecase struct {
	repo ports.TweetRepository
}

func NewGetUserTweetsUsecase(repo ports.TweetRepository) *GetUserTweetsUsecase {
	return &GetUserTweetsUsecase{repo: repo}
}

func (uc *GetUserTweetsUsecase) Execute(ctx context.Context, userID int64) ([]*domain.Tweet, error) {
	common.Logger().Debug("[GetUserTweetsUsecase] Execute called", zap.Int64("user_id", userID))
	if userID <= 0 {
		return nil, fmt.Errorf("user id must be greater than zero")
	}
	return uc.repo.FindAllByUserID(userID)
}
