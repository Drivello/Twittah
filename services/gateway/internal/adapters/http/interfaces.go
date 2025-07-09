package http

type KafkaTopicLister interface {
	Topics() ([]string, error)
}
