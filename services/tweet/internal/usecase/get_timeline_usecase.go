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
}

func NewGetTimelineUseCase(repo ports.TweetRepository, userService ports.UserServicePort) *GetTimelineUseCase {
	return &GetTimelineUseCase{repo: repo, userService: userService}
}

func (uc *GetTimelineUseCase) Execute(ctx context.Context, userID int64) ([]*domain.Tweet, error) {
	common.Logger().Debug("[GetTimelineUseCase] Execute called", zap.Any("user_id", userID))
	if userID <= 0 {
		return nil, fmt.Errorf("user id must be greater than zero")
	}
	following, err := uc.userService.GetFollowing(ctx, userID)
	if err != nil {
		return nil, err
	}
	userIDs := make([]int64, len(following)+1)
	userIDs[0] = userID
	for i, u := range following {
		userIDs[i+1] = u.ID
	}
	return uc.repo.FindAllByMultipleUserIDs(userIDs)
}
