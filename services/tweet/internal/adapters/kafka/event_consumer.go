package kafka

import (
	"context"

	"github.com/Drivello/Twittah/services/tweet/config"
	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// UserEventConsumer implements sarama.ConsumerGroupHandler for user events (follow/unfollow/etc)
type UserEventConsumer struct {
	Producer        sarama.SyncProducer
	Config          config.ConsumerConfig
	WorkQueue       *common.WorkQueue
	eventDispatcher *KafkaEventDispatcher
}

// NewEventConsumer initializes the consumer and registers handlers per event type
func NewEventConsumer(cfg config.ConsumerConfig, producer sarama.SyncProducer, wq *common.WorkQueue, kafkaEventDispatcher *KafkaEventDispatcher) *UserEventConsumer {
	c := &UserEventConsumer{
		Producer:        producer,
		Config:          cfg,
		WorkQueue:       wq,
		eventDispatcher: kafkaEventDispatcher,
	}

	return c
}

// Setup and Cleanup implement sarama.ConsumerGroupHandler
func (c *UserEventConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (c *UserEventConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim processes messages from the user events topic

func (c *UserEventConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	topic := claim.Topic()
	common.Logger().Info("[UserEventConsumer] ConsumeClaim started for user events topic", zap.String("topic", topic))

	process := func(item common.WorkItem) {
		ctx := item.Ctx
		msg, _ := item.Msg.(sarama.ConsumerMessage)

		kafkaEventError := c.eventDispatcher.Dispatch(ctx, msg.Value)

		if kafkaEventError != nil {
			common.Logger().Error("[KafkaConsumer] Failed to process event, sending to DLQ. Error: ",
				zap.Error(kafkaEventError.Error))

			// dlqErr := SendToDLQ(c.Producer, c.Config.DLQTopic, msg.Value)
			// if dlqErr != nil {
			// 	common.Logger().Error("Failed to send message to DLQ", zap.Error(dlqErr))
			// }
		}

		sess.MarkMessage(&msg, "")
	}

	c.WorkQueue.Start(process, c.Config.RetryConfig.MaxRetryDuration)
	defer c.WorkQueue.Stop()

	for msg := range claim.Messages() {
		var req KafkaEventRequest

		ctx, cancel := context.WithTimeout(sess.Context(), c.Config.RetryConfig.MaxRetryDuration)
		defer cancel()

		err := c.WorkQueue.Submit(common.WorkItem{
			Ctx:     ctx,
			Request: req,
			Msg:     *msg,
		})

		if err != nil {
			common.Logger().Error("[UserEventConsumer] Failed to process user event message",
				zap.Error(err))
			dlqErr := SendToDLQ(c.Producer, c.Config.DLQTopic, msg.Value)
			if dlqErr != nil {
				common.Logger().Error("Failed to send message to DLQ (queue full)", zap.Error(dlqErr))
			}
			sess.MarkMessage(msg, "")
		} else {
			common.Logger().Info("[UserEventConsumer] User event message processed successfully", zap.String("event_type", req.EventType))
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
