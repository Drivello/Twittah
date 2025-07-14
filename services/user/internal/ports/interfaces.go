package ports

import "context"

type EntClient interface {
	Ping(ctx context.Context) error
}

type KafkaTopicLister interface {
	Topics() ([]string, error)
}
