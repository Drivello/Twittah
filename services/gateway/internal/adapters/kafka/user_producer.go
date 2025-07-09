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

// UserEventProducer publishes user-related events to Kafka.
type UserEventProducer struct {
	producer sarama.SyncProducer
	topic    string
}

// PublishUserCreateRequest publishes a user creation request event to Kafka.
// username: New user's username.
// email: New user's email.
// password: New user's password.
// Returns error if publishing fails.
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

// NewUserEventProducer creates a new UserEventProducer.
// brokers: Kafka broker addresses.
// topic: Kafka topic for user events.
// Returns a pointer to UserEventProducer and error if any.
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
// PublishUserCreated publishes a user_created event to Kafka (TODO).
// id: User ID.
// username: Username.
// Returns error if publishing fails.
func (p *UserEventProducer) PublishUserCreated(id, username string) error {
	// TODO: Marshal and send user_created event to Kafka
	return nil
}

// TODO: Implement user_deleted event
// PublishUserDeleted publishes a user_deleted event to Kafka (TODO).
// id: User ID.
// Returns error if publishing fails.
func (p *UserEventProducer) PublishUserDeleted(id string) error {
	// TODO: Marshal and send user_deleted event to Kafka
	return nil
}
