package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Drivello/Twittah/services/user/ent"
	"go.uber.org/zap"
)

// Config holds all configuration for the UserService.
type Config struct {
	PostgresDSN  string
	RedisAddr    string
	KafkaBrokers []string
	ServicePort  string
}

func InitLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
	return logger
}

func InitEntClient(postgresDSN string) (*ent.Client, error) {
	client, err := ent.Open("postgres", postgresDSN)
	if err != nil {
		return nil, err
	}
	if err := client.Schema.Create(context.Background()); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}

func ParseKafkaBrokers(brokers []string) []string {
	if len(brokers) == 1 && strings.Contains(brokers[0], ",") {
		return strings.Split(brokers[0], ",")
	}
	return brokers
}

// LoadConfig loads configuration from environment variables or .env file.
// LoadConfig loads configuration from environment variables. Returns a pointer to Config and error if any variable is missing.
func LoadConfig() (*Config, error) {
	postgresDSN := os.Getenv("POSTGRES_DSN")
	if postgresDSN == "" {
		return nil, fmt.Errorf("POSTGRES_DSN env var required")
	}
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		return nil, fmt.Errorf("REDIS_ADDR env var required")
	}
	brokersStr := os.Getenv("KAFKA_BROKERS")
	if brokersStr == "" {
		return nil, fmt.Errorf("KAFKA_BROKERS env var required")
	}
	brokers := strings.Split(brokersStr, ",")
	servicePort := os.Getenv("USER_SERVICE_PORT")
	if servicePort == "" {
		return nil, fmt.Errorf("USER_SERVICE_PORT env var required")
	}
	return &Config{
		PostgresDSN:  postgresDSN,
		RedisAddr:    redisAddr,
		KafkaBrokers: brokers,
		ServicePort:  servicePort,
	}, nil
}
