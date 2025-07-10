package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/Drivello/Twittah/services/auth/config"
	"github.com/Drivello/Twittah/services/auth/internal/common"
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
	UserUC    ports.UserUseCasePort
	Producer  sarama.SyncProducer
	Config    config.RetryConfig
	EventType string
}

func (c *UserCreateConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (c *UserCreateConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim processes messages from the Kafka topic and attempts to create users.
// On failure after retries, the message is sent to the DLQ (if enabled).
func (c *UserCreateConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	topic := claim.Topic()
	common.Logger().Debug("[KafkaConsumer] [UserCreateConsumer] Listening topic", zap.String("topic", topic))

	workerCount := 8 // puedes parametrizar esto con config
	queueSize := 64
	wq := NewWorkQueue(workerCount, queueSize)

	// Función que ejecuta el registro y DLQ (corre en cada worker)
	process := func(item WorkItem) {
		req := item.Request
		ctx := item.Ctx
		msg, _ := item.Msg.(sarama.ConsumerMessage)

		err := retryWithTimeoutAndJitter(ctx, c.Config, func() error {
			user := &domain.User{
				Username: req.Username,
				Email:    req.Email,
				Password: req.Password,
			}
			_, err := c.UserUC.RegisterUser(ctx, user)
			common.Logger().Debug("[KafkaConsumer] RegisterUser error", zap.Error(err))
			return err
		})

		if err != nil {
			common.Logger().Error("[KafkaConsumer] Failed to create user after retries, would send to DLQ",
				zap.Error(err),
				zap.String("dlq_topic", c.Config.DLQTopic))

			event := buildEventWithType(req, c.EventType)
			dlqErr := sendToDLQ(c.Producer, c.Config.DLQTopic, event)
			if dlqErr != nil {
				common.Logger().Error("[KafkaConsumer] Failed to send message to DLQ", zap.Error(dlqErr))
			}

		}

		// Marcar mensaje como procesado (éxito o DLQ)
		sess.MarkMessage(&msg, "")
	}

	// Start WorkQueue with retry duration
	wq.Start(process, c.Config.MaxRetryDuration)
	defer wq.Stop()

	for msg := range claim.Messages() {
		var req UserCreateRequest
		if err := json.Unmarshal(msg.Value, &req); err != nil {
			common.Logger().Error("[KafkaConsumer] Failed to unmarshal user create request", zap.Error(err))
			sess.MarkMessage(msg, "")
			continue
		}

		// Submit WorkItem with a simple background context (timeout is in worker)
		err := wq.Submit(WorkItem{
			Ctx:     context.Background(),
			Request: req,
			Msg:     *msg,
		})
		if err != nil {
			common.Logger().Warn("[KafkaConsumer] WorkQueue full, would send to DLQ", zap.Error(err))

			event := buildEventWithType(req, c.EventType)
			dlqErr := sendToDLQ(c.Producer, c.Config.DLQTopic, event)
			if dlqErr != nil {
				common.Logger().Error("Failed to send message to DLQ (queue full)", zap.Error(dlqErr))
			}

			common.Logger().Error("[KafkaConsumer] Dropping message because WorkQueue is full and DLQ is disabled")
			sess.MarkMessage(msg, "")
		}
	}
	return nil
}

// retryWithTimeoutAndJitter retries a function using exponential backoff and jitter until timeout.
func retryWithTimeoutAndJitter(ctx context.Context, cfg config.RetryConfig, fn func() error) error {
	backoff := cfg.InitialBackoff

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry timed out after %v: %w", cfg.MaxRetryDuration, ctx.Err())
		default:
			err := fn()
			if err == nil {
				return nil // success
			}

			common.Logger().Warn("Retryable error, will retry",
				zap.Error(err),
				zap.Duration("next_backoff", backoff))

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

	common.Logger().Info("Message sent to DLQ", zap.String("topic", topic))
	return nil
}

func buildEventWithType(req UserCreateRequest, eventType string) map[string]interface{} {
	return map[string]interface{}{
		"event_type": eventType,
		"username":   req.Username,
		"email":      req.Email,
		"password":   req.Password,
	}
}
