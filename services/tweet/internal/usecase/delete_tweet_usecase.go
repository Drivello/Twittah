package usecase

import (
	"context"
	"fmt"

	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"github.com/Drivello/Twittah/services/tweet/internal/ports"
	"go.uber.org/zap"
)

type DeleteTweetUsecase struct {
	repo          ports.TweetRepository
	userService   ports.UserServicePort
	timelineCache ports.TimelineCachePort
}

func NewDeleteTweetUsecase(repo ports.TweetRepository, userService ports.UserServicePort, timelineCache ports.TimelineCachePort) *DeleteTweetUsecase {
	return &DeleteTweetUsecase{repo: repo, userService: userService, timelineCache: timelineCache}
}

func (uc *DeleteTweetUsecase) Execute(ctx context.Context, tweetID int64) error {
	common.Logger().Debug("[DeleteTweetUsecase] Execute called", zap.Int64("tweet_id", tweetID))
	if tweetID <= 0 {
		return fmt.Errorf("tweet id must be greater than zero")
	}

	tweet, err := uc.repo.FindByID(tweetID)
	if err != nil {
		return err
	}

	err = uc.repo.Delete(tweetID)

	if err != nil {
		return err
	}

	followers, err := uc.userService.GetFollowers(ctx, tweet.AuthorID)
	if err != nil {
		common.Logger().Warnf("No se pudo obtener seguidores para invalidar cache: %v", err)
	} else {
		for _, follower := range followers {
			_ = uc.timelineCache.InvalidateTimeline(ctx, fmt.Sprintf("%d", follower.ID))
		}
	}
	_ = uc.timelineCache.InvalidateTimeline(ctx, fmt.Sprintf("%d", tweet.AuthorID))
	return nil
}
