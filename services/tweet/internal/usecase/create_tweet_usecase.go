package usecase

import (
	"context"
	"fmt"

	"github.com/Drivello/Twittah/services/tweet/internal/domain"
	"github.com/Drivello/Twittah/services/tweet/internal/ports"
	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"go.uber.org/zap"
)

type CreateTweetUsecase struct {
	repo ports.TweetRepository
}

func NewCreateTweetUsecase(repo ports.TweetRepository) *CreateTweetUsecase {
	return &CreateTweetUsecase{repo: repo}
}

func (uc *CreateTweetUsecase) CreateTweet(ctx context.Context, tweet *domain.Tweet) error {
	common.Logger().Debug("[CreateTweetUsecase] CreateTweet called", zap.Any("tweet", tweet))
	if tweet == nil || len(tweet.Content) < 1 {
		return fmt.Errorf("tweet content must not be empty")
	}
	return uc.repo.Save(tweet)
}
