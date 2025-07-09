// Package kafka provides Kafka producers for TweetService.
package kafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// TimelineProducer produces timelines.updated events to Kafka.
type TimelineProducer struct {
	producer sarama.SyncProducer
	logger   *zap.Logger
}

// NewTimelineProducer creates a new TimelineProducer.
func NewTimelineProducer(producer sarama.SyncProducer, logger *zap.Logger) *TimelineProducer {
	return &TimelineProducer{producer: producer, logger: logger}
}

// ProduceTimelineUpdated publishes a timelines.updated event to Kafka.
func (tp *TimelineProducer) ProduceTimelineUpdated(topic string, event interface{}) error {
	value, err := json.Marshal(event)
	if err != nil {
		tp.logger.Sugar().Errorw("failed to marshal timeline event", "error", err)
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}
	_, _, err = tp.producer.SendMessage(msg)
	if err != nil {
		tp.logger.Sugar().Errorw("failed to produce timeline event", "error", err)
		return err
	}
	return nil
}
