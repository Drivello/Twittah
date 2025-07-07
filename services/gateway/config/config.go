package config

import (
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
		port = "8081"
	}
	brokersStr := os.Getenv("KAFKA_BROKERS")
	brokers := []string{"localhost:9092"}
	if brokersStr != "" {
		brokers = strings.Split(brokersStr, ",")
	}
	followsTopic := os.Getenv("KAFKA_FOLLOWS_TOPIC")
	if followsTopic == "" {
		followsTopic = "follows"
	}
	userTopic := os.Getenv("KAFKA_USER_TOPIC")
	if userTopic == "" {
		userTopic = "user_created"
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
