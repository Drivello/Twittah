package usecase

import (
	"context"
	"fmt"

	"github.com/Drivello/Twittah/services/gateway/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
	"go.uber.org/zap"
)

// CreateTweetUseCase implements the user logic and fulfills the hexagonal port CreateTweetUseCasePort.

type CreateTweetUseCase struct {
	Producer ports.EventProducerPort[kafka.TweetPayload]
}

var _ ports.CreateTweetUseCasePort = (*CreateTweetUseCase)(nil)

// NewCreateTweetUseCase creates a new CreateTweetUseCase.
func NewCreateTweetUseCase(producer ports.EventProducerPort[kafka.TweetPayload]) *CreateTweetUseCase {
	return &CreateTweetUseCase{Producer: producer}
}

// Execute handles unfollow logic.
func (uc *CreateTweetUseCase) Execute(ctx context.Context, authorID int64, content string) error {
	if len(content) < 1 {
		return fmt.Errorf("tweet content must not be empty")
	}
	request := kafka.KafkaTweetCreatePayload{
		AuthorID: authorID,
		Content:  content,
	}
	common.Logger().Debug("[Gateway] CreateTweetUseCase.Execute called", zap.Int64("author_id", authorID), zap.String("content", content))
	return uc.Producer.PublishEvent(kafka.KafkaEventRequest[kafka.TweetPayload]{
		EventType: "tweets.create",
		Payload:   request,
	})
}
