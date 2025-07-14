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
	repo          ports.TweetRepository
	userService   ports.UserServicePort
	timelineCache ports.TimelineCachePort
}

func NewCreateTweetUsecase(repo ports.TweetRepository, userService ports.UserServicePort, timelineCache ports.TimelineCachePort) *CreateTweetUsecase {
	return &CreateTweetUsecase{repo: repo, userService: userService, timelineCache: timelineCache}
}

func (uc *CreateTweetUsecase) Execute(ctx context.Context, tweet *domain.Tweet) error {
	common.Logger().Debug("[CreateTweetUsecase] Execute called", zap.Any("tweet", tweet))
	if tweet == nil || len(tweet.Content) < 1 {
		return fmt.Errorf("tweet content must not be empty")
	}
	if tweet.AuthorID <= 0 {
		return fmt.Errorf("tweet author id must be greater than zero")
	}
	err := uc.repo.Save(tweet)
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
