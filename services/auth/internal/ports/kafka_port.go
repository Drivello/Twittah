package ports

type UserEventProducerPort interface {
	PublishUserCreated(id int64, username string, eventType string) error
}
