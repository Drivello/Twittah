package kafka

import (
	"context"
	"github.com/Drivello/Twittah/services/auth/config"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// StartKafkaConsumers starts both the main user creation consumer and the DLQ worker in separate goroutines.
// This keeps main.go clean and readable.
// StartKafkaConsumers wires and starts both Kafka consumers as background goroutines.
func StartKafkaConsumers(cfg *config.Config, userConsumer *UserCreateConsumer, dlqWorker *DLQWorker) {
	// Start the main user creation consumer goroutine
	go func() {
		group, err := sarama.NewConsumerGroup(cfg.KafkaBrokers, "auth-user-created-group", nil)
		if err != nil {
			zap.L().Fatal("[KafkaSetup] Failed to create Kafka consumer group", zap.Error(err))
		}
		zap.L().Info("[KafkaSetup] Starting user_created consumer group")
		for {
			err := group.Consume(
				context.Background(),
				[]string{"user_created"},
				userConsumer,
			)
			if err != nil {
				zap.L().Error("[KafkaSetup] Kafka user consumer error", zap.Error(err))
			}
		}
	}()

	// Start the DLQ worker goroutine
	go func() {
		group, err := sarama.NewConsumerGroup(cfg.KafkaBrokers, "auth-user-dlq-group", nil)
		if err != nil {
			zap.L().Fatal("[KafkaSetup] Failed to create DLQ Kafka consumer group", zap.Error(err))
		}
		zap.L().Info("[KafkaSetup] Starting DLQ worker consumer group")
		for {
			err := group.Consume(
				context.Background(),
				[]string{cfg.DLQWorker.SourceTopic},
				dlqWorker,
			)
			if err != nil {
				zap.L().Error("[KafkaSetup] DLQ worker consumer error", zap.Error(err))
			}
		}
	}()
}
