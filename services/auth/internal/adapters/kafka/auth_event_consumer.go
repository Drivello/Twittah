package kafka

import (
	"context"

	"github.com/Drivello/Twittah/services/auth/config"
	"github.com/Drivello/Twittah/services/auth/internal/common"
	"github.com/Drivello/Twittah/services/auth/internal/ports"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// AuthConsumer handles Kafka messages for user-related events.
type AuthConsumer struct {
	AuthProducer             ports.EventProducerPort
	registerUserUseCasesPort ports.RegisterUserUseCasesPort
	Config                   config.ConsumerConfig
	WorkQueue                *common.WorkQueue
	KafkaEventDispatcher     *KafkaEventDispatcher
}

func NewAuthConsumer(producer ports.EventProducerPort, cfg config.ConsumerConfig, wq *common.WorkQueue, authUseCases ports.RegisterUserUseCasesPort) *AuthConsumer {
	consumer := &AuthConsumer{
		AuthProducer:             producer,
		registerUserUseCasesPort: authUseCases,
		Config:                   cfg,
		WorkQueue:                wq,
		KafkaEventDispatcher:     NewKafkaEventDispatcher(authUseCases),
	}

	return consumer
}

func (c *AuthConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (c *AuthConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *AuthConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	topic := claim.Topic()
	common.Logger().Debug("[KafkaConsumer] [UserCreateConsumer] Listening topic", zap.String("topic", topic))

	process := func(item common.WorkItem) {
		req, _ := item.Request.(KafkaEventRequest)
		ctx := item.Ctx
		msg, _ := item.Msg.(sarama.ConsumerMessage)

		err := RetryWithTimeoutAndJitter(ctx, c.Config.RetryConfig, func() error {
			err := c.KafkaEventDispatcher.Dispatch(ctx, msg.Value)
			if err != nil {
				return err.Error
			}
			return nil
		})

		if err != nil {
			common.Logger().Error("[KafkaConsumer] Failed to process event after retries, sending to DLQ",
				zap.String("event_type", req.EventType),
				zap.Error(err))

			dlqErr := SendToDLQ(c.AuthProducer, req.EventType, msg.Value)
			if dlqErr != nil {
				common.Logger().Error("Failed to send message to DLQ", zap.Error(dlqErr))
			}
		}

		sess.MarkMessage(&msg, "")
	}

	c.WorkQueue.Start(process, c.Config.RetryConfig.MaxRetryDuration)
	defer c.WorkQueue.Stop()

	for msg := range claim.Messages() {
		var req KafkaEventRequest

		ctx, cancel := context.WithTimeout(sess.Context(), c.Config.RetryConfig.MaxRetryDuration)

		err := c.WorkQueue.Submit(common.WorkItem{
			Ctx:     ctx,
			Request: req,
			Msg:     *msg,
		})
		if err != nil {
			common.Logger().Warn("[KafkaConsumer] WorkQueue full, sending to DLQ",
				zap.Error(err), zap.Int32("partition", msg.Partition), zap.Int64("offset", msg.Offset))

			dlqErr := SendToDLQ(c.AuthProducer, req.EventType, msg.Value)
			if dlqErr != nil {
				common.Logger().Error("Failed to send message to DLQ (queue full)", zap.Error(dlqErr))
			}
			sess.MarkMessage(msg, "")
		}
		cancel()
	}
	return nil
}
