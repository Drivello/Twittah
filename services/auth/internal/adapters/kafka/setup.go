package kafka

import (
	"context"

	"github.com/Drivello/Twittah/services/auth/config"
	"github.com/Drivello/Twittah/services/auth/internal/common"
	"github.com/Drivello/Twittah/services/auth/internal/ports"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// StartKafkaConsumers starts both the main user creation consumer and the DLQ worker in separate goroutines.
// This keeps main.go clean and readable.
// StartKafkaConsumers wires and starts both Kafka consumers as background goroutines.
var (
	userConsumerGroup sarama.ConsumerGroup
	dlqConsumerGroup  sarama.ConsumerGroup
)

func StartKafkaConsumers(cfg *config.Config, repo ports.UserRepository, userUC ports.UserUseCasePort, producer sarama.SyncProducer) {
	common.Logger().Debug("[KafkaSetup] Entrando a StartKafkaConsumers")
	userConsumer := &UserCreateConsumer{
		UserUC:    userUC,
		Producer:  producer,
		Config:    cfg.Retry,
		EventType: cfg.UserCreatedEventType,
	}
	dlqWorker := &DLQWorker{
		Repo:     repo,
		Producer: producer,
		Config:   cfg.DLQWorker,
	}

	// Start the main user creation consumer goroutine
	common.Logger().Debug("[KafkaSetup] Antes de lanzar goroutine de user consumer", zap.String("group", cfg.KafkaAuthConsumerGroupName), zap.String("topic", cfg.KafkaUserEventsTopic))
	go func() {
		common.Logger().Debug("[KafkaSetup] Dentro de goroutine de user consumer")
		group, err := sarama.NewConsumerGroup(cfg.KafkaBrokers, cfg.KafkaAuthConsumerGroupName, nil)
		if err != nil {
			common.Logger().Fatal("[KafkaSetup] Failed to create Kafka consumer group", zap.Error(err))
		}
		userConsumerGroup = group
		common.Logger().Info("[KafkaSetup] Starting" + cfg.KafkaUserEventsTopic + " consumer group")
		for {
			common.Logger().Debug("[KafkaSetup] Antes de llamar a group.Consume para user consumer")
			err := group.Consume(
				context.Background(),
				[]string{cfg.KafkaUserEventsTopic},
				userConsumer,
			)
			if err != nil {
				common.Logger().Error("[KafkaSetup] Kafka user consumer error", zap.Error(err))
			}
			common.Logger().Debug("[KafkaSetup] Terminó una iteración de group.Consume para user consumer")
		}
	}()

	// Start the DLQ worker goroutine
	common.Logger().Debug("[KafkaSetup] Antes de lanzar goroutine de DLQ worker", zap.String("group", cfg.KafkaAuthDLQConsumerGroupName), zap.String("topic", cfg.DLQWorker.SourceTopic))
	go func() {
		common.Logger().Debug("[KafkaSetup] Dentro de goroutine de DLQ worker")
		group, err := sarama.NewConsumerGroup(cfg.KafkaBrokers, cfg.KafkaAuthDLQConsumerGroupName, nil)
		if err != nil {
			common.Logger().Fatal("[KafkaSetup] Failed to create DLQ Kafka consumer group", zap.Error(err))
		}
		dlqConsumerGroup = group
		common.Logger().Info("[KafkaSetup] Starting DLQ worker consumer group")
		for {
			common.Logger().Debug("[KafkaSetup] Antes de llamar a group.Consume para DLQ worker")
			err := group.Consume(
				context.Background(),
				[]string{cfg.DLQWorker.SourceTopic},
				dlqWorker,
			)
			if err != nil {
				common.Logger().Error("[KafkaSetup] DLQ worker consumer error", zap.Error(err))
			}
			common.Logger().Debug("[KafkaSetup] Terminó una iteración de group.Consume para DLQ worker")
		}
	}()
}

// CloseKafka shuts down the Kafka consumer groups and producer gracefully.
func CloseKafka(producer sarama.SyncProducer) {
	if userConsumerGroup != nil {
		if err := userConsumerGroup.Close(); err != nil {
			common.Logger().Error("Failed to close user consumer group", zap.Error(err))
		}
	}
	if dlqConsumerGroup != nil {
		if err := dlqConsumerGroup.Close(); err != nil {
			common.Logger().Error("Failed to close DLQ consumer group", zap.Error(err))
		}
	}
	if producer != nil {
		if err := producer.Close(); err != nil {
			common.Logger().Error("Failed to close Kafka producer", zap.Error(err))
		}
	}
}
