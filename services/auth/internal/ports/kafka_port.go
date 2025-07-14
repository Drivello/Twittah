package ports

type EventProducerPort interface {
	PublishEvent(eventType string, payload interface{}) error
}
