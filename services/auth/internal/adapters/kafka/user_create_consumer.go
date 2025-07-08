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

type UserCreateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserCreateConsumer handles Kafka messages for user creation events.

type UserCreateConsumer struct {
	UserUC   ports.UserUseCasePort
	Producer sarama.SyncProducer
	Config   config.RetryConfig
}

func (c *UserCreateConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (c *UserCreateConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim processes messages from the Kafka topic and attempts to create users.
// On failure after retries, the message is sent to the DLQ.
func (c *UserCreateConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	zap.L().Info("[KafkaConsumer] ConsumeClaim started for user_created topic")
	for msg := range claim.Messages() {
		// Process each message in the Kafka topic
		var req UserCreateRequest
		if err := json.Unmarshal(msg.Value, &req); err != nil {
			zap.L().Error("[KafkaConsumer] Failed to unmarshal user create request", zap.Error(err))
			sess.MarkMessage(msg, "")
			continue // Skip invalid message
		}

		ctx, cancel := context.WithTimeout(sess.Context(), c.Config.MaxRetryDuration)
		defer cancel()

		event := struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		}

			// Attempt user creation with retry and backoff
		err := retryWithTimeoutAndJitter(ctx, c.Config, func() error {
			user := &domain.User{
				Username: req.Username,
				Email:    req.Email,
				Password: req.Password,
			}
			_, err := c.UserUC.RegisterUser(ctx, user)
			return err
		})

		if err != nil {
			// If all retries fail, send to DLQ
			zap.L().Error("[KafkaConsumer] Failed to create user after retries, sending to DLQ", zap.Error(err), zap.String("dlq_topic", c.Config.DLQTopic))
			dlqErr := sendToDLQ(c.Producer, c.Config.DLQTopic, event)
			if dlqErr != nil {
				zap.L().Error("[KafkaConsumer] Failed to send message to DLQ", zap.Error(dlqErr))
			}
			sess.MarkMessage(msg, "")
			continue
		}

		// Successfully processed message
		sess.MarkMessage(msg, "")
	}
	return nil
}

// retryWithTimeoutAndJitter retries a function using exponential backoff and jitter until timeout.
func retryWithTimeoutAndJitter(ctx context.Context, cfg config.RetryConfig, fn func() error) error {
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

			zap.L().Warn("Retryable error, will retry", zap.Error(err), zap.Duration("next_backoff", backoff))

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

// sendToDLQ serializes the event and sends it to the DLQ topic using the provided Kafka producer.
func sendToDLQ(producer sarama.SyncProducer, topic string, event interface{}) error {
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
