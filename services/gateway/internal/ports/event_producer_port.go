package ports

type EventProducerPort[T any] interface {
	PublishEvent(eventType string, payload T) error
}
