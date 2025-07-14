package kafka

import (
	"context"

	"github.com/Drivello/Twittah/services/user/internal/common"
	"github.com/Drivello/Twittah/services/user/internal/domain"
	"github.com/Drivello/Twittah/services/user/internal/ports"
	"go.uber.org/zap"
)

func HandleUserCreated(ctx context.Context, createUserUC ports.CreateUserUseCasePort, payload KafkaUserCreatedPayload) error {
	common.Logger().Debug("[UserEventConsumer] User created event received", zap.Int64("user_id", payload.UserID))
	user := &domain.User{
		ID:       payload.UserID,
		Username: payload.Username,
	}
	return createUserUC.Execute(ctx, user)
}

// HandleUserFollow procesa un evento users.follow
func HandleUserFollow(ctx context.Context, followUC ports.FollowUserUseCasePort, payload KafkaFollowPayload) error {
	common.Logger().Debug("[UserEventConsumer] User follow event received", zap.Int64("follower_id", payload.FollowerID), zap.Int64("followee_id", payload.FolloweeID))
	followerID := payload.FollowerID
	followeeID := payload.FolloweeID
	return followUC.Execute(ctx, followerID, followeeID)
}

// HandleUserUnfollow procesa un evento users.unfollow
func HandleUserUnfollow(ctx context.Context, unfollowUC ports.UnfollowUserUseCasePort, payload KafkaFollowPayload) error {
	common.Logger().Debug("[UserEventConsumer] User unfollow event received", zap.Int64("follower_id", payload.FollowerID), zap.Int64("followee_id", payload.FolloweeID))
	followerID := payload.FollowerID
	followeeID := payload.FolloweeID
	return unfollowUC.Execute(ctx, followerID, followeeID)
}
