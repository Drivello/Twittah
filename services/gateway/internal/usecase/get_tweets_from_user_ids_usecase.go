package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"github.com/Drivello/Twittah/services/gateway/internal/domain"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
	"go.uber.org/zap"
)

type GetTweetsFromMultipleUserIDsUseCase struct {
	TweetService ports.TweetServicePort
}

func NewGetTweetsFromMultipleUserIDsUseCase(tweetService ports.TweetServicePort) *GetTweetsFromMultipleUserIDsUseCase {
	return &GetTweetsFromMultipleUserIDsUseCase{TweetService: tweetService}
}

func (uc *GetTweetsFromMultipleUserIDsUseCase) Execute(ctx context.Context, userIDs []int64) ([]*domain.Tweet, error) {
	common.Logger().Debug("[Gateway] GetTweetsFromMultipleUserIDsUseCase.Execute called", zap.Any("user_ids", userIDs))
	tweets, err := uc.TweetService.GetTweetsFromMultipleUserIDs(ctx, userIDs)
	if err != nil {
		common.Logger().Debug("[Gateway] GetTweetsFromMultipleUserIDsUseCase.Execute error", zap.Any("user_ids", userIDs), zap.Error(err))
		return nil, err
	}
	common.Logger().Debug("[Gateway] GetTweetsFromMultipleUserIDsUseCase.Execute success", zap.Any("user_ids", userIDs), zap.Any("tweets", tweets))
	return tweets, nil
}
