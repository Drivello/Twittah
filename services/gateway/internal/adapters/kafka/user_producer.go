package kafka

import (
	// "encoding/json"
	"github.com/IBM/sarama"
)

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
