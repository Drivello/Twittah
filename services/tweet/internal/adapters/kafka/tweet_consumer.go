// Package kafka provides Kafka consumers for TweetService.
package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/Drivello/Twittah/services/tweet/config"
	"github.com/Drivello/Twittah/services/tweet/internal/ports"
	"github.com/Drivello/Twittah/services/tweet/internal/usecase"
)

// ConsumerGroupHandler handles Kafka events for TweetService.
// ConsumerGroupHandler handles Kafka events for TweetService, including retries and DLQ.
type ConsumerGroupHandler struct {
	publishUC    *usecase.PublishTweet
	timelineUC   *usecase.GetTimeline
	cache        ports.Cache
	logger       *zap.Logger
	dlqProducer  sarama.SyncProducer // Producer for sending failed messages to DLQ
}

// NewConsumerGroupHandler creates a new ConsumerGroupHandler.
// NewConsumerGroupHandler creates a new ConsumerGroupHandler.
func NewConsumerGroupHandler(publishUC *usecase.PublishTweet, timelineUC *usecase.GetTimeline, cache ports.Cache, logger *zap.Logger, dlqProducer sarama.SyncProducer) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{
		publishUC:  publishUC,
		timelineUC: timelineUC,
		cache:      cache,
		logger:     logger,
		dlqProducer: dlqProducer,
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim processes messages from Kafka.
// ConsumeClaim processes messages from Kafka, handling retries and sending failed messages to a DLQ.
func (h *ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	retryCfg := config.LoadRetryConfig()
	for msg := range claim.Messages() {
		h.logger.Sugar().Infow("kafka event received", "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset)

		ctx, cancel := context.WithTimeout(sess.Context(), retryCfg.MaxRetryDuration)
		defer cancel()

		err := retryWithTimeoutAndJitter(ctx, retryCfg, func() error {
			switch msg.Topic {
			case "tweets.published":
				return h.handleTweetPublished(msg)
			case "follows.created":
				return h.handleFollowCreated(msg)
			case "follows.deleted":
				return h.handleFollowDeleted(msg)
			default:
				h.logger.Sugar().Warnw("unknown kafka topic", "topic", msg.Topic)
				return nil
			}
		}, h.logger)

		if err != nil {
			h.logger.Sugar().Errorw("Failed to process message after retries, sending to DLQ", "topic", msg.Topic, "error", err, "dlq_topic", retryCfg.DLQTopic)
			if h.dlqProducer != nil {
				dlqMsg := &sarama.ProducerMessage{
					Topic: retryCfg.DLQTopic,
					Value: sarama.ByteEncoder(msg.Value),
				}
				_, _, prodErr := h.dlqProducer.SendMessage(dlqMsg)
				if prodErr != nil {
					h.logger.Sugar().Errorw("failed to send message to DLQ", "error", prodErr)
				}
			}
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

// handleTweetPublished processes a tweet published event from Kafka.
// Returns an error if processing fails.
func (h *ConsumerGroupHandler) handleTweetPublished(msg *sarama.ConsumerMessage) error {
	// TODO: Unmarshal event and call publishUC if needed
	h.logger.Sugar().Infow("handleTweetPublished stub", "value", string(msg.Value))
	return nil // Cambia esto por error real si implementas lógica
}

// handleFollowCreated processes a follow created event from Kafka.
// Returns an error if processing fails.
func (h *ConsumerGroupHandler) handleFollowCreated(msg *sarama.ConsumerMessage) error {
	// TODO: Unmarshal event and update timeline cache
	h.logger.Sugar().Infow("handleFollowCreated stub", "value", string(msg.Value))
	return nil // Cambia esto por error real si implementas lógica
}

// handleFollowDeleted processes a follow deleted event from Kafka.
// Returns an error if processing fails.
func (h *ConsumerGroupHandler) handleFollowDeleted(msg *sarama.ConsumerMessage) error {
	// TODO: Unmarshal event and update timeline cache
	h.logger.Sugar().Infow("handleFollowDeleted stub", "value", string(msg.Value))
	return nil // Cambia esto por error real si implementas lógica
}
