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
	Retry        RetryConfig
	DLQWorker    DLQWorkerConfig
}

// getEnv returns the value of the environment variable or a default if not set.
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
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
	port := getEnv("AUTH_PORT", "8083")
	brokers := getEnv("KAFKA_BROKERS", "")
	if brokers == "" {
		zap.L().Fatal("KAFKA_BROKERS env var required")
	}

	retryCfg := RetryConfig{
		MaxRetryDuration: getEnvAsDuration("RETRY_MAX_DURATION", 2*time.Hour),
		MaxBackoff:       getEnvAsDuration("RETRY_MAX_BACKOFF", 30*time.Second),
		InitialBackoff:   getEnvAsDuration("RETRY_INITIAL_BACKOFF", 500*time.Millisecond),
		DLQTopic:         getEnv("RETRY_DLQ_TOPIC", "user_created_dlq"),
	}

	dlqWorkerCfg := DLQWorkerConfig{
		MaxRetryDuration: getEnvAsDuration("DLQ_WORKER_MAX_DURATION", 1*time.Hour),
		MaxBackoff:       getEnvAsDuration("DLQ_WORKER_MAX_BACKOFF", 5*time.Minute),
		InitialBackoff:   getEnvAsDuration("DLQ_WORKER_INITIAL_BACKOFF", 5*time.Second),
		SourceTopic:      getEnv("DLQ_WORKER_SOURCE_TOPIC", "user_created_dlq"),
		FinalDLQTopic:    getEnv("DLQ_WORKER_FINAL_DLQ_TOPIC", "user_created_dlq_final"),
	}

	return &Config{
		Port:         port,
		KafkaBrokers: []string{brokers},
		Retry:        retryCfg,
		DLQWorker:    dlqWorkerCfg,
	}
}
