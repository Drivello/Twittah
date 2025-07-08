package kafka

import (
	"encoding/json"

	"github.com/Drivello/Twittah/services/user/internal/domain"
	"github.com/Drivello/Twittah/services/user/internal/ports"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type UserEventProducer struct {
	producer sarama.SyncProducer
	topic    string
}

var _ ports.UserEventPublisher = (*UserEventProducer)(nil)

func NewUserEventProducer(brokers []string, topic string) (*UserEventProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &UserEventProducer{producer: producer, topic: topic}, nil
}

func (p *UserEventProducer) PublishFollow(followerID, followeeID string) error {
	event := domain.FollowEvent{
		EventType:  "follow_created",
		FollowerID: followerID,
		FolloweeID: followeeID,
	}
	return p.publishEvent(event)
}

func (p *UserEventProducer) PublishUnfollow(followerID, followeeID string) error {
	event := domain.FollowEvent{
		EventType:  "follow_deleted",
		FollowerID: followerID,
		FolloweeID: followeeID,
	}
	return p.publishEvent(event)
}

func (p *UserEventProducer) publishEvent(event domain.FollowEvent) error {
	msgBytes, err := json.Marshal(event)
	if err != nil {
		zap.S().Errorw("Failed to marshal follow event", "error", err)
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(msgBytes),
	}
	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		zap.S().Errorw("Failed to send Kafka event", "error", err)
	}
	return err
}
