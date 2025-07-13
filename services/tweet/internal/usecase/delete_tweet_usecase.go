package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/tweet/internal/ports"
)

type DeleteTweetUsecase struct {
	repo ports.TweetRepository
}

func NewDeleteTweetUsecase(repo ports.TweetRepository) *DeleteTweetUsecase {
	return &DeleteTweetUsecase{repo: repo}
}

func (uc *DeleteTweetUsecase) DeleteTweet(ctx context.Context, tweetID int64) error {
	return uc.repo.Delete(tweetID)
}
