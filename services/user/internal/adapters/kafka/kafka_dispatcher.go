package kafka

import (
	"context"
	"encoding/json"

	"github.com/Drivello/Twittah/services/user/internal/common"
	"github.com/Drivello/Twittah/services/user/internal/metrics"
	"github.com/Drivello/Twittah/services/user/internal/ports"
	"go.uber.org/zap"
)

type KafkaEventDispatcher struct {
	UserCreatedUC ports.CreateUserUseCasePort
	FollowUC      ports.FollowUserUseCasePort
	UnfollowUC    ports.UnfollowUserUseCasePort
}

func NewKafkaEventDispatcher(createUserUC ports.CreateUserUseCasePort, followUC ports.FollowUserUseCasePort, unfollowUC ports.UnfollowUserUseCasePort) *KafkaEventDispatcher {
	return &KafkaEventDispatcher{
		UserCreatedUC: createUserUC,
		FollowUC:      followUC,
		UnfollowUC:    unfollowUC,
	}
}

func (d *KafkaEventDispatcher) Dispatch(ctx context.Context, data []byte) *KafkaEventError {

	var req KafkaEventRequest
	if err := json.Unmarshal(data, &req); err != nil {
		common.Logger().Error("[KafkaEventDispatcher] Invalid event message. Skipping.", zap.Error(err))
		return nil
	}
	common.Logger().Debug("[KafkaEventDispatcher] Dispatching event", zap.String("event_type", req.EventType))

	switch req.EventType {
	case "users.created", "users.created.dlq":
		var userCreatedRequest KafkaUserCreatedRequest

		if err := json.Unmarshal(data, &userCreatedRequest); err != nil {
			return &KafkaEventError{
				EventType: req.EventType,
				Error:     err,
			}
		}

		err := HandleUserCreated(ctx, d.UserCreatedUC, userCreatedRequest.Payload)
		if err != nil {
			return &KafkaEventError{
				EventType: req.EventType,
				Error:     err,
			}
		}

		return nil
	case "users.follow", "users.follow.dlq":
		var followRequest KafkaFollowRequest
		metrics.FollowsConsumerTotal.Inc()

		if err := json.Unmarshal(data, &followRequest); err != nil {
			metrics.FollowsErrorsTotal.Inc()
			return &KafkaEventError{
				EventType: req.EventType,
				Error:     err,
			}
		}
		err := HandleUserFollow(ctx, d.FollowUC, followRequest.Payload)
		if err != nil {
			metrics.FollowsErrorsTotal.Inc()
			return &KafkaEventError{
				EventType: req.EventType,
				Error:     err,
			}
		}

		metrics.FollowsCreatedTotal.Inc()
		return nil
	case "users.unfollow", "users.unfollow.dlq":
		var unfollowRequest KafkaFollowRequest
		metrics.UnfollowsConsumerTotal.Inc()

		if err := json.Unmarshal(data, &unfollowRequest); err != nil {
			metrics.UnfollowsErrorsTotal.Inc()
			return &KafkaEventError{
				EventType: req.EventType,
				Error:     err,
			}
		}
		err := HandleUserUnfollow(ctx, d.UnfollowUC, unfollowRequest.Payload)
		if err != nil {
			metrics.UnfollowsErrorsTotal.Inc()
			return &KafkaEventError{
				EventType: req.EventType,
				Error:     err,
			}
		}

		metrics.UnfollowsDoneTotal.Inc()
		return nil
	default:
		common.Logger().Error("[KafkaEventDispatcher] Unknown event type. Skipping.", zap.String("event_type", req.EventType))
		return nil
	}

}
