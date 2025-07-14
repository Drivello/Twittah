package ports

import "context"

// EntClient abstracts the DB client for health checks or other operations
// Implemented by adapters/postgres
// Used by adapters/http
// Should be mocked in tests
type EntClient interface {
	Ping(ctx context.Context) error
}

// KafkaTopicLister abstracts Kafka topic listing for health checks
// Implemented by adapters/kafka
// Used by adapters/http
// Should be mocked in tests
type KafkaTopicLister interface {
	Topics() ([]string, error)
}
