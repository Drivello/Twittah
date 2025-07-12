package config

import "time"

// ConsumerConfig holds configuration for Kafka consumers.
type ConsumerConfig struct {
	Group       string
	Topic       string
	DLQTopic    string
	TTL         time.Duration
	RetryConfig RetryConfig
	DLQConfig   DLQConfig
}

// RetryConfig holds configuration for retry logic in Kafka consumers.
type RetryConfig struct {
	MaxRetryDuration time.Duration
	MaxBackoff       time.Duration
	InitialBackoff   time.Duration
}

// DLQConfig holds configuration for DLQ topics.
type DLQConfig struct {
	DLQTopic    string
	SourceTopic string
	TTL         time.Duration
}
