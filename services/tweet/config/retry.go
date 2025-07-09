package config

import (
	"os"
	"time"
)

type RetryConfig struct {
	MaxRetryDuration time.Duration
	MaxBackoff       time.Duration
	InitialBackoff   time.Duration
	DLQTopic         string
}

type DLQWorkerConfig struct {
	MaxRetryDuration time.Duration
	MaxBackoff       time.Duration
	InitialBackoff   time.Duration
	SourceTopic      string
	FinalDLQTopic    string
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return defaultVal
	}
	return d
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func LoadRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetryDuration: getEnvAsDuration("RETRY_MAX_DURATION", 2*time.Hour),
		MaxBackoff:       getEnvAsDuration("RETRY_MAX_BACKOFF", 30*time.Second),
		InitialBackoff:   getEnvAsDuration("RETRY_INITIAL_BACKOFF", 500*time.Millisecond),
		DLQTopic:         getEnv("RETRY_DLQ_TOPIC", "tweets.dlq"),
	}
}

func LoadDLQWorkerConfig() DLQWorkerConfig {
	return DLQWorkerConfig{
		MaxRetryDuration: getEnvAsDuration("DLQ_WORKER_MAX_DURATION", 1*time.Hour),
		MaxBackoff:       getEnvAsDuration("DLQ_WORKER_MAX_BACKOFF", 5*time.Minute),
		InitialBackoff:   getEnvAsDuration("DLQ_WORKER_INITIAL_BACKOFF", 5*time.Second),
		SourceTopic:      getEnv("DLQ_WORKER_SOURCE_TOPIC", "tweets.dlq"),
		FinalDLQTopic:    getEnv("DLQ_WORKER_FINAL_DLQ_TOPIC", "tweets.dlq.final"),
	}
}
