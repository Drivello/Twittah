package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/gateway/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
)

// UserUseCase implementa la lógica de usuario y cumple con el puerto hexagonal UserUseCasePort

type CreateTweetUseCase struct {
	Producer ports.EventProducerPort[kafka.TweetPayload]
}

var _ ports.CreateTweetUseCasePort = (*CreateTweetUseCase)(nil)

// NewTweetUseCase crea un nuevo UserUseCase.
func NewCreateTweetUseCase(producer ports.EventProducerPort[kafka.TweetPayload]) *CreateTweetUseCase {
	return &CreateTweetUseCase{Producer: producer}
}

// Execute handles unfollow logic.
func (uc *CreateTweetUseCase) Execute(ctx context.Context, authorID int64, content string) error {
	// request := kafka.KafkaFollowRequest{
	// 	EventType: "tweet.create",
	// 	Payload: kafka.KafkaFollowPayload{
	// 		FollowerID: followerID,
	// 		FolloweeID: followeeID,
	// 	},
	// }
	// requestBytes, err := json.Marshal(request)
	// if err != nil {
	// 	return err
	// }
	// return uc.Producer.PublishEventRequest(ctx, requestBytes)
	return nil
}
