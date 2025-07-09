// Package kafka provides Kafka consumers for TweetService.
package kafka

import (
	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/Drivello/Twittah/services/tweet/internal/ports"
	"github.com/Drivello/Twittah/services/tweet/internal/usecase"
)

// ConsumerGroupHandler handles Kafka events for TweetService.
type ConsumerGroupHandler struct {
	publishUC  *usecase.PublishTweet
	timelineUC *usecase.GetTimeline
	cache      ports.Cache
	logger     *zap.Logger
}

// NewConsumerGroupHandler creates a new ConsumerGroupHandler.
func NewConsumerGroupHandler(publishUC *usecase.PublishTweet, timelineUC *usecase.GetTimeline, cache ports.Cache, logger *zap.Logger) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{
		publishUC:  publishUC,
		timelineUC: timelineUC,
		cache:      cache,
		logger:     logger,
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim processes messages from Kafka.
func (h *ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.logger.Sugar().Infow("kafka event received", "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset)
		switch msg.Topic {
		case "tweets.published":
			h.handleTweetPublished(msg)
		case "follows.created":
			h.handleFollowCreated(msg)
		case "follows.deleted":
			h.handleFollowDeleted(msg)
		default:
			h.logger.Sugar().Warnw("unknown kafka topic", "topic", msg.Topic)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (h *ConsumerGroupHandler) handleTweetPublished(msg *sarama.ConsumerMessage) {
	// TODO: Unmarshal event and call publishUC if needed
	h.logger.Sugar().Infow("handleTweetPublished stub", "value", string(msg.Value))
}

func (h *ConsumerGroupHandler) handleFollowCreated(msg *sarama.ConsumerMessage) {
	// TODO: Unmarshal event and update timeline cache
	h.logger.Sugar().Infow("handleFollowCreated stub", "value", string(msg.Value))
}

func (h *ConsumerGroupHandler) handleFollowDeleted(msg *sarama.ConsumerMessage) {
	// TODO: Unmarshal event and update timeline cache
	h.logger.Sugar().Infow("handleFollowDeleted stub", "value", string(msg.Value))
}
