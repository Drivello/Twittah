package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

func SendToDLQ(producer sarama.SyncProducer, topic string, event interface{}) error {
	if producer == nil {
		return fmt.Errorf("DLQ producer is nil")
	}
	if topic == "" {
		return fmt.Errorf("DLQ topic is empty")
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal DLQ payload: %w", err)
	}

	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(payload),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to DLQ: %w", err)
	}

	common.Logger().Info("Message sent to DLQ", zap.String("topic", topic))
	return nil
}
