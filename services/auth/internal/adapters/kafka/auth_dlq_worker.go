package kafka

import (
	"context"
	"time"

	"github.com/Drivello/Twittah/services/auth/config"
	"github.com/Drivello/Twittah/services/auth/internal/common"
	"github.com/Drivello/Twittah/services/auth/internal/ports"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// DLQWorker handles messages from the Dead Letter Queue (DLQ) and retries processing.
type DLQWorker struct {
	Producer             ports.EventProducerPort
	AuthUC               ports.RegisterUserUseCasesPort
	Config               config.ConsumerConfig
	kafkaEventDispatcher *KafkaEventDispatcher
}

// Setup and Cleanup implement sarama.ConsumerGroupHandler for DLQWorker.
func (w *DLQWorker) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (w *DLQWorker) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// NewDLQWorker creates a DLQWorker and registers all handlers.
func NewDLQWorker(producer ports.EventProducerPort, authUC ports.RegisterUserUseCasesPort, cfg config.ConsumerConfig) *DLQWorker {
	worker := &DLQWorker{
		Producer:             producer,
		AuthUC:               authUC,
		Config:               cfg,
		kafkaEventDispatcher: NewKafkaEventDispatcher(authUC),
	}

	return worker
}

// ConsumeClaim processes messages from the DLQ topic and retries them.
func (w *DLQWorker) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	common.Logger().Info("[DLQWorker] ConsumeClaim started for DLQ topic")
	for msg := range claim.Messages() {

		// TTL check: skip old messages
		if w.Config.DLQConfig.TTL > 0 {
			msgAge := time.Since(msg.Timestamp)
			if msgAge > w.Config.TTL {
				common.Logger().Warn("[DLQWorker] Discarding DLQ message by TTL", zap.Duration("age", msgAge), zap.Duration("ttl", w.Config.TTL))
				sess.MarkMessage(msg, "")
				continue
			}
		}

		// Execute handler
		ctx, cancel := context.WithTimeout(sess.Context(), w.Config.RetryConfig.MaxRetryDuration)
		defer cancel()

		err := w.kafkaEventDispatcher.Dispatch(ctx, msg.Value)
		if err != nil {
			common.Logger().Error("[DLQWorker] Failed to process DLQ message, discarding",
				zap.Error(err.Error))
		} else {
			common.Logger().Info("[DLQWorker] DLQ message processed successfully")
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
