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
	KafkaUserTopic      string
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

	userTopic := os.Getenv("KAFKA_USER_TOPIC")
	if userTopic == "" {
		return nil, fmt.Errorf("KAFKA_USER_TOPIC env var required")
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
		KafkaUserTopic:      userTopic,
		KafkaTweetTopic:     tweetTopic,
		LogLevel:            logLevel,
		UserMicroserviceURL: userMicroserviceURL,
	}, nil
}

func InitKafkaProducers(cfg *GatewayConfig) (
	*kafka.EventProducer[kafka.AuthPayload],
	*kafka.EventProducer[kafka.UserPayload],
	*kafka.EventProducer[kafka.TweetPayload],
	error,
) {
	// Auth producer solo soporta KafkaUserCreatePayload
	authProducer, err := kafka.NewEventProducer[kafka.AuthPayload](
		cfg.KafkaBrokers,
		cfg.KafkaAuthTopic,
		kafka.AuthEventRequestBuilder,
	)
	if err != nil {
		common.Logger().Errorw("Failed to create auth kafka producer", "error", err)
		return nil, nil, nil, err
	}

	// User producer soporta KafkaUserCreatedPayload y KafkaFollowPayload
	userProducer, err := kafka.NewEventProducer[kafka.UserPayload](
		cfg.KafkaBrokers,
		cfg.KafkaUserTopic,
		kafka.UserEventRequestBuilder,
	)
	if err != nil {
		common.Logger().Errorw("Failed to create user kafka producer", "error", err)
		return nil, nil, nil, err
	}

	// Tweet producer solo soporta KafkaTweetCreatePayload
	tweetProducer, err := kafka.NewEventProducer[kafka.TweetPayload](
		cfg.KafkaBrokers,
		cfg.KafkaTweetTopic,
		kafka.TweetEventRequestBuilder,
	)
	if err != nil {
		common.Logger().Errorw("Failed to create tweet kafka producer", "error", err)
		return nil, nil, nil, err
	}

	return authProducer, userProducer, tweetProducer, nil
}
