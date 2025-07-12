package kafka

import (
	"fmt"

	"github.com/Drivello/Twittah/services/auth/internal/common"
	"github.com/Drivello/Twittah/services/auth/internal/ports"
	"go.uber.org/zap"
)

func SendToDLQ(producer ports.EventProducerPort, eventType string, payload interface{}) error {
	err := producer.PublishEvent(eventType, payload.(KafkaEventRequest).Payload)
	if err != nil {
		return fmt.Errorf("failed to send message to DLQ: %w", err)
	}

	common.Logger().Info("Message sent to DLQ", zap.String("event_type", eventType))
	return nil
}
