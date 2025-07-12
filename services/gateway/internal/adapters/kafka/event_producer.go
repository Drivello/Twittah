package kafka

import (
	"encoding/json"

	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"github.com/IBM/sarama"
)

type EventProducer[T any] struct {
	Producer sarama.SyncProducer
	Topic    string
	Builder  KafkaEventRequestBuilder[T]
}

// KafkaEventRequestBuilder defines a function that builds Kafka event requests.
type KafkaEventRequestBuilder[T any] func(eventType string, payload T) (KafkaEventRequest[T], error)

func NewEventProducer[T any](brokers []string, topic string, builder KafkaEventRequestBuilder[T]) (*EventProducer[T], error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		common.Logger().Error("[EventProducer] Failed to create SyncProducer", "error", err)
		return nil, err
	}

	common.Logger().Info("[EventProducer] SyncProducer created successfully", "brokers", brokers, "topic", topic)
	return &EventProducer[T]{
		Producer: producer,
		Topic:    topic,
		Builder:  builder,
	}, nil
}

// PublishEvent publishes a KafkaEventRequest directly to Kafka.
func (p *EventProducer[T]) PublishEvent(event KafkaEventRequest[T]) error {
	common.Logger().Debug("[EventProducer] Publishing event "+event.EventType+" to topic "+p.Topic, "event_type", event.EventType)

	value, err := json.Marshal(event)
	if err != nil {
		common.Logger().Error("[EventProducer] Failed to marshal event "+event.EventType, "error", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.Topic,
		Value: sarama.ByteEncoder(value),
	}
	_, _, err = p.Producer.SendMessage(msg)
	if err != nil {
		common.Logger().Error("[EventProducer] Failed to publish event "+event.EventType, "error", err)
		return err
	}

	common.Logger().Debug("[EventProducer] Event " + event.EventType + " published")
	return nil
}
