package kafka

import (
	"context"

	"github.com/Drivello/Twittah/services/tweet/config"
	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"github.com/IBM/sarama"
)

// NewSyncProducer creates a SyncProducer for Kafka
func NewSyncProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	return sarama.NewSyncProducer(brokers, config)
}

func StartConsumerGroup(brokers []string, group string, topics []string, handler sarama.ConsumerGroupHandler) error {
	config := sarama.NewConfig()
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
	tweetConsumerGroup sarama.ConsumerGroup
	userConsumerGroup  sarama.ConsumerGroup
)

func StartKafkaConsumers(ctx context.Context, cfg *config.Config, kafkaEventDispatcher *KafkaEventDispatcher, producer sarama.SyncProducer, workQueue *common.WorkQueue) {
	common.Logger().Debug("[KafkaSetup] Entrando a StartKafkaConsumers")

	tweetConsumer := NewEventConsumer(cfg.KafkaTweetConsumerConfig, producer, workQueue, kafkaEventDispatcher)
	userConsumer := NewEventConsumer(cfg.KafkaUserConsumerConfig, producer, workQueue, kafkaEventDispatcher)

	// Tweet consumer goroutine (factory)
	StartKafkaConsumerGoroutine(
		ctx,
		cfg.KafkaBrokers,
		cfg.KafkaTweetConsumerConfig.Group,
		[]string{cfg.KafkaTweetConsumerConfig.Topic, cfg.KafkaTweetConsumerConfig.DLQTopic},
		tweetConsumer,
		"[TweetEventConsumer]",
		&tweetConsumerGroup,
	)

	// User consumer goroutine (factory)
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
