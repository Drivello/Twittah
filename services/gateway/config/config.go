package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/Drivello/Twittah/services/gateway/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/gateway/internal/common"
)

// GatewayConfig holds all configuration for the Gateway Service.
type GatewayConfig struct {
	Port                string
	KafkaBrokers        []string
	KafkaAuthTopic      string
	KafkaFollowTopic    string
	KafkaTweetTopic     string
	LogLevel            string
	UserMicroserviceURL string
}

// LoadConfig loads all required configuration from environment variables.
// Returns a pointer to GatewayConfig and error if any variable is missing.
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

	followTopic := os.Getenv("KAFKA_FOLLOWS_TOPIC")
	if followTopic == "" {
		return nil, fmt.Errorf("KAFKA_FOLLOWS_TOPIC env var required")
	}
	authTopic := os.Getenv("KAFKA_AUTH_TOPIC")
	if authTopic == "" {
		return nil, fmt.Errorf("KAFKA_AUTH_TOPIC env var required")
	}
	tweetTopic := os.Getenv("KAFKA_TWEET_TOPIC")
	if tweetTopic == "" {
		return nil, fmt.Errorf("KAFKA_TWEET_TOPIC env var required")
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	userMicroserviceURL := os.Getenv("USER_MICROSERVICE_URL")
	if userMicroserviceURL == "" {
		return nil, fmt.Errorf("USER_MICROSERVICE_URL env var required")
	}
	return &GatewayConfig{
		Port:                port,
		KafkaBrokers:        brokers,
		KafkaAuthTopic:      authTopic,
		KafkaFollowTopic:    followTopic,
		KafkaTweetTopic:     tweetTopic,
		LogLevel:            logLevel,
		UserMicroserviceURL: userMicroserviceURL,
	}, nil
}

func InitKafkaProducers(cfg *GatewayConfig) (*kafka.AuthProducer, *kafka.UserEventProducer, *kafka.TweetEventProducer, error) {
	authProducer, err := kafka.NewAuthProducer(cfg.KafkaBrokers, cfg.KafkaAuthTopic)
	if err != nil {
		common.Logger().Errorw("Failed to create auth kafka producer", "error", err)
		return nil, nil, nil, err
	}
	userProducer, err := kafka.NewUserEventProducer(cfg.KafkaBrokers, cfg.KafkaFollowTopic)
	if err != nil {
		common.Logger().Errorw("Failed to create user kafka producer", "error", err)
		return nil, nil, nil, err
	}
	tweetProducer, err := kafka.NewTweetEventProducer(cfg.KafkaBrokers, cfg.KafkaTweetTopic)
	if err != nil {
		common.Logger().Errorw("Failed to create tweet kafka producer", "error", err)
		return nil, nil, nil, err
	}
	return authProducer, userProducer, tweetProducer, nil
}
