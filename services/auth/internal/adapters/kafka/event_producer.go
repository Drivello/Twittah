package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/Drivello/Twittah/services/auth/internal/common"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// EventProducer wraps a sarama.SyncProducer for publishing events.
type EventProducer struct {
	Producer sarama.SyncProducer
	Topic    string
}

// NewEventProducer creates a new EventProducer for the given brokers and topic.
func NewEventProducer(brokers []string, topic string) (*EventProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		common.Logger().Error("[EventProducer] Failed to create SyncProducer", zap.Error(err))
		return nil, err
	}
	common.Logger().Info("[EventProducer] SyncProducer created successfully", zap.Strings("brokers ", brokers), zap.String("topic ", topic))
	return &EventProducer{Producer: producer, Topic: topic}, nil
}

// BuildKafkaEventRequest construye la request adecuada para cada tipo de evento Kafka.
func BuildKafkaEventRequest(eventType string, payload interface{}) (interface{}, error) {
	switch eventType {
	case "users.created":
		userPayload, ok := payload.(KafkaUserCreatedPayload)
		if !ok {
			// Intentar type assertion usando el package explícito
			if uptr, ok2 := payload.(*KafkaUserCreatedPayload); ok2 && uptr != nil {
				userPayload = *uptr
				ok = true
			}
		}
		if !ok {
			return nil, fmt.Errorf("payload must be KafkaUserCreatedPayload for users.created, got %T", payload)
		}
		return KafkaUserCreatedRequest{
			EventType: eventType,
			Payload:   userPayload,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported event type: %s", eventType)
	}
}

// PublishEvent publica un evento a Kafka usando el builder centralizado.
func (p *EventProducer) PublishEvent(eventType string, payload interface{}) error {
	common.Logger().Debug("[UserEventProducer] Publishing event "+eventType+" to topic "+p.Topic, zap.String("event_type", eventType))
	event, err := BuildKafkaEventRequest(eventType, payload)
	if err != nil {
		common.Logger().Error("[UserEventProducer] BuildKafkaEventRequest error", zap.String("event_type", eventType), zap.Error(err))
		return err
	}
	value, err := json.Marshal(event)
	if err != nil {
		common.Logger().Error("[UserEventProducer] Failed to marshal event "+eventType, zap.Error(err))
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: p.Topic,
		Value: sarama.ByteEncoder(value),
	}
	_, _, err = p.Producer.SendMessage(msg)
	if err != nil {
		common.Logger().Error("[UserEventProducer] Failed to publish event "+eventType, zap.Error(err))
		return err
	}
	common.Logger().Debug("[UserEventProducer] event " + eventType + " published")
	return nil
}
