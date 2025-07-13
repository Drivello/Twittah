package kafka

import (
	"context"

	"github.com/Drivello/Twittah/services/user/config"
	"github.com/Drivello/Twittah/services/user/internal/common"
	"github.com/Drivello/Twittah/services/user/internal/ports"
	"github.com/IBM/sarama"
)

// NewSyncProducer creates a SyncProducer for Kafka
func NewSyncProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Producer.Return.Successes = true
	return sarama.NewSyncProducer(brokers, config)
}

func StartConsumerGroup(brokers []string, group string, topics []string, handler sarama.ConsumerGroupHandler) error {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Version = sarama.V2_1_0_0
	consumer, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		return err
	}
	ctx := context.Background()
	for {
		if err := consumer.Consume(ctx, topics, handler); err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

var (
	userConsumerGroup sarama.ConsumerGroup
)

func StartKafkaConsumers(ctx context.Context, cfg *config.Config, repo ports.UserRepository, kafkaEventDispatcher *KafkaEventDispatcher, producer sarama.SyncProducer, workQueue *common.WorkQueue) {
	common.Logger().Debug("[KafkaSetup] Entrando a StartKafkaConsumers")

	userConsumer := NewEventConsumer(cfg.KafkaUserConsumerConfig, producer, workQueue, kafkaEventDispatcher)

	StartKafkaConsumerGoroutine(
		ctx,
		cfg.KafkaBrokers,
		cfg.KafkaUserConsumerConfig.Group,
		[]string{cfg.KafkaUserConsumerConfig.Topic, cfg.KafkaUserConsumerConfig.DLQTopic},
		userConsumer,
		"[UserEventConsumer]",
		&userConsumerGroup,
	)

}
