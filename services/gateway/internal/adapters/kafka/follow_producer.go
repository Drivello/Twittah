package kafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type FollowEventProducer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewFollowEventProducer creates a new FollowEventProducer.
// brokers: Kafka broker addresses.
// topic: Kafka topic for follow events.
// Returns a pointer to FollowEventProducer and error if any.
func NewFollowEventProducer(brokers []string, topic string) (*FollowEventProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &FollowEventProducer{producer: producer, topic: topic}, nil
}

type FollowEvent struct {
	EventType  string `json:"event_type"`
	FollowerID string `json:"follower_id"`
	FolloweeID string `json:"followee_id"`
}

// PublishFollow publishes a follow event to Kafka.
// followerID: ID of the user following.
// followeeID: ID of the user being followed.
// Returns error if publishing fails.
func (p *FollowEventProducer) PublishFollow(followerID, followeeID string) error {
	event := FollowEvent{
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
func (p *FollowEventProducer) PublishUnfollow(followerID, followeeID string) error {
	event := FollowEvent{
		EventType:  "follow_deleted",
		FollowerID: followerID,
		FolloweeID: followeeID,
	}
	return p.publishEvent(event)
}

func (p *FollowEventProducer) publishEvent(event FollowEvent) error {
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
		zap.S().Errorw("Failed to send follow event to Kafka", "error", err)
	}
	return err
}
