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
	authConsumerGroup sarama.ConsumerGroup
	dlqConsumerGroup  sarama.ConsumerGroup
)

func StartKafkaConsumers(ctx context.Context, cfg *config.Config, repo ports.AuthRepositoryPort, registerUserUC ports.RegisterUserUseCasesPort, producer ports.EventProducerPort, workQueue *common.WorkQueue) {
	common.Logger().Debug("[KafkaSetup] Entrando a StartKafkaConsumers")

	authConsumer := NewAuthConsumer(producer, cfg.KafkaUserConsumerConfig, workQueue, registerUserUC)
	authDLQWorker := NewDLQWorker(producer, registerUserUC, cfg.KafkaUserConsumerConfig)

	// User consumer goroutine
	go func() {
		group, err := sarama.NewConsumerGroup(cfg.KafkaBrokers, cfg.KafkaUserConsumerConfig.Group, nil)
		if err != nil {
			common.Logger().Fatal("[KafkaSetup] Failed to create Kafka consumer group", zap.Error(err))
		}
		authConsumerGroup = group
		common.Logger().Info("[KafkaSetup] Starting " + cfg.KafkaUserConsumerConfig.Topic + " consumer group")

		for {
			select {
			case <-ctx.Done():
				common.Logger().Info("[KafkaSetup] Context canceled, shutting down user consumer")
				return
			default:
				err := group.Consume(
					ctx,
					[]string{cfg.KafkaUserConsumerConfig.Topic},
					authConsumer,
				)
				if err != nil {
					if err == context.Canceled {
						common.Logger().Info("[KafkaSetup] User consumer context canceled, exiting")
						return
					}
					common.Logger().Error("[KafkaSetup] Kafka user consumer error", zap.Error(err))
				}
			}
		}
	}()

	// DLQ worker goroutine
	go func() {
		group, err := sarama.NewConsumerGroup(cfg.KafkaBrokers, cfg.KafkaUserConsumerConfig.Group+"-dlq", nil)
		if err != nil {
			common.Logger().Fatal("[KafkaSetup] Failed to create DLQ Kafka consumer group", zap.Error(err))
		}
		dlqConsumerGroup = group
		common.Logger().Info("[KafkaSetup] Starting DLQ worker " + cfg.KafkaUserConsumerConfig.DLQTopic + " consumer group")

		for {
			select {
			case <-ctx.Done():
				common.Logger().Info("[KafkaSetup] Context canceled, shutting down DLQ worker")
				return
			default:
				err := group.Consume(
					ctx,
					[]string{cfg.KafkaUserConsumerConfig.DLQTopic},
					authDLQWorker,
				)
				if err != nil {
					if err == context.Canceled {
						common.Logger().Info("[KafkaSetup] DLQ worker context canceled, exiting")
						return
					}
					common.Logger().Error("[KafkaSetup] DLQ worker consumer error", zap.Error(err))
				}
			}
		}
	}()
}

// CloseKafka closes Kafka producer and consumers gracefully.
func CloseKafka(producer sarama.SyncProducer) {
	logger := common.Logger()

	// Cerrar producer
	if producer != nil {
		if err := producer.Close(); err != nil {
			logger.Error("[Kafka] Failed to close producer", zap.Error(err))
		} else {
			logger.Info("[Kafka] Producer closed successfully")
		}
	}

	// Cerrar authConsumerGroup (si existe)
	if authConsumerGroup != nil {
		if err := authConsumerGroup.Close(); err != nil {
			logger.Error("[Kafka] Failed to close auth consumer group", zap.Error(err))
		} else {
			logger.Info("[Kafka] Auth consumer group closed successfully")
		}
	}

	// Cerrar dlqConsumerGroup (si existe)
	if dlqConsumerGroup != nil {
		if err := dlqConsumerGroup.Close(); err != nil {
			logger.Error("[Kafka] Failed to close DLQ consumer group", zap.Error(err))
		} else {
			logger.Info("[Kafka] DLQ consumer group closed successfully")
		}
	}
}
