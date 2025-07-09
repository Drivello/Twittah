// Package kafka provides Kafka consumers for TweetService.
package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/Drivello/Twittah/services/tweet/config"
	"github.com/Drivello/Twittah/services/tweet/internal/domain"
	"github.com/Drivello/Twittah/services/tweet/internal/ports"
	"github.com/Drivello/Twittah/services/tweet/internal/usecase"
)

// parseTimeOrNow parses an RFC3339 string or returns time.Now() if invalid
func parseTimeOrNow(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Now()
	}
	return t
}

// ConsumerGroupHandler handles Kafka events for TweetService.
// ConsumerGroupHandler handles Kafka events for TweetService, including retries and DLQ.
type ConsumerGroupHandler struct {
	publishUC   *usecase.PublishTweet
	timelineUC  *usecase.GetTimeline
	cache       ports.Cache
	logger      *zap.Logger
	dlqProducer sarama.SyncProducer // Producer for sending failed messages to DLQ
}

// NewConsumerGroupHandler creates a new ConsumerGroupHandler.
// NewConsumerGroupHandler creates a new ConsumerGroupHandler.
func NewConsumerGroupHandler(publishUC *usecase.PublishTweet, timelineUC *usecase.GetTimeline, cache ports.Cache, logger *zap.Logger, dlqProducer sarama.SyncProducer) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{
		publishUC:   publishUC,
		timelineUC:  timelineUC,
		cache:       cache,
		logger:      logger,
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
	var dto TweetPublishedEventDTO
	if err := json.Unmarshal(msg.Value, &dto); err != nil {
		h.logger.Sugar().Errorw("invalid tweet published event", "error", err)
		return err
	}
	domainTweet := &domain.Tweet{
		ID:        dto.ID,
		AuthorID:  dto.AuthorID,
		Content:   dto.Content,
		CreatedAt: parseTimeOrNow(dto.CreatedAt),
	}
	return h.publishUC.Execute(context.Background(), domainTweet, nil)
}

// handleFollowCreated processes a follow created event from Kafka.
// Returns an error if processing fails.
func (h *ConsumerGroupHandler) handleFollowCreated(msg *sarama.ConsumerMessage) error {
	var dto FollowCreatedEventDTO
	if err := json.Unmarshal(msg.Value, &dto); err != nil {
		h.logger.Sugar().Errorw("invalid follow created event", "error", err)
		return err
	}
	// TODO: Implementar lógica de actualización de timeline para follow creado
	// Por ahora, no se realiza ninguna acción
	return nil
}

// handleFollowDeleted processes a follow deleted event from Kafka.
// Returns an error if processing fails.
func (h *ConsumerGroupHandler) handleFollowDeleted(msg *sarama.ConsumerMessage) error {
	var dto FollowDeletedEventDTO
	if err := json.Unmarshal(msg.Value, &dto); err != nil {
		h.logger.Sugar().Errorw("invalid follow deleted event", "error", err)
		return err
	}
	// TODO: Implementar lógica de actualización de timeline para follow eliminado
	// Por ahora, no se realiza ninguna acción
	return nil
}
