package ports

import (
	"context"
	"github.com/Drivello/Twittah/services/gateway/internal/adapters/kafka"
)

// TweetPublisherPort define el puerto para publicar tweets (ej: a Kafka)
type TweetPublisherPort interface {
	PublishTweet(ctx context.Context, event kafka.TweetPublishedEventDTO) error
}
