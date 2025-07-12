package config

import (
	"strings"
)

// Config holds all configuration for the UserService.
type Config struct {
	// Service
	ServicePort string
	LogLevel    string

	// Database & Redis
	PostgresDSN string
	RedisAddr   string

	// Kafka
	KafkaBrokers []string

	// Kafka User Config
	KafkaUserConsumerConfig ConsumerConfig
}

// LoadConfig loads configuration from environment variables or .env file.
// LoadConfig loads configuration from environment variables. Returns a pointer to Config and error if any variable is missing.
func LoadConfig() *Config {

	// Service
	servicePort := mustGetEnv("USER_SERVICE_PORT")
	logLevel := getEnvOrDefault("LOG_LEVEL", "info")

	// Database & Redis
	postgresDSN := mustGetEnv("USER_POSTGRES_DSN")
	redisAddr := mustGetEnv("REDIS_ADDR")

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

	return &Config{
		// Service
		ServicePort: servicePort,
		LogLevel:    logLevel,

		// Database
		PostgresDSN: postgresDSN,
		RedisAddr:   redisAddr,

		// Kafka
		KafkaBrokers: brokers,

		// Kafka User Config
		KafkaUserConsumerConfig: kafkaUserConsumerConfig,
	}
}
