package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/gateway/internal/domain"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
)

// TweetUseCase implementa la lógica de tweets y cumple con el puerto hexagonal TweetUseCasePort

type TweetUseCase struct {
	Producer ports.TweetEventProducerPort
}

var _ ports.TweetUseCasePort = (*TweetUseCase)(nil)

func (uc *TweetUseCase) PublishTweet(ctx context.Context, event domain.KafkaEventRequest) error {
	return uc.Producer.PublishTweet(ctx, event)
}

type TweetUsecase struct {
	Publisher ports.TweetEventProducerPort
}

func NewTweetUsecase(publisher ports.TweetEventProducerPort) *TweetUsecase {
	return &TweetUsecase{Publisher: publisher}
}

type TweetInput struct {
	AuthorID string
	Content  string
}

func (uc *TweetUsecase) PublishTweet(ctx context.Context, input TweetInput) error {
	event := domain.KafkaEventRequest{
		EventType: "tweet_created",
		Payload: map[string]interface{}{
			"author_id": input.AuthorID,
			"content":   input.Content,
		},
	}
	return uc.Publisher.PublishTweet(ctx, event)
}
