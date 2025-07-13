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
	RedisAddr   string

	// Kafka
	KafkaBrokers []string

	// Kafka User Consumer
	KafkaUserConsumerConfig ConsumerConfig

	// Kafka Tweet Consumer
	KafkaTweetConsumerConfig ConsumerConfig
}

func LoadConfig() *Config {
	// Service
	port := mustGetEnv("TWEET_PORT")
	logLevel := getEnvOrDefault("LOG_LEVEL", "info")

	// Database
	postgresDSN := mustGetEnv("TWEET_POSTGRES_DSN")
	redisAddr := mustGetEnv("TWEET_REDIS_ADDR")

	// Kafka
	brokersStr := mustGetEnv("KAFKA_BROKERS")
	brokers := strings.Split(brokersStr, ",")

	// Kafka User Config
	kafkaUserConsumerConfig := ConsumerConfig{
		Group:    mustGetEnv("KAFKA_USER_GROUP"),
		Topic:    mustGetEnv("KAFKA_USER_TOPIC"),
		DLQTopic: mustGetEnv("KAFKA_USER_TOPIC_DLQ"),
		TTL:      getEnvAsDuration("KAFKA_USER_TOPIC_DLQ_TTL", 0),
		RetryConfig: RetryConfig{
			MaxRetryDuration: getEnvAsDuration("KAFKA_USER_TOPIC_RETRY_MAX_DURATION", 0),
			MaxBackoff:       getEnvAsDuration("KAFKA_USER_TOPIC_RETRY_MAX_BACKOFF", 0),
			InitialBackoff:   getEnvAsDuration("KAFKA_USER_TOPIC_RETRY_INITIAL_BACKOFF", 0),
		},
		DLQConfig: DLQConfig{
			DLQTopic:    mustGetEnv("KAFKA_USER_TOPIC_DLQ"),
			SourceTopic: mustGetEnv("KAFKA_USER_TOPIC"),
			TTL:         getEnvAsDuration("KAFKA_USER_TOPIC_DLQ_TTL", 0),
		},
	}

	// Kafka Tweet Consumer Config
	kafkaTweetConsumerConfig := ConsumerConfig{
		Group:    mustGetEnv("KAFKA_TWEET_CONSUMER_GROUP_NAME"),
		Topic:    mustGetEnv("KAFKA_TWEET_TOPIC"),
		DLQTopic: mustGetEnv("KAFKA_TWEET_DLQ_TOPIC"),
		TTL:      getEnvAsDuration("KAFKA_TWEET_DLQ_TTL", 168*time.Hour),
		RetryConfig: RetryConfig{
			MaxRetryDuration: getEnvAsDuration("KAFKA_TWEET_RETRY_MAX_DURATION", 2*time.Hour),
			MaxBackoff:       getEnvAsDuration("KAFKA_TWEET_RETRY_MAX_BACKOFF", 300*time.Second),
			InitialBackoff:   getEnvAsDuration("KAFKA_TWEET_RETRY_INITIAL_BACKOFF", 5*time.Second),
		},
		DLQConfig: DLQConfig{
			DLQTopic:    mustGetEnv("KAFKA_TWEET_DLQ_TOPIC"),
			SourceTopic: mustGetEnv("KAFKA_TWEET_TOPIC"),
			TTL:         getEnvAsDuration("KAFKA_TWEET_DLQ_TTL", 168*time.Hour),
		},
	}

	return &Config{
		Port:                     port,
		LogLevel:                 logLevel,
		PostgresDSN:              postgresDSN,
		RedisAddr:                redisAddr,
		KafkaBrokers:             brokers,
		KafkaUserConsumerConfig:  kafkaUserConsumerConfig,
		KafkaTweetConsumerConfig: kafkaTweetConsumerConfig,
	}
}
