package config

import (
	"fmt"
	"os"
	"strings"
	"github.com/Drivello/Twittah/services/gateway/internal/adapters/kafka"
	"go.uber.org/zap"
)

type GatewayConfig struct {
	Port            string
	KafkaBrokers    []string
	KafkaFollowsTopic string
	KafkaUserTopic  string
}

func LoadConfig() (*GatewayConfig, error) {
	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		return nil, fmt.Errorf("GATEWAY_PORT env var required")
	}
	brokersStr := os.Getenv("KAFKA_BROKERS")
	if brokersStr == "" {
		return nil, fmt.Errorf("KAFKA_BROKERS env var required")
	}
	brokers := strings.Split(brokersStr, ",")
	followsTopic := os.Getenv("KAFKA_FOLLOWS_TOPIC")
	if followsTopic == "" {
		return nil, fmt.Errorf("KAFKA_FOLLOWS_TOPIC env var required")
	}
	userTopic := os.Getenv("KAFKA_USER_TOPIC")
	if userTopic == "" {
		return nil, fmt.Errorf("KAFKA_USER_TOPIC env var required")
	}
	return &GatewayConfig{
		Port: port,
		KafkaBrokers: brokers,
		KafkaFollowsTopic: followsTopic,
		KafkaUserTopic: userTopic,
	}, nil
}

func InitLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

func InitKafkaProducers(cfg *GatewayConfig) (*kafka.FollowEventProducer, *kafka.UserEventProducer, error) {
	followProducer, err := kafka.NewFollowEventProducer(cfg.KafkaBrokers, cfg.KafkaFollowsTopic)
	if err != nil {
		return nil, nil, err
	}
	userProducer, err := kafka.NewUserEventProducer(cfg.KafkaBrokers, cfg.KafkaUserTopic)
	if err != nil {
		return nil, nil, err
	}
	return followProducer, userProducer, nil
}
