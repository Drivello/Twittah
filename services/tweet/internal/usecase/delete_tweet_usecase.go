package usecase

import (
	"context"
	"fmt"

	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"github.com/Drivello/Twittah/services/tweet/internal/ports"
	"go.uber.org/zap"
)

type DeleteTweetUsecase struct {
	repo ports.TweetRepository
}

func NewDeleteTweetUsecase(repo ports.TweetRepository) *DeleteTweetUsecase {
	return &DeleteTweetUsecase{repo: repo}
}

func (uc *DeleteTweetUsecase) Execute(ctx context.Context, tweetID int64) error {
	common.Logger().Debug("[DeleteTweetUsecase] Execute called", zap.Int64("tweet_id", tweetID))
	if tweetID <= 0 {
		return fmt.Errorf("tweet id must be greater than zero")
	}

	return uc.repo.Delete(tweetID)
}
