package kafka

import (
	"github.com/Drivello/Twittah/services/user/internal/common"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// DefaultKafkaWorkProcess is a generic process function for WorkQueue to handle Kafka WorkItems
func DefaultKafkaWorkProcess(eventDispatcher *KafkaEventDispatcher, producer interface{}, dlqTopic string) func(common.WorkItem) {
	return func(item common.WorkItem) {
		ctx := item.Ctx
		msg, ok := item.Msg.(sarama.ConsumerMessage)
		if !ok {
			common.Logger().Error("[WorkerPool] WorkItem.Msg is not a sarama.ConsumerMessage")
			return
		}

		kafkaEventError := eventDispatcher.Dispatch(ctx, msg.Value)
		if kafkaEventError != nil {
			common.Logger().Error("[WorkerPool] Failed to process event, sending to DLQ.", zap.Error(kafkaEventError.Error))
			if prod, ok := producer.(sarama.SyncProducer); ok {
				dlqErr := SendToDLQ(prod, dlqTopic, msg.Value)
				if dlqErr != nil {
					common.Logger().Error("[WorkerPool] Failed to send message to DLQ", zap.Error(dlqErr))
				}
			}
		}
	}
}
