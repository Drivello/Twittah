package kafka

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"github.com/Drivello/Twittah/services/tweet/internal/ports"
	"go.uber.org/zap"
)

type KafkaEventDispatcher struct {
	UserCreatedUC  ports.CreateUserUseCasePort
	TweetCreatedUC ports.CreateTweetUsecasePort
	TweetDeletedUC ports.DeleteTweetUsecasePort
}

func NewKafkaEventDispatcher(
	createUserUC ports.CreateUserUseCasePort,
	tweetCreatedUC ports.CreateTweetUsecasePort,
	tweetDeletedUC ports.DeleteTweetUsecasePort,
) *KafkaEventDispatcher {
	return &KafkaEventDispatcher{
		UserCreatedUC:  createUserUC,
		TweetCreatedUC: tweetCreatedUC,
		TweetDeletedUC: tweetDeletedUC,
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
	case "tweets.create", "tweets.create.dlq":
		var tweetCreatedRequest KafkaCreateTweetRequest
		if err := json.Unmarshal(data, &tweetCreatedRequest); err != nil {
			return &KafkaEventError{
				EventType: req.EventType,
				Error:     err,
			}
		}
		err := HandleTweetCreated(ctx, d.TweetCreatedUC, tweetCreatedRequest.Payload)
		return returnEventError(err, req)
	case "tweets.delete", "tweets.delete.dlq":
		var tweetDeletedRequest KafkaDeleteTweetRequest
		if err := json.Unmarshal(data, &tweetDeletedRequest); err != nil {
			return &KafkaEventError{
				EventType: req.EventType,
				Error:     err,
			}
		}
		err := HandleTweetDeleted(ctx, d.TweetDeletedUC, tweetDeletedRequest.Payload)
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
