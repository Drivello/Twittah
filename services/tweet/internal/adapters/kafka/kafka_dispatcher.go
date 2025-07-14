package kafka

import (
	"context"
	"encoding/json"

	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"github.com/Drivello/Twittah/services/tweet/internal/metrics"
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
	case "tweets.create", "tweets.create.dlq":
		var tweetCreatedRequest KafkaCreateTweetRequest
		metrics.TweetsCreationConsumedTotal.Inc()

		if err := json.Unmarshal(data, &tweetCreatedRequest); err != nil {
			metrics.TweetsCreatedErrorTotal.Inc()
			return &KafkaEventError{
				EventType: req.EventType,
				Error:     err,
			}
		}

		err := HandleTweetCreated(ctx, d.TweetCreatedUC, tweetCreatedRequest.Payload)
		if err != nil {
			metrics.TweetsCreatedErrorTotal.Inc()
			return &KafkaEventError{
				EventType: req.EventType,
				Error:     err,
			}
		}

		metrics.TweetsCreatedTotal.Inc()
		return nil
	case "tweets.delete", "tweets.delete.dlq":
		var tweetDeletedRequest KafkaDeleteTweetRequest
		metrics.TweetsDeletionConsumedTotal.Inc()

		if err := json.Unmarshal(data, &tweetDeletedRequest); err != nil {
			metrics.TweetsDeletedErrorTotal.Inc()
			return &KafkaEventError{
				EventType: req.EventType,
				Error:     err,
			}
		}

		err := HandleTweetDeleted(ctx, d.TweetDeletedUC, tweetDeletedRequest.Payload)
		if err != nil {
			metrics.TweetsDeletedErrorTotal.Inc()
			return &KafkaEventError{
				EventType: req.EventType,
				Error:     err,
			}
		}

		metrics.TweetsDeletedTotal.Inc()
		return nil
	default:
		common.Logger().Error("[KafkaEventDispatcher] Unknown event type. Skipping.", zap.String("event_type", req.EventType))
		return nil
	}

}
