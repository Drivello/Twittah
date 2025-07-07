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

func (p *FollowEventProducer) PublishFollow(followerID, followeeID string) error {
	event := FollowEvent{
		EventType:  "follow_created",
		FollowerID: followerID,
		FolloweeID: followeeID,
	}
	return p.publishEvent(event)
}

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
