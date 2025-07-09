package kafka

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

// TweetPublishedEventDTO is the DTO for tweet creation events sent to Kafka
// (copiado del adapter de TweetService para consistencia)
type TweetPublishedEventDTO struct {
	ID        string `json:"id"`
	AuthorID  string `json:"author_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// TweetProducer implementa ports.TweetPublisherPort
// topic debe coincidir con el que consume TweetService
// Uso: NewTweetProducer(sarama.SyncProducer, "tweets.published")
type TweetProducer struct {
	Producer sarama.SyncProducer
	Topic    string
}

func NewTweetProducer(producer sarama.SyncProducer, topic string) *TweetProducer {
	return &TweetProducer{Producer: producer, Topic: topic}
}

// PublishTweet implementa ports.TweetPublisherPort
func (p *TweetProducer) PublishTweet(ctx context.Context, event TweetPublishedEventDTO) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: p.Topic,
		Value: sarama.ByteEncoder(payload),
	}
	_, _, err = p.Producer.SendMessage(msg)
	return err
}
