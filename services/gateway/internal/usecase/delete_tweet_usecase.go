package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/gateway/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
	"go.uber.org/zap"
)

type DeleteTweetUseCase struct {
	Producer ports.EventProducerPort[kafka.TweetPayload]
}

var _ ports.DeleteTweetUseCasePort = (*DeleteTweetUseCase)(nil)

func NewDeleteTweetUseCase(producer ports.EventProducerPort[kafka.TweetPayload]) *DeleteTweetUseCase {
	return &DeleteTweetUseCase{Producer: producer}
}

func (uc *DeleteTweetUseCase) Execute(ctx context.Context, tweetID int64) error {
	request := kafka.KafkaTweetDeletePayload{
		TweetID: tweetID,
	}
	common.Logger().Debug("[Gateway] DeleteTweetUseCase.Execute called", zap.Int64("tweet_id", tweetID))
	return uc.Producer.PublishEvent(kafka.KafkaEventRequest[kafka.TweetPayload]{
		EventType: "tweets.delete",
		Payload:   request,
	})
}
