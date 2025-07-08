package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type UserCreatedEvent struct {
	EventType string `json:"event_type"`
	ID        string `json:"id"`
	Username  string `json:"username"`
}

type UserEventProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewUserEventProducer(brokers []string, topic string) (*UserEventProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &UserEventProducer{producer: producer, topic: topic}, nil
}

func (p *UserEventProducer) PublishUserCreated(id, username string) error {
	event := UserCreatedEvent{
		EventType: "user_created",
		ID:        id,
		Username:  username,
	}
	msgBytes, err := json.Marshal(event)
	if err != nil {
		zap.S().Errorw("Failed to marshal user_created event", "error", err)
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(msgBytes),
	}
	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		zap.S().Errorw("Failed to send user_created event", "error", err)
	}
	return err
}
