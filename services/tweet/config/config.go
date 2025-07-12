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

	// Kafka Tweet Consumer Config
	tweetConsumerConfig := ConsumerConfig{
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
		KafkaTweetConsumerConfig: tweetConsumerConfig,
	}
}
