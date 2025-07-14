package usecase

import (
	"context"
	"fmt"

	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"github.com/Drivello/Twittah/services/tweet/internal/domain"
	"github.com/Drivello/Twittah/services/tweet/internal/ports"
	"go.uber.org/zap"
)

type GetTimelineUseCase struct {
	repo        ports.TweetRepository
	userService ports.UserServicePort
	timelineCache ports.TimelineCachePort
}

func NewGetTimelineUseCase(repo ports.TweetRepository, userService ports.UserServicePort, timelineCache ports.TimelineCachePort) *GetTimelineUseCase {
	return &GetTimelineUseCase{repo: repo, userService: userService, timelineCache: timelineCache}
}

func (uc *GetTimelineUseCase) Execute(ctx context.Context, userID int64) ([]*domain.Tweet, error) {
	common.Logger().Debug("[GetTimelineUseCase] Execute called", zap.Any("user_id", userID))
	if userID <= 0 {
		return nil, fmt.Errorf("user id must be greater than zero")
	}
	// Intentar obtener del cache
	cached, err := uc.timelineCache.GetTimeline(ctx, fmt.Sprintf("%d", userID))
	if err == nil && cached != nil {
		common.Logger().Infof("Cache HIT for user %d", userID)
		// Convertir []domain.Tweet a []*domain.Tweet
		result := make([]*domain.Tweet, len(cached))
		for i := range cached {
			result[i] = &cached[i]
		}
		return result, nil
	} else if err != nil {
		common.Logger().Warnf("Error al consultar cache timeline para user %d: %v", userID, err)
	}
	common.Logger().Infof("Cache MISS for user %d", userID)
	following, err := uc.userService.GetFollowing(ctx, userID)
	if err != nil {
		return nil, err
	}
	userIDs := make([]int64, len(following)+1)
	userIDs[0] = userID
	for i, u := range following {
		if i+1 >= len(userIDs) {
			common.Logger().Errorf("Index out of range: i+1=%d, len(userIDs)=%d, following=%v", i+1, len(userIDs), following)
			break
		}
		userIDs[i+1] = u.ID
	}
	tweets, err := uc.repo.FindAllByMultipleUserIDs(userIDs)
	if err != nil {
		return nil, err
	}
	// Cachear resultado
	if tweets != nil {
		deref := make([]domain.Tweet, len(tweets))
		for i, t := range tweets {
			deref[i] = *t
		}
		errCache := uc.timelineCache.SetTimeline(ctx, fmt.Sprintf("%d", userID), deref)
		if errCache != nil {
			common.Logger().Warnf("Error cacheando timeline para user %d: %v", userID, errCache)
		}
	}
	return tweets, nil
}
