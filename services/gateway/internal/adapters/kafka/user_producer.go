package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type UserCreateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserEventProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func (p *UserEventProducer) PublishUserCreateRequest(username, email, password string) error {
	event := UserCreateRequest{
		Username: username,
		Email:    email,
		Password: password,
	}
	msgBytes, err := json.Marshal(event)
	if err != nil {
		zap.S().Errorw("Failed to marshal user create request", "error", err)
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(msgBytes),
	}
	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		zap.S().Errorw("Failed to send user create request to Kafka", "error", err)
	}
	return err
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

// TODO: Implement user_created event
func (p *UserEventProducer) PublishUserCreated(id, username string) error {
	// TODO: Marshal and send user_created event to Kafka
	return nil
}

// TODO: Implement user_deleted event
func (p *UserEventProducer) PublishUserDeleted(id string) error {
	// TODO: Marshal and send user_deleted event to Kafka
	return nil
}
