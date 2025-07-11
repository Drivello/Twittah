package kafka

import (
	"context"
	"encoding/json"

	"github.com/Drivello/Twittah/services/auth/config"
	"github.com/Drivello/Twittah/services/auth/internal/common"
	"github.com/Drivello/Twittah/services/auth/internal/ports"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// AuthConsumer handles Kafka messages for user-related events.
type AuthConsumer struct {
	AuthProducer sarama.SyncProducer
	AuthUseCases ports.AuthUseCasesPort
	Config       config.RetryConfig
	WorkQueue    *common.WorkQueue
	handlers     map[string]func(context.Context, map[string]interface{}) error
}

func NewAuthConsumer(producer sarama.SyncProducer, cfg config.RetryConfig, wq *common.WorkQueue, authUseCases ports.AuthUseCasesPort) *AuthConsumer {
	consumer := &AuthConsumer{
		AuthProducer: producer,
		AuthUseCases: authUseCases,
		Config:       cfg,
		WorkQueue:    wq,
		handlers:     make(map[string]func(context.Context, map[string]interface{}) error),
	}

	// Register handlers for each event type
	consumer.handlers["users.create"] = func(ctx context.Context, payload map[string]interface{}) error {
		return HandleUserCreate(ctx, consumer.AuthUseCases, payload)
	}

	return consumer
}

func (c *AuthConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (c *AuthConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *AuthConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	topic := claim.Topic()
	common.Logger().Debug("[KafkaConsumer] [UserCreateConsumer] Listening topic", zap.String("topic", topic))

	process := func(item common.WorkItem) {
		req, _ := item.Request.(UserEventRequest)
		ctx := item.Ctx
		msg, _ := item.Msg.(sarama.ConsumerMessage)

		common.Logger().Debug("Raw Kafka message", zap.ByteString("payload", msg.Value), zap.String("event_type", req.EventType))

		handler, ok := c.handlers[req.EventType]
		if !ok {
			common.Logger().Error("Unknown event type. Skipping.",
				zap.String("event_type", req.EventType))
			sess.MarkMessage(&msg, "")
			return
		}

		err := RetryWithTimeoutAndJitter(ctx, c.Config, func() error {
			return handler(ctx, req.Payload)
		})

		if err != nil {
			common.Logger().Error("[KafkaConsumer] Failed to process event after retries, sending to DLQ",
				zap.String("event_type", req.EventType),
				zap.Error(err))

			dlqErr := SendToDLQ(c.AuthProducer, c.Config.DLQTopic, req)
			if dlqErr != nil {
				common.Logger().Error("Failed to send message to DLQ", zap.Error(dlqErr))
			}
		}

		sess.MarkMessage(&msg, "")
	}

	c.WorkQueue.Start(process, c.Config.MaxRetryDuration)
	defer c.WorkQueue.Stop()

	for msg := range claim.Messages() {
		var req UserEventRequest

		for k := range c.handlers {
			common.Logger().Debug("Registered handler", zap.String("key", k))
		}

		common.Logger().Debug("Parsed UserEventRequest",
			zap.String("event_type", req.EventType),
			zap.Any("payload", req.Payload))

		if err := json.Unmarshal(msg.Value, &req); err != nil {
			common.Logger().Error("[KafkaConsumer] Failed to unmarshal user event request",
				zap.Error(err), zap.Int32("partition", msg.Partition), zap.Int64("offset", msg.Offset))
			sess.MarkMessage(msg, "")
			continue
		}

		if req.EventType == "" {
			common.Logger().Error("Missing event_type in message. Skipping.",
				zap.Int32("partition", msg.Partition), zap.Int64("offset", msg.Offset))
			sess.MarkMessage(msg, "")
			continue
		}

		err := c.WorkQueue.Submit(common.WorkItem{
			Ctx:     context.Background(),
			Request: req,
			Msg:     *msg,
		})
		if err != nil {
			common.Logger().Warn("[KafkaConsumer] WorkQueue full, sending to DLQ",
				zap.Error(err), zap.Int32("partition", msg.Partition), zap.Int64("offset", msg.Offset))

			dlqErr := SendToDLQ(c.AuthProducer, c.Config.DLQTopic, req)
			if dlqErr != nil {
				common.Logger().Error("Failed to send message to DLQ (queue full)", zap.Error(dlqErr))
			}
			sess.MarkMessage(msg, "")
		}
	}
	return nil
}
