package ports

type UserEventProducerPort interface {
	PublishUserCreated(id int64, username string) error
}
