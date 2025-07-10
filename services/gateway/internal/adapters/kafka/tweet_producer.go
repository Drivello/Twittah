package kafka

import (
	"context"
	"encoding/json"

	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"github.com/Drivello/Twittah/services/gateway/internal/domain"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
	"github.com/IBM/sarama"
)

// TweetPublishedEventDTO is the DTO for tweet creation events sent to Kafka
// (copiado del adapter de TweetService para consistencia)

type TweetEventProducer struct {
	Producer         sarama.SyncProducer
	createTweetTopic string
}

// Verifica en compile-time que implementa el puerto hexagonal
var _ ports.TweetEventProducerPort = (*TweetEventProducer)(nil)

// NewTweetEventProducer creates a new TweetEventProducer for publishing tweet creation events.
// brokers: Kafka broker addresses.
// topic: Kafka topic for tweet events.
// Returns a pointer to TweetEventProducer.
func NewTweetEventProducer(brokers []string, topic string) (*TweetEventProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	common.Logger().Debug("TweetEventProducer created", "topic ", topic)
	return &TweetEventProducer{Producer: producer, createTweetTopic: topic}, nil
}

// PublishTweet implementa ports.TweetEventProducerPort
func (p *TweetEventProducer) PublishTweet(ctx context.Context, event domain.TweetPublishEventDTO) error {
	payload, err := json.Marshal(event)
	if err != nil {
		common.Logger().Errorw("Failed to marshal tweet event", "error", err)
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: p.createTweetTopic,
		Value: sarama.ByteEncoder(payload),
	}
	_, _, err = p.Producer.SendMessage(msg)
	if err != nil {
		common.Logger().Errorw("Failed to send tweet event to Kafka", "error", err)
	}
	return err
}
