package http

import "context"

// EntClient abstracts DB ping for healthcheck
// Implement Ping(ctx) error in your Ent client wrapper
// Example implementation: db.DB.PingContext(ctx)
type EntClient interface {
	Ping(ctx context.Context) error
}

// RedisPinger abstracts Redis ping
// Implement Ping(ctx) error in your Redis client wrapper
type RedisPinger interface {
	Ping(ctx context.Context) error
}

// KafkaTopicLister abstracts Kafka topic listing
// Implement Topics() ([]string, error) in your Kafka client wrapper
type KafkaTopicLister interface {
	Topics() ([]string, error)
}
