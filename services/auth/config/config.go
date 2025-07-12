package config

import (
	"strings"
	"time"
)

type Config struct {
	// Service
	Port     string
	LogLevel string

	// Database
	PostgresDSN string

	// Kafka
	KafkaBrokers   []string
	KafkaUserTopic string

	// Kafka Auth Consumer
	KafkaUserConsumerConfig ConsumerConfig
}

// LoadConfig loads configuration from environment variables or .env file.
func LoadConfig() *Config {
	// Service
	port := mustGetEnv("AUTH_PORT")
	logLevel := getEnvOrDefault("LOG_LEVEL", "info")

	// Database
	postgresDSN := mustGetEnv("AUTH_POSTGRES_DSN")

	// Kafka
	brokersStr := mustGetEnv("KAFKA_BROKERS")
	brokers := strings.Split(brokersStr, ",")
	userTopic := mustGetEnv("KAFKA_USER_TOPIC")

	// Kafka Auth Consumer Config
	authConsumerConfig := ConsumerConfig{
		Group:    mustGetEnv("KAFKA_AUTH_CONSUMER_GROUP_NAME"),
		Topic:    mustGetEnv("KAFKA_AUTH_TOPIC"),
		DLQTopic: mustGetEnv("KAFKA_AUTH_DLQ_TOPIC"),
		TTL:      getEnvAsDuration("KAFKA_AUTH_DLQ_TTL", 168*time.Hour),
		RetryConfig: RetryConfig{
			MaxRetryDuration: getEnvAsDuration("KAFKA_AUTH_RETRY_MAX_DURATION", 2*time.Hour),
			MaxBackoff:       getEnvAsDuration("KAFKA_AUTH_RETRY_MAX_BACKOFF", 300*time.Second),
			InitialBackoff:   getEnvAsDuration("KAFKA_AUTH_RETRY_INITIAL_BACKOFF", 5*time.Second),
		},
		DLQConfig: DLQConfig{
			DLQTopic:    mustGetEnv("KAFKA_AUTH_DLQ_TOPIC"),
			SourceTopic: mustGetEnv("KAFKA_USER_TOPIC"),
			TTL:         getEnvAsDuration("KAFKA_AUTH_DLQ_TTL", 168*time.Hour),
		},
	}

	return &Config{
		Port:                    port,
		LogLevel:                logLevel,
		PostgresDSN:             postgresDSN,
		KafkaBrokers:            brokers,
		KafkaUserTopic:          userTopic,
		KafkaUserConsumerConfig: authConsumerConfig,
	}
}
