package kafka

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Drivello/Twittah/services/user/internal/common"
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
		common.Logger().Error("[KafkaEventDispatcher] Invalid event message, skipping", zap.Error(err))
		return &KafkaEventError{
			EventType: req.EventType,
			Error:     err,
		}
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
		return returnEventError(err, req)

	case "users.follow", "users.follow.dlq":
		var followRequest KafkaFollowRequest
		if err := json.Unmarshal(data, &followRequest); err != nil {
			return &KafkaEventError{
				EventType: req.EventType,
				Error:     err,
			}
		}
		err := HandleUserFollow(ctx, d.FollowUC, followRequest.Payload)
		return returnEventError(err, req)
	case "users.unfollow", "users.unfollow.dlq":
		var unfollowRequest KafkaFollowRequest
		if err := json.Unmarshal(data, &unfollowRequest); err != nil {
			return &KafkaEventError{
				EventType: req.EventType,
				Error:     err,
			}
		}
		err := HandleUserUnfollow(ctx, d.UnfollowUC, unfollowRequest.Payload)
		return returnEventError(err, req)
	default:
		common.Logger().Error("[KafkaEventDispatcher] Unknown event type", zap.String("event_type", req.EventType))
		return &KafkaEventError{
			EventType: "unknown.event",
			Error:     errors.New("unknown event type"),
		}
	}

}

func returnEventError(err error, req KafkaEventRequest) *KafkaEventError {
	if err != nil {
		return &KafkaEventError{
			EventType: req.EventType,
			Error:     err,
		}
	}
	return nil
}
