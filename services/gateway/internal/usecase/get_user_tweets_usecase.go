package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"github.com/Drivello/Twittah/services/gateway/internal/domain"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
	"go.uber.org/zap"
)

type GetUserTweetsUseCase struct {
	TweetService ports.TweetServicePort
}

func NewGetUserTweetsUseCase(tweetService ports.TweetServicePort) *GetUserTweetsUseCase {
	return &GetUserTweetsUseCase{TweetService: tweetService}
}

func (uc *GetUserTweetsUseCase) Execute(ctx context.Context, userID int64) ([]*domain.Tweet, error) {
	common.Logger().Debug("[Gateway] GetUserTweetsUseCase.Execute called", zap.Int64("user_id", userID))
	tweets, err := uc.TweetService.GetTweetsFromUserID(ctx, userID)
	if err != nil {
		common.Logger().Debug("[Gateway] GetUserTweetsUseCase.Execute error", zap.Int64("user_id", userID), zap.Error(err))
		return nil, err
	}
	common.Logger().Debug("[Gateway] GetUserTweetsUseCase.Execute success", zap.Int64("user_id", userID), zap.Any("tweets", tweets))
	return tweets, nil
}
