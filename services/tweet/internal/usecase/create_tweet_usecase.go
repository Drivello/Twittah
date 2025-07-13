package usecase

import (
	"context"
	"fmt"

	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"github.com/Drivello/Twittah/services/tweet/internal/domain"
	"github.com/Drivello/Twittah/services/tweet/internal/ports"
	"go.uber.org/zap"
)

type CreateTweetUsecase struct {
	repo ports.TweetRepository
}

func NewCreateTweetUsecase(repo ports.TweetRepository) *CreateTweetUsecase {
	return &CreateTweetUsecase{repo: repo}
}

func (uc *CreateTweetUsecase) Execute(ctx context.Context, tweet *domain.Tweet) error {
	common.Logger().Debug("[CreateTweetUsecase] Execute called", zap.Any("tweet", tweet))
	if tweet == nil || len(tweet.Content) < 1 {
		return fmt.Errorf("tweet content must not be empty")
	}
	if tweet.AuthorID <= 0 {
		return fmt.Errorf("tweet author id must be greater than zero")
	}
	return uc.repo.Save(tweet)
}
