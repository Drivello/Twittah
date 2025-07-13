package usecase

import (
	"context"
	"fmt"

	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"github.com/Drivello/Twittah/services/tweet/internal/domain"
	"github.com/Drivello/Twittah/services/tweet/internal/ports"
	"go.uber.org/zap"
)

type GetTweetsFromMultipleUserIDsUsecase struct {
	repo ports.TweetRepository
}

func NewGetTweetsFromMultipleUserIDsUsecase(repo ports.TweetRepository) *GetTweetsFromMultipleUserIDsUsecase {
	return &GetTweetsFromMultipleUserIDsUsecase{repo: repo}
}

func (uc *GetTweetsFromMultipleUserIDsUsecase) Execute(ctx context.Context, userIDs []int64) ([]*domain.Tweet, error) {
	common.Logger().Debug("[GetTweetsFromMultipleUserIDsUsecase] Execute called", zap.Any("user_ids", userIDs))
	if len(userIDs) == 0 {
		return nil, fmt.Errorf("at least one user id must be provided")
	}
	for _, id := range userIDs {
		if id <= 0 {
			return nil, fmt.Errorf("user id must be greater than zero")
		}
	}
	return uc.repo.FindAllByMultipleUserIDs(userIDs)
}
