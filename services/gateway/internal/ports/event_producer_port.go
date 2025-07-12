package ports

import "github.com/Drivello/Twittah/services/gateway/internal/adapters/kafka"

type EventProducerPort[T any] interface {
	PublishEvent(event kafka.KafkaEventRequest[T]) error
}
