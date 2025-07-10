package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Drivello/Twittah/services/user/internal/ports"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// -----------------------------
// Retry configuration via environment
// -----------------------------
type RetryConfig struct {
	MaxRetryDuration time.Duration
	MaxBackoff       time.Duration
	InitialBackoff   time.Duration
	DLQTopic         string
}

func LoadRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetryDuration: getEnvAsDuration("RETRY_MAX_DURATION", 2*time.Hour),
		MaxBackoff:       getEnvAsDuration("RETRY_MAX_BACKOFF", 30*time.Second),
		InitialBackoff:   getEnvAsDuration("RETRY_INITIAL_BACKOFF", 500*time.Millisecond),
		DLQTopic:         getEnv("RETRY_DLQ_TOPIC", "user.created.dlq"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if valStr := os.Getenv(key); valStr != "" {
		val, err := time.ParseDuration(valStr)
		if err != nil {
			zap.S().Warnw("Invalid duration format, using default", "key", key, "value", valStr, "error", err)
			return defaultVal
		}
		return val
	}
	return defaultVal
}

// -----------------------------
// Event definition
// -----------------------------
type UserCreatedEvent struct {
	EventType string `json:"event_type"`
	ID        string `json:"id"`
	Username  string `json:"username"`
}

// -----------------------------
// Kafka Consumer Handler
// -----------------------------
type UserCreatedHandler struct {
	Repo     ports.UserRepository
	Producer sarama.SyncProducer
	Config   RetryConfig
}

func (h *UserCreatedHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *UserCreatedHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *UserCreatedHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var event UserCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			zap.S().Errorw("Failed to unmarshal user_created event", "error", err)
			sess.MarkMessage(msg, "")
			continue
		}

		if event.EventType == "user_created" {
			ctx, cancel := context.WithTimeout(sess.Context(), h.Config.MaxRetryDuration)
			defer cancel()

			err := retryWithTimeoutAndJitter(ctx, h.Config, func() error {
				return h.Repo.InsertUser(ctx, event.ID, event.Username)
			})

			if err != nil {
				zap.S().Errorw("Failed to insert user after retries, sending to DLQ",
					"error", err,
					"user_id", event.ID,
					"dlq_topic", h.Config.DLQTopic,
				)

				dlqErr := sendToDLQ(h.Producer, h.Config.DLQTopic, event)
				if dlqErr != nil {
					zap.S().Errorw("Failed to send message to DLQ", "error", dlqErr)
				}
			}
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}

// -----------------------------
// Retry with timeout and jitter
// -----------------------------
func retryWithTimeoutAndJitter(ctx context.Context, cfg RetryConfig, fn func() error) error {
	backoff := cfg.InitialBackoff
	rand.Seed(time.Now().UnixNano()) // for jitter

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry timed out after %v: %w", cfg.MaxRetryDuration, ctx.Err())
		default:
			err := fn()
			if err == nil {
				return nil // success
			}

			zap.S().Warnw("Retryable error, will retry", "error", err, "next_backoff", backoff)

			// Sleep with jitter (randomize between 50% and 150% of backoff)
			jitter := time.Duration(rand.Int63n(int64(backoff))) / 2
			sleepDuration := backoff + jitter
			time.Sleep(sleepDuration)

			// Exponential backoff up to max
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
// Send to DLQ
// -----------------------------
func sendToDLQ(producer sarama.SyncProducer, topic string, event UserCreatedEvent) error {
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

	zap.S().Infow("Message sent to DLQ", "topic", topic, "user_id", event.ID)
	return nil
}
