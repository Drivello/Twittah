package kafka

import (
	"context"
	"encoding/json"

	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"github.com/Drivello/Twittah/services/gateway/internal/domain"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
	"github.com/IBM/sarama"
)

type UserEventProducer struct {
	producer sarama.SyncProducer
	topic    string
}

// Verifica en compile-time que implementa el puerto hexagonal
var _ ports.UserEventProducerPort = (*UserEventProducer)(nil)

// NewUserEventProducer creates a new UserEventProducer.
// brokers: Kafka broker addresses.
// topic: Kafka topic for user events.
// Returns a pointer to UserEventProducer and error if any.
func NewUserEventProducer(brokers []string, topic string) (*UserEventProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		common.Logger().Errorw("Failed to create follow kafka producer", "error", err)
		return nil, err
	}
	return &UserEventProducer{producer: producer, topic: topic}, nil
}



// PublishFollow publishes a follow event to Kafka.
// followerID: ID of the user following.
// followeeID: ID of the user being followed.
// Returns error if publishing fails.
func (p *UserEventProducer) PublishFollow(ctx context.Context, followerID, followeeID string) error {
	event := domain.UserEvent{
		EventType:  "follow_created",
		FollowerID: followerID,
		FolloweeID: followeeID,
	}
	return p.publishEvent(event)
}

// PublishUnfollow publishes an unfollow event to Kafka.
// followerID: ID of the user unfollowing.
// followeeID: ID of the user being unfollowed.
// Returns error if publishing fails.
func (p *UserEventProducer) PublishUnfollow(ctx context.Context, followerID, followeeID string) error {
	event := domain.UserEvent{
		EventType:  "follow_deleted",
		FollowerID: followerID,
		FolloweeID: followeeID,
	}
	return p.publishEvent(event)
}

func (p *UserEventProducer) publishEvent(event domain.UserEvent) error {
	msgBytes, err := json.Marshal(event)
	if err != nil {
		common.Logger().Errorw("Failed to marshal follow event", "error", err)
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(msgBytes),
	}
	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		common.Logger().Errorw("Failed to send follow event to Kafka", "error", err)
	}
	return err
}
