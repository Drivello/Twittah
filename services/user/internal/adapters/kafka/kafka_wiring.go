package kafka

import (
	"context"

	"github.com/Drivello/Twittah/services/user/internal/ports"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
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

func StartAllKafkaConsumers(
	repo ports.UserRepository,
	brokers []string,
	logger *zap.Logger,
) (producer sarama.SyncProducer, err error) {
	producer, err = NewSyncProducer(brokers)
	if err != nil {
		return nil, err
	}

	retryCfg := LoadRetryConfig()
	dlqCfg := LoadDLQWorkerConfig()

	userCreatedHandler := &UserCreatedHandler{
		Repo:     repo,
		Producer: producer,
		Config:   retryCfg,
	}
	dlqHandler := &DLQWorker{
		Repo:     repo,
		Producer: producer,
		Config:   dlqCfg,
	}

	go func() {
		err := StartConsumerGroup(
			brokers,
			"user-service-group",
			[]string{"user_created"},
			userCreatedHandler,
		)
		if err != nil {
			logger.Sugar().Fatalw("user_created consumer stopped", "error", err)
		}
	}()

	go func() {
		err := StartConsumerGroup(
			brokers,
			"user-dlq-worker-group",
			[]string{dlqCfg.SourceTopic},
			dlqHandler,
		)
		if err != nil {
			logger.Sugar().Fatalw("DLQ worker consumer stopped", "error", err)
		}
	}()

	return producer, nil
}
