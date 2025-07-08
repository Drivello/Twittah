package kafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// UserCreatedEvent represents the event structure for a created user.
type UserCreatedEvent struct {
	EventType string `json:"event_type"`
	ID        int64  `json:"id"`
	Username  string `json:"username"`
}

// UserEventProducer wraps a sarama.SyncProducer for publishing user events.
type UserEventProducer struct {
	Producer sarama.SyncProducer
	Topic    string
}

// NewUserEventProducer creates a new UserEventProducer for the given brokers and topic.
func NewUserEventProducer(brokers []string, topic string) (*UserEventProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		zap.L().Error("[KafkaProducer] Failed to create SyncProducer", zap.Error(err))
		return nil, err
	}
	zap.L().Info("[KafkaProducer] SyncProducer created successfully", zap.Strings("brokers", brokers), zap.String("topic", topic))
	return &UserEventProducer{Producer: producer, Topic: topic}, nil
}

// PublishUserCreated publishes a user_created event to Kafka.
func (p *UserEventProducer) PublishUserCreated(id int64, username string) error {
	zap.L().Info("[Kafka] Publishing user_created event", zap.Int64("id", id), zap.String("username", username))
	event := UserCreatedEvent{
		EventType: "user_created",
		ID:        id,
		Username:  username,
	}
	value, err := json.Marshal(event)
	if err != nil {
		zap.L().Error("[Kafka] Failed to marshal user_created event", zap.Error(err))
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: p.Topic,
		Value: sarama.ByteEncoder(value),
	}
	partition, offset, err := p.Producer.SendMessage(msg)
	if err != nil {
		zap.L().Error("[Kafka] Failed to publish user_created event", zap.Error(err))
		return err
	}
	zap.L().Info("[Kafka] user_created event published", zap.Int64("id", id), zap.String("username", username), zap.Int32("partition", partition), zap.Int64("offset", offset))
	return nil
}
