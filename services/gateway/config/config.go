package config

import (
	"strings"

	"github.com/Drivello/Twittah/services/gateway/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/gateway/internal/common"
)

// GatewayConfig holds all configuration for the Gateway Service.
type GatewayConfig struct {
	Port                 string
	KafkaBrokers         []string
	KafkaAuthTopic       string
	KafkaUserTopic       string
	KafkaTweetTopic      string
	LogLevel             string
	UserMicroserviceURL  string
	TweetMicroserviceURL string
}

// LoadConfig loads all required configuration from environment variables.
// Returns a pointer to GatewayConfig and error if any variable is missing.
func LoadConfig() (*GatewayConfig, error) {
	// Service
	port := mustGetEnv("GATEWAY_PORT")
	logLevel := getEnvOrDefault("LOG_LEVEL", "info")

	// Kafka
	brokersStr := mustGetEnv("KAFKA_BROKERS")
	brokers := strings.Split(brokersStr, ",")

	// Kafka Topics
	userTopic := mustGetEnv("KAFKA_USER_TOPIC")
	authTopic := mustGetEnv("KAFKA_AUTH_TOPIC")
	tweetTopic := mustGetEnv("KAFKA_TWEET_TOPIC")

	// Microservices
	userMicroserviceURL := mustGetEnv("USER_MICROSERVICE_URL")
	tweetMicroserviceURL := mustGetEnv("TWEET_MICROSERVICE_URL")

	return &GatewayConfig{
		Port:                 port,
		KafkaBrokers:         brokers,
		KafkaAuthTopic:       authTopic,
		KafkaUserTopic:       userTopic,
		KafkaTweetTopic:      tweetTopic,
		LogLevel:             logLevel,
		UserMicroserviceURL:  userMicroserviceURL,
		TweetMicroserviceURL: tweetMicroserviceURL,
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
