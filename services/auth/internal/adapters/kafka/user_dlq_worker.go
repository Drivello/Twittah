package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/Drivello/Twittah/services/auth/config"
	"github.com/Drivello/Twittah/services/auth/internal/domain"
	"github.com/Drivello/Twittah/services/auth/internal/ports"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// DLQWorker handles messages from the Dead Letter Queue (DLQ) and retries user creation.

type DLQWorker struct {
	Repo     ports.UserRepository
	Producer sarama.SyncProducer
	Config   config.DLQWorkerConfig
}

// Setup and Cleanup implement sarama.ConsumerGroupHandler for DLQWorker.
func (w *DLQWorker) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (w *DLQWorker) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim processes messages from the DLQ topic and retries user creation.
// If retries fail, the message is sent to the final DLQ.
func (w *DLQWorker) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	zap.L().Info("[DLQWorker] ConsumeClaim started for DLQ topic")
	for msg := range claim.Messages() {
		// Process each message from the DLQ topic
		var event struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			zap.L().Error("[DLQWorker] Invalid DLQ message, skipping", zap.Error(err))
			sess.MarkMessage(msg, "")
			continue // Skip invalid DLQ message
		}

		ctx, cancel := context.WithTimeout(sess.Context(), w.Config.MaxRetryDuration)
		defer cancel()

		// Attempt user creation with retry and backoff
		err := retryWithTimeoutAndJitterDQL(ctx, w.Config, func() error {
			user := &domain.User{
				Username: event.Username,
				Email:    event.Email,
				Password: event.Password,
			}
			_, err := w.Repo.CreateUser(ctx, user)
			return err
		})

		if err != nil {
			// If all retries fail, send to final DLQ
			zap.L().Error("[DLQWorker] Failed to process DLQ message after retries, sending to final DLQ", zap.Error(err), zap.String("final_dlq", w.Config.FinalDLQTopic))
			dlqErr := requeueToDLQ(w.Producer, w.Config.FinalDLQTopic, event)
			if dlqErr != nil {
				zap.L().Error("[DLQWorker] Failed to send message to final DLQ", zap.Error(dlqErr))
			}
		}

		// Successfully processed or handled message
		sess.MarkMessage(msg, "")
	}
	return nil
}

// retryWithTimeoutAndJitterDQL retries a function using exponential backoff and jitter until timeout for DLQ worker.
func retryWithTimeoutAndJitterDQL(ctx context.Context, cfg config.DLQWorkerConfig, fn func() error) error {
	backoff := cfg.InitialBackoff

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry timed out after %v: %w", cfg.MaxRetryDuration, ctx.Err())
		default:
			err := fn()
			if err == nil {
				return nil
			}

			zap.L().Warn("DLQ worker retryable error, will retry", zap.Error(err), zap.Duration("next_backoff", backoff))

			// Sleep con jitter (aleatoriza entre 50% y 150% del backoff)
			jitter := time.Duration(rand.Int63n(int64(backoff))) / 2
			sleepDuration := backoff + jitter
			time.Sleep(sleepDuration)

			// Exponential backoff hasta el máximo
			if backoff < cfg.MaxBackoff {
				backoff *= 2
				if backoff > cfg.MaxBackoff {
					backoff = cfg.MaxBackoff
				}
			}
		}
	}
}

// requeueToDLQ serializes the event and sends it to the DLQ topic using the provided Kafka producer.
func requeueToDLQ(producer sarama.SyncProducer, topic string, event interface{}) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal DLQ payload: %w", err)
	}

	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(payload),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to DLQ: %w", err)
	}

	zap.L().Info("Message sent to DLQ", zap.String("topic", topic))
	return nil
}
