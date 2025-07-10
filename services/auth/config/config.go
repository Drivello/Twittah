package config

import (
	"os"
	"time"

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
	FinalDLQTopic    string
}

// Config holds all environment configuration for the Auth service.
type Config struct {
	Port         string
	KafkaBrokers []string
	DatabaseURL  string
	Retry        RetryConfig
	DLQWorker    DLQWorkerConfig
}

// getEnvAsDuration parses an environment variable as a time.Duration or returns the default.
func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if valStr := os.Getenv(key); valStr != "" {
		val, err := time.ParseDuration(valStr)
		if err != nil {
			zap.L().Error("Invalid duration for %s: %s, using default %v", zap.String("key", key), zap.String("value", valStr), zap.Duration("default", defaultVal))
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

	if port == "" {
		zap.L().Fatal("AUTH_PORT env var required")
	}
	if brokers == "" {
		zap.L().Fatal("KAFKA_BROKERS env var required")
	}
	if databaseURL == "" {
		zap.L().Fatal("AUTH_POSTGRES_DSN env var required")
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
		SourceTopic:      os.Getenv("DLQ_WORKER_SOURCE_TOPIC"),
		FinalDLQTopic:    os.Getenv("DLQ_WORKER_FINAL_DLQ_TOPIC"),
	}

	return &Config{
		Port:         port,
		KafkaBrokers: []string{brokers},
		DatabaseURL:  databaseURL,
		Retry:        retryCfg,
		DLQWorker:    dlqWorkerCfg,
	}
}
