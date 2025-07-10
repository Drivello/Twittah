package kafka

import (
	"context"
	"encoding/json"

	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"github.com/Drivello/Twittah/services/gateway/internal/domain"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
	"github.com/IBM/sarama"
)

type UserCreateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthEventProducer publishes user-related events to Kafka.
type AuthEventProducer struct {
	producer      sarama.SyncProducer
	registerTopic string
}

// Verifica en compile-time que implementa el puerto hexagonal
var _ ports.AuthEventProducerPort = (*AuthEventProducer)(nil)

// NewAuthEventProducer creates a new AuthEventProducer.
// brokers: Kafka broker addresses.
// topic: Kafka topic for user events.
// Returns a pointer to AuthEventProducer and error if any.
func NewAuthEventProducer(brokers []string, topic string) (*AuthEventProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	common.Logger().Debug("AuthEventProducer created", "topic ", topic)
	return &AuthEventProducer{producer: producer, registerTopic: topic}, nil
}

// PublishUserCreateRequest publishes a user creation request event to Kafka.
// username: New user's username.
// email: New user's email.
// password: New user's password.
// Returns error if publishing fails.
// Uses domain.UserCreateRequest as the DTO.
func (p *AuthEventProducer) PublishUserCreateRequest(ctx context.Context, username, email, password string) error {
	event := domain.UserCreateRequest{
		Username: username,
		Email:    email,
		Password: password,
	}
	msgBytes, err := json.Marshal(event)
	if err != nil {
		common.Logger().Errorw("Failed to marshal user create request", "error", err)
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: p.registerTopic,
		Value: sarama.ByteEncoder(msgBytes),
	}
	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		common.Logger().Errorw("Failed to send user create request to Kafka", "error", err)
	}
	common.Logger().Debug("User create request sent to Kafka", "username", username)
	return err
}
