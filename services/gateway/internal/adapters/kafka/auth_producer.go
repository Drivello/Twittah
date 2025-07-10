package kafka

import (
	"context"
	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
	"github.com/IBM/sarama"
)

// UserCreateRequest representa el evento de creación de usuario
type UserCreateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthProducer publica eventos relacionados con usuarios a Kafka
type AuthProducer struct {
	producer      sarama.SyncProducer
	registerTopic string
}

// Verifica en compile-time que implementa el puerto hexagonal
var _ ports.AuthEventProducerPort = (*AuthProducer)(nil)

// NewAuthProducer crea un nuevo AuthProducer
func NewAuthProducer(brokers []string, topic string) (*AuthProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		common.Logger().Error("Failed to create Kafka SyncProducer", "error", err)
		return nil, err
	}
	return &AuthProducer{
		producer:      producer,
		registerTopic: topic,
	}, nil
}

// PublishUserCreateRequest publishes a user creation event (already constructed) to Kafka
func (p *AuthProducer) PublishUserCreateRequest(ctx context.Context, eventBytes []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: p.registerTopic,
		Value: sarama.ByteEncoder(eventBytes),
	}

	_, _, err := p.producer.SendMessage(msg)
	if err != nil {
		common.Logger().Error("Failed to send user create request to Kafka", "error", err)
		return err
	}

	common.Logger().Debug("User create request sent to Kafka", "topic", p.registerTopic)
	return nil
}
