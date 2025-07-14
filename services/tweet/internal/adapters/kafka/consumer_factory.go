package kafka

import (
	"context"

	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"github.com/IBM/sarama"
)

// StartKafkaConsumerGoroutine centraliza la lógica de arranque y ciclo de vida de un consumidor Kafka en una goroutine.
func StartKafkaConsumerGoroutine(
	ctx context.Context,
	brokers []string,
	groupName string,
	topics []string,
	consumer sarama.ConsumerGroupHandler,
	logPrefix string,
	groupVar *sarama.ConsumerGroup,
) {
	go func() {
		group, err := sarama.NewConsumerGroup(brokers, groupName, nil)
		if err != nil {
			common.Logger().Fatalf("[%s] Failed to create Kafka consumer group: %v", logPrefix, err)
		}
		if groupVar != nil {
			*groupVar = group
		}
		common.Logger().Infof("[%s] Starting consumer group: %s", logPrefix, groupName)
		for {
			select {
			case <-ctx.Done():
				common.Logger().Infof("[%s] Context canceled, shutting down consumer", logPrefix)
				return
			default:
				common.Logger().Debug("[%s] Consuming topics", logPrefix, "topics", topics)
				err := group.Consume(ctx, topics, consumer)
				if err != nil {
					if err == context.Canceled {
						common.Logger().Infof("[%s] Consumer context canceled, exiting", logPrefix)
						return
					}
					common.Logger().Errorf("[%s] Consumer error: %v", logPrefix, err)
				}
			}
		}
	}()
}
