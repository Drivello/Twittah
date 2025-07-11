package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"github.com/Drivello/Twittah/services/auth/internal/common"
)

func SendToDLQ(producer sarama.SyncProducer, topic string, event interface{}) error {
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
