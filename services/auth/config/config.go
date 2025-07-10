package config

import (
	"os"
	"time"

	"github.com/Drivello/Twittah/services/auth/internal/common"
	"go.uber.org/zap"
)

// RetryConfig holds configuration for retry logic in Kafka consumers.
type RetryConfig struct {
	MaxRetryDuration time.Duration
	MaxBackoff       time.Duration
	InitialBackoff   time.Duration
	DLQTopic         string
}

// DLQWorkerConfig holds configuration for the DLQ worker retry and topics.
type DLQWorkerConfig struct {
	MaxRetryDuration time.Duration
	MaxBackoff       time.Duration
	InitialBackoff   time.Duration
	SourceTopic      string
	DLQMessageTTL    time.Duration // TTL para mensajes de la DLQ
}

// Config holds all environment configuration for the Auth service.
type Config struct {
	Port                          string
	KafkaBrokers                  []string
	DatabaseURL                   string
	KafkaUserEventsTopic          string
	UserCreatedEventType          string
	KafkaAuthConsumerGroupName    string
	KafkaAuthDLQConsumerGroupName string
	Retry                         RetryConfig
	DLQWorker                     DLQWorkerConfig
}

// getEnvAsDuration parses an environment variable as a time.Duration or returns the default.
func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if valStr := os.Getenv(key); valStr != "" {
		val, err := time.ParseDuration(valStr)
		if err != nil {
			common.Logger().Error("Invalid duration for %s: %s, using default %v", zap.String("key", key), zap.String("value", valStr), zap.Duration("default", defaultVal))
			return defaultVal
		}
		return val
	}
	return defaultVal
}

// LoadConfig loads configuration from environment variables and returns a Config struct.
func LoadConfig() *Config {
	port := os.Getenv("AUTH_PORT")
	brokers := os.Getenv("KAFKA_BROKERS")
	databaseURL := os.Getenv("AUTH_POSTGRES_DSN")
	kafkaUserEventsTopic := os.Getenv("KAFKA_USER_TOPIC")
	userCreatedEventType := os.Getenv("USER_CREATED_EVENT_TYPE")
	kafkaAuthConsumerGroupName := os.Getenv("KAFKA_AUTH_CONSUMER_GROUP_NAME")
	kafkaAuthDLQConsumerGroupName := os.Getenv("KAFKA_AUTH_DLQ_CONSUMER_GROUP_NAME")
	if userCreatedEventType == "" {
		userCreatedEventType = "user_created_to_replicate"
	}
	if port == "" {
		common.Logger().Fatal("AUTH_PORT env var required")
	}
	if brokers == "" {
		common.Logger().Fatal("KAFKA_BROKERS env var required")
	}
	if databaseURL == "" {
		common.Logger().Fatal("AUTH_POSTGRES_DSN env var required")
	}
	if kafkaUserEventsTopic == "" {
		common.Logger().Fatal("KAFKA_USER_TOPIC env var required")
	}
	if kafkaAuthConsumerGroupName == "" {
		common.Logger().Fatal("KAFKA_AUTH_CONSUMER_GROUP_NAME env var required")
	}
	if kafkaAuthDLQConsumerGroupName == "" {
		common.Logger().Fatal("KAFKA_AUTH_DLQ_CONSUMER_GROUP_NAME env var required")
	}
	if kafkaUserEventsTopic == "" {
		common.Logger().Fatal("KAFKA_USER_TOPIC env var required")
	}
	if kafkaAuthConsumerGroupName == "" {
		common.Logger().Fatal("KAFKA_AUTH_DLQ_CONSUMER_GROUP_NAME env var required")
	}

	if kafkaAuthDLQConsumerGroupName == "" {
		common.Logger().Fatal("KAFKA_AUTH_DLQ_CONSUMER_GROUP_NAME env var required")
	}

	retryCfg := RetryConfig{
		MaxRetryDuration: getEnvAsDuration("RETRY_MAX_DURATION", 0),
		MaxBackoff:       getEnvAsDuration("RETRY_MAX_BACKOFF", 0),
		InitialBackoff:   getEnvAsDuration("RETRY_INITIAL_BACKOFF", 0),
		DLQTopic:         os.Getenv("RETRY_DLQ_TOPIC"),
	}
	dlqWorkerCfg := DLQWorkerConfig{
		MaxRetryDuration: getEnvAsDuration("DLQ_WORKER_MAX_DURATION", 0),
		MaxBackoff:       getEnvAsDuration("DLQ_WORKER_MAX_BACKOFF", 0),
		InitialBackoff:   getEnvAsDuration("DLQ_WORKER_INITIAL_BACKOFF", 0),
		SourceTopic:      os.Getenv("RETRY_DLQ_TOPIC"),
		DLQMessageTTL:    getEnvAsDuration("DLQ_WORKER_MESSAGE_TTL", 0),
	}

	return &Config{
		Port:                          port,
		KafkaBrokers:                  []string{brokers},
		DatabaseURL:                   databaseURL,
		KafkaUserEventsTopic:          kafkaUserEventsTopic,
		UserCreatedEventType:          userCreatedEventType,
		KafkaAuthConsumerGroupName:    kafkaAuthConsumerGroupName,
		KafkaAuthDLQConsumerGroupName: kafkaAuthDLQConsumerGroupName,
		Retry:                         retryCfg,
		DLQWorker:                     dlqWorkerCfg,
	}
}
