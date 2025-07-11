package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Drivello/Twittah/services/auth/config"
	"github.com/Drivello/Twittah/services/auth/internal/common"
	"github.com/Drivello/Twittah/services/auth/internal/ports"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// DLQWorker handles messages from the Dead Letter Queue (DLQ) and retries processing.
type DLQWorker struct {
	AuthUC   ports.AuthUseCasesPort
	Producer sarama.SyncProducer
	Config   config.DLQWorkerConfig

	handlers map[string]func(context.Context, map[string]interface{}) error
}

// Setup and Cleanup implement sarama.ConsumerGroupHandler for DLQWorker.
func (w *DLQWorker) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (w *DLQWorker) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// NewDLQWorker creates a DLQWorker and registers all handlers.
func NewDLQWorker(authUC ports.AuthUseCasesPort, producer sarama.SyncProducer, cfg config.DLQWorkerConfig) *DLQWorker {
	worker := &DLQWorker{
		AuthUC:   authUC,
		Producer: producer,
		Config:   cfg,
		handlers: make(map[string]func(context.Context, map[string]interface{}) error),
	}

	// Register handlers for each event type
	worker.handlers["users.create"] = func(ctx context.Context, payload map[string]interface{}) error {
		return HandleUserCreate(ctx, worker.AuthUC, payload)
	}

	return worker
}

// ConsumeClaim processes messages from the DLQ topic and retries them.
func (w *DLQWorker) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	common.Logger().Info("[DLQWorker] ConsumeClaim started for DLQ topic")
	for msg := range claim.Messages() {
		// Deserialize the DLQ message
		var req UserEventRequest
		if err := json.Unmarshal(msg.Value, &req); err != nil {
			common.Logger().Error("[DLQWorker] Invalid DLQ message, skipping", zap.Error(err))
			sess.MarkMessage(msg, "")
			continue
		}

		if req.EventType == "" {
			common.Logger().Error("[DLQWorker] Missing event_type in message, skipping")
			sess.MarkMessage(msg, "")
			continue
		}

		// TTL check: skip old messages
		if w.Config.DLQMessageTTL > 0 {
			msgAge := time.Since(msg.Timestamp)
			if msgAge > w.Config.DLQMessageTTL {
				common.Logger().Warn("[DLQWorker] Discarding DLQ message by TTL", zap.Duration("age", msgAge), zap.Duration("ttl", w.Config.DLQMessageTTL))
				sess.MarkMessage(msg, "")
				continue
			}
		}

		// Find handler
		handler, exists := w.handlers[req.EventType]
		if !exists {
			common.Logger().Error("[DLQWorker] No handler registered for event_type", zap.String("event_type", req.EventType))
			sess.MarkMessage(msg, "")
			continue
		}

		// Execute handler
		ctx, cancel := context.WithTimeout(sess.Context(), w.Config.MaxRetryDuration)
		defer cancel()

		err := handler(ctx, req.Payload)
		if err != nil {
			common.Logger().Error("[DLQWorker] Failed to process DLQ message, discarding",
				zap.String("event_type", req.EventType),
				zap.Error(err))
		} else {
			common.Logger().Info("[DLQWorker] DLQ message processed successfully", zap.String("event_type", req.EventType))
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
