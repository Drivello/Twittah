package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"github.com/Drivello/Twittah/services/gateway/internal/domain"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
	"go.uber.org/zap"
)

type GetTimelineUseCase struct {
	TweetService ports.TweetServicePort
}

func NewGetTimelineUseCase(tweetService ports.TweetServicePort) *GetTimelineUseCase {
	return &GetTimelineUseCase{TweetService: tweetService}
}

func (uc *GetTimelineUseCase) Execute(ctx context.Context, userID int64) ([]*domain.Tweet, error) {
	common.Logger().Debug("[Gateway] GetTimelineUseCase.Execute called", zap.Int64("user_id", userID))
	tweets, err := uc.TweetService.GetTimeline(ctx, userID)
	if err != nil {
		common.Logger().Debug("[Gateway] GetTimelineUseCase.Execute error", zap.Int64("user_id", userID), zap.Error(err))
		return nil, err
	}
	common.Logger().Debug("[Gateway] GetTimelineUseCase.Execute success", zap.Int64("user_id", userID), zap.Any("tweets", tweets))
	return tweets, nil
}
