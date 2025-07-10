package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/Drivello/Twittah/services/user/internal/ports"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// -----------------------------
// DLQ Worker Configuration
// -----------------------------
type DLQWorkerConfig struct {
	MaxRetryDuration time.Duration
	MaxBackoff       time.Duration
	InitialBackoff   time.Duration
	SourceTopic      string
}

func LoadDLQWorkerConfig() DLQWorkerConfig {
	return DLQWorkerConfig{
		MaxRetryDuration: getEnvAsDuration("DLQ_WORKER_MAX_DURATION", 1*time.Hour),
		MaxBackoff:       getEnvAsDuration("DLQ_WORKER_MAX_BACKOFF", 5*time.Minute),
		InitialBackoff:   getEnvAsDuration("DLQ_WORKER_INITIAL_BACKOFF", 5*time.Second),
		SourceTopic:      getEnv("DLQ_WORKER_SOURCE_TOPIC", "user.created.dlq"),
	}
}

// -----------------------------
// Worker Struct
// -----------------------------
type DLQWorker struct {
	Repo     ports.UserRepository
	Producer sarama.SyncProducer
	Config   DLQWorkerConfig
}

// -----------------------------
// ConsumerGroupHandler for DLQ
// -----------------------------
func (w *DLQWorker) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (w *DLQWorker) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (w *DLQWorker) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var event UserCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			zap.S().Errorw("Invalid DLQ message, skipping", "error", err)
			sess.MarkMessage(msg, "")
			continue
		}

		ctx, cancel := context.WithTimeout(sess.Context(), w.Config.MaxRetryDuration)
		defer cancel()

		err := retryWithTimeoutAndJitterDQL(ctx, RetryConfig{
			MaxRetryDuration: w.Config.MaxRetryDuration,
			MaxBackoff:       w.Config.MaxBackoff,
			InitialBackoff:   w.Config.InitialBackoff,
		}, func() error {
			return w.Repo.InsertUser(ctx, event.ID, event.Username)
		})

		if err != nil {
			zap.S().Errorw("Failed to process DLQ message after retries, sending to final DLQ",
				"error", err, "user_id", event.ID, "final_dlq", w.Config.SourceTopic)

			dlqErr := requeueToDLQ(w.Producer, w.Config.SourceTopic, event)
			if dlqErr != nil {
				zap.S().Errorw("Failed to requeue", "error", dlqErr)
			}
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}

// -----------------------------
// Retry with backoff & jitter
// -----------------------------
func retryWithTimeoutAndJitterDQL(ctx context.Context, cfg RetryConfig, fn func() error) error {
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

			zap.S().Warnw("DLQ worker retryable error, will retry", "error", err, "next_backoff", backoff)

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

// -----------------------------
// Send to DLQ (reutilizado)
// -----------------------------
func requeueToDLQ(producer sarama.SyncProducer, topic string, event UserCreatedEvent) error {
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

	zap.S().Infow("Requeued message to DLQ", "topic", topic, "user_id", event.ID)
	return nil
}
