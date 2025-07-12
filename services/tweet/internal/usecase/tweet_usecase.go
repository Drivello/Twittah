package usecase

import (
	"github.com/Drivello/Twittah/services/tweet/internal/domain"
	"github.com/Drivello/Twittah/services/tweet/internal/ports"
)

type TweetUsecase struct {
	repo ports.TweetRepository
}

func NewTweetUsecase(repo ports.TweetRepository) *TweetUsecase {
	return &TweetUsecase{repo: repo}
}

func (uc *TweetUsecase) CreateTweet(tweet *domain.Tweet) error {
	return uc.repo.Save(tweet)
}

func (uc *TweetUsecase) DeleteTweet(tweetID int64) error {
	return uc.repo.Delete(tweetID)
}
